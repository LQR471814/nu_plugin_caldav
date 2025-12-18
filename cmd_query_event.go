package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"errors"
	"runtime"
	"sync"

	"github.com/LQR471814/nu_plugin_caldav/internal/db"
	"github.com/LQR471814/nu_plugin_caldav/internal/dto"
	"github.com/LQR471814/nu_plugin_caldav/internal/nuconv"
	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
)

var default_nosync = nu.ToValue(false)

var queryEventsCmd = &nu.Command{
	Signature: nu.PluginSignature{
		Name:        "caldav query events",
		Category:    "Network",
		Desc:        "Reads raw events objects from a given calendar.",
		SearchTerms: caldavKeywordsQuery("events"),
		Named: []nu.Flag{
			{
				Long:    "no-sync",
				Short:   'f',
				Desc:    "Query events without syncing.",
				Shape:   syntaxshape.Boolean(),
				Default: &default_nosync,
			},
		},
		RequiredPositional: []nu.PositionalArg{
			{
				Name:  "calendar_path",
				Desc:  "The `path` attribute of the calendar record returned by `caldav query calendars`.",
				Shape: syntaxshape.String(),
			},
		},
		InputOutputTypes: []nu.InOutTypes{
			{
				In:  types.Nothing(),
				Out: nuconv.EventObjectListType,
			},
		},
	},
	OnRun: queryEventsCmdExec,
}

func init() {
	commands = append(commands, queryEventsCmd)
}

func queryEventsCmdExec(ctx context.Context, call *nu.ExecCommand) (err error) {
	// parse flags
	calendarPath, err := tryCast[string](call.Positional[0])
	if err != nil {
		return
	}
	nosync := false
	v, ok := call.FlagValue("no-sync")
	if ok {
		nosync = v.Value.(bool)
	}

	// execution
	client, err := getClient(ctx, call)
	if err != nil {
		return
	}

	if nosync {
		var out nu.Value
		out, err = fetchNoSync(ctx, client, calendarPath)
		if err != nil {
			return
		}
		err = call.ReturnValue(ctx, out)
		return
	}

	driver, qry, err := db.Open(ctx)
	if err != nil {
		return
	}
	defer driver.Close()
	m := syncManager{
		ctx:          ctx,
		client:       client,
		driver:       driver,
		qry:          qry,
		calendarPath: calendarPath,
	}
	err = m.sync()
	if err != nil {
		return
	}

	output, err := call.ReturnListStream(ctx)
	if err != nil {
		return
	}
	defer close(output)

	workerCount := runtime.NumCPU()

	errs := make(chan error)
	events := make(chan db.ReadEventsRow, runtime.NumCPU())
	wg := sync.WaitGroup{}

	for range workerCount {
		wg.Add(1)
		go func() { // process events concurrently and send them to output stream
			defer wg.Done()
			for e := range events {
				var obj dto.EventObject
				decoder := gob.NewDecoder(bytes.NewBuffer(e.Dto))
				err = decoder.Decode(&obj)
				if err != nil {
					errs <- err
					continue
				}
				nuobj, err := nuconv.EventObjectToNu(obj)
				if err != nil {
					errs <- err
					continue
				}
				output <- nuobj
			}
		}()
	}

	go func() { // pull events in from database and send them to be processed
		err = qry.ReadEvents(ctx, calendarPath, events)
		if err != nil {
			errs <- err
		}
		close(events)
	}()

	go func() { // only close errors channel after confirming all workers have exited
		wg.Wait()
		close(errs)
	}()

	var errlist []error // collect errors, return at end
	for e := range errs {
		errlist = append(errlist, e)
	}
	if len(errlist) > 0 {
		err = errors.Join(errlist...)
		output <- nu.ToValue(err)
	}
	return
}

func fetchNoSync(ctx context.Context, client *caldav.Client, calendarPath string) (out nu.Value, err error) {
	objects, err := client.QueryCalendar(ctx, calendarPath, &caldav.CalendarQuery{
		CompRequest: caldav.CalendarCompRequest{
			Name:     ical.CompEvent,
			AllProps: true,
		},
	})
	if err != nil {
		return
	}
	dtoObjects := dto.NewEventObjectList(objects)
	out, err = nuconv.EventObjectListToNu(dtoObjects)
	return
}

type syncManager struct {
	ctx          context.Context
	client       *caldav.Client
	driver       *sql.DB
	qry          *db.Queries
	calendarPath string
}

func (m syncManager) performSync(txqry *db.Queries, syncToken string) (nextSyncToken string, err error) {
	resp, err := m.client.SyncCollection(m.ctx, m.calendarPath, &caldav.SyncQuery{
		SyncToken: syncToken,
		CompRequest: caldav.CalendarCompRequest{
			Name:     ical.CompEvent,
			AllProps: true,
		},
	})
	if err != nil {
		return
	}
	nextSyncToken = resp.SyncToken

	// sync deletes
	var deletePaths []string
	for _, path := range resp.Deleted {
		deletePaths = append(deletePaths, path)
	}
	err = txqry.DeleteEvents(m.ctx, deletePaths)
	if err != nil {
		return
	}

	// sync puts
	var updatePaths []string
	for _, u := range resp.Updated {
		updatePaths = append(updatePaths, u.Path)
	}
	updateObjects, err := m.client.MultiGetCalendar(m.ctx, m.calendarPath, &caldav.CalendarMultiGet{
		Paths: updatePaths,
		CompRequest: caldav.CalendarCompRequest{
			Name:     ical.CompEvent,
			AllProps: true,
		},
	})
	if err != nil {
		return
	}
	for _, obj := range updateObjects {
		buf := bytes.NewBuffer(nil)
		encoder := gob.NewEncoder(buf)
		err = encoder.Encode(dto.NewEventObject(obj))
		if err != nil {
			return
		}
		err = txqry.PutEvent(m.ctx, db.PutEventParams{
			Path:         obj.Path,
			CalendarPath: m.calendarPath,
			Dto:          buf.Bytes(),
		})
		if err != nil {
			return
		}
	}

	return
}

func (m syncManager) sync() (err error) {
	tx, err := m.driver.BeginTx(m.ctx, nil)
	if err != nil {
		return
	}
	defer tx.Rollback()
	txqry := m.qry.WithTx(tx)

	syncToken, err := txqry.ReadCalendar(m.ctx, m.calendarPath)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}
	nextSyncToken, err := m.performSync(txqry, syncToken.String)
	if err != nil {
		return
	}
	err = txqry.PutCalendar(m.ctx, db.PutCalendarParams{
		Path: m.calendarPath,
		SyncToken: sql.NullString{
			String: nextSyncToken,
			Valid:  true,
		},
	})
	if err != nil {
		return
	}

	tx.Commit()
	return
}
