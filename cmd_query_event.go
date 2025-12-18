package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"errors"
	"time"

	"github.com/LQR471814/nu_plugin_caldav/events"
	"github.com/LQR471814/nu_plugin_caldav/internal/db"
	"github.com/LQR471814/nu_plugin_caldav/internal/dto"
	"github.com/LQR471814/nu_plugin_caldav/internal/nuconv"
	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
)

var default_start_time = nu.ToValue(time.Time{})
var default_end_time = nu.ToValue(events.MAX_TIME)
var default_nosync = nu.ToValue(false)

var queryEventsCmd = &nu.Command{
	Signature: nu.PluginSignature{
		Name:        "caldav query events",
		Category:    "Network",
		Desc:        "Reads raw events objects from a given calendar.",
		SearchTerms: []string{"caldav", "query", "events"},
		Named: []nu.Flag{
			{
				Long:    "start",
				Short:   's',
				Desc:    "Filter for all events after this start time.",
				Shape:   syntaxshape.DateTime(),
				Default: &default_start_time,
			},
			{
				Long:    "end",
				Short:   'e',
				Desc:    "Filter for all events before this end time.",
				Shape:   syntaxshape.DateTime(),
				Default: &default_end_time,
			},
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
	start := time.Time{}
	v, ok := call.FlagValue("start")
	if ok {
		start = v.Value.(time.Time)
	}
	end := events.MAX_TIME
	v, ok = call.FlagValue("end")
	if ok {
		end = v.Value.(time.Time)
	}
	nosync := false
	v, ok = call.FlagValue("no-sync")
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
		out, err = fetchNoSync(ctx, client, calendarPath, start, end)
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
		start:        start,
		end:          end,
		calendarPath: calendarPath,
	}
	err = m.sync()
	if err != nil {
		return
	}
	var objects dto.EventObjectList
	events, err := qry.ReadEvents(ctx, calendarPath)
	for _, e := range events {
		var o dto.EventObject
		decoder := gob.NewDecoder(bytes.NewBuffer(e.Dto))
		err = decoder.Decode(&o)
		if err != nil {
			return
		}
		objects = append(objects, o)
	}
	out, err := nuconv.EventObjectListToNu(objects)
	if err != nil {
		return
	}
	err = call.ReturnValue(ctx, out)
	return
}

func fetchNoSync(ctx context.Context, client *caldav.Client, calendarPath string, start, end time.Time) (out nu.Value, err error) {
	objects, err := client.QueryCalendar(ctx, calendarPath, &caldav.CalendarQuery{
		CompFilter: caldav.CompFilter{
			Name: ical.CompCalendar,
			Comps: []caldav.CompFilter{{
				Name:  ical.CompEvent,
				Start: start,
				End:   end,
			}},
		},
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
	start, end   time.Time
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

	syncToken, err := txqry.ReadCalendarSyncToken(m.ctx, m.calendarPath)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}
	nextSyncToken, err := m.performSync(txqry, syncToken.String)
	if err != nil {
		return
	}
	err = txqry.UpdateCalSyncToken(m.ctx, db.UpdateCalSyncTokenParams{
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
