package main

import (
	"context"
	"database/sql"
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
var default_text_match = nu.ToValue("")
var default_text_match_negate = nu.ToValue(false)
var default_refresh = nu.ToValue(false)

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
				Long:  "no-sync",
				Short: 'f',
				Desc:  "Force fetch all events from the server without syncing.",
				Shape: syntaxshape.Boolean(),
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

type queryEventsContext struct {
	ctx        context.Context
	driver     *sql.DB
	qry        *db.Queries
	start, end time.Time
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

	// setup

	client, err := getClient(ctx, call)
	if err != nil {
		return
	}
	driver, qry, err := db.Open(ctx)
	if err != nil {
		return
	}
	defer driver.Close()

	// execution

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

	var replicaObjects []dto.EventObject

	// each calendar object only ever stores one unique VEVENT object.
	//
	// exception:
	// if the VEVENT has recurrence overrides, the recurrence overrides will
	// come with the original VEVENT as separate VEVENT components.
	for _, obj := range objects {
		replica := dto.EventObject{
			ObjectPath: &obj.Path,
		}

		for _, component := range obj.Data.Children {
			if component.Name != ical.CompEvent {
				continue
			}
			event := events.Event{
				Event:    ical.Event{Component: component},
				Timezone: time.Local,
			}
			prop := component.Props.Get(ical.PropRecurrenceID)
			if prop != nil {
				replica.Overrides = append(replica.Overrides, dto.NewEvent(event))
				continue
			}
			replica.Main = dto.NewEvent(event)
		}

		replicaObjects = append(replicaObjects, replica)
	}

	out, err := nuconv.EventObjectListToNu(replicaObjects)
	if err != nil {
		return
	}
	err = call.ReturnValue(ctx, out)
	return
}

func (c queryEventsContext) fetchAll() {
}

func (c queryEventsContext) sync() {
}
