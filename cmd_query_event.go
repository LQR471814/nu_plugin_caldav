package main

import (
	"context"
	"time"

	"github.com/LQR471814/nu_plugin_caldav/events"
	"github.com/LQR471814/nu_plugin_caldav/internal/nutypes"
	"github.com/LQR471814/nu_plugin_caldav/internal/nutypes/conversions"
	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
)

var default_start_time = nu.ToValue(time.Time{})
var default_end_time = nu.ToValue(events.MAX_TIME)
var default_text_match_negate = nu.ToValue(false)

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
				Long:    "text-match",
				Short:   't',
				Desc:    "Filter for events that contain (or do not contain, if --text-match-negate is set) a particular string.",
				Shape:   syntaxshape.String(),
				Default: &default_end_time,
			},
			{
				Long:    "text-match-negate",
				Short:   'n',
				Desc:    "Flip the condition text-match.",
				Shape:   syntaxshape.Boolean(),
				Default: &default_text_match_negate,
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
				Out: conversions.EventObjectReplicaListType,
			},
		},
	},
	OnRun: queryEventsCmdExec,
}

func init() {
	commands = append(commands, queryEventsCmd)
}

func queryEventsCmdExec(ctx context.Context, call *nu.ExecCommand) (err error) {
	client, err := getClient(ctx, call)
	if err != nil {
		return
	}
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
	textMatch := ""
	v, ok = call.FlagValue("text-match")
	if ok {
		textMatch = v.Value.(string)
	}
	textMatchNegate := false
	v, ok = call.FlagValue("text-match-negate")
	if ok {
		textMatchNegate = v.Value.(bool)
	}
	var propFilters []caldav.PropFilter
	if textMatch != "" {
		propFilters = append(propFilters, caldav.PropFilter{
			TextMatch: &caldav.TextMatch{
				Text:            textMatch,
				NegateCondition: textMatchNegate,
			},
		})
	}

	objects, err := client.QueryCalendar(ctx, calendarPath, &caldav.CalendarQuery{
		CompFilter: caldav.CompFilter{
			Name:  ical.CompCalendar,
			Props: propFilters,
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

	var replicaObjects []nutypes.EventObjectReplica

	// each calendar object only ever stores one unique VEVENT object.
	//
	// exception:
	// if the VEVENT has recurrence overrides, the recurrence overrides will
	// come with the original VEVENT as separate VEVENT components.
	for _, obj := range objects {
		replica := nutypes.EventObjectReplica{
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
				replica.Overrides = append(replica.Overrides, nutypes.NewEventReplica(event))
				continue
			}
			replica.Main = nutypes.NewEventReplica(event)
		}

		replicaObjects = append(replicaObjects, replica)
	}

	out, err := conversions.EventObjectReplicaListToNu(replicaObjects)
	if err != nil {
		return
	}
	err = call.ReturnValue(ctx, out)
	return
}
