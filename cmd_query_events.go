package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
)

var queryEventsCmd = &nu.Command{
	Signature: nu.PluginSignature{
		Name:        "caldav query events",
		Category:    "Network",
		Desc:        "Reads events for a given calendar from CalDAV.",
		SearchTerms: []string{"caldav", "query", "events"},
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
				Out: eventsType,
			},
		},
	},
	OnRun: queryEventsCmdExec,
}

func init() {
	commands = append(commands, queryEventsCmd)
}

func fetchCalendarObjects(ctx context.Context, client *caldav.Client, calendarPath string) ([]caldav.CalendarObject, error) {
	res, err := client.QueryCalendar(ctx, calendarPath, &caldav.CalendarQuery{
		CompFilter: caldav.CompFilter{
			Name: ical.CompCalendar,
			Comps: []caldav.CompFilter{{
				Name:  ical.CompEvent,
				Start: time.Time{},
				End:   max_time,
			}},
		},
		CompRequest: caldav.CalendarCompRequest{
			Name: ical.CompCalendar,
			Comps: []caldav.CalendarCompRequest{{
				Name: ical.CompEvent,
				Props: []string{
					ical.PropUID,
					ical.PropSummary,
					ical.PropDescription,
					ical.PropLocation,
					ical.PropDateTimeStart,
					ical.PropDateTimeEnd,
					ical.PropCategories,
					ical.PropRecurrenceDates,
					ical.PropRecurrenceID,
					ical.PropRecurrenceRule,
					ical.PropTrigger,
				},
			}},
		},
	})
	return res, err
}

func queryEventsCmdExec(ctx context.Context, call *nu.ExecCommand) (err error) {
	client, err := getClientFromEnv(ctx, call)
	if err != nil {
		return
	}

	calendarPath, err := tryCast[string](call.Positional[0])
	if err != nil {
		return
	}

	res, err := fetchCalendarObjects(ctx, client, calendarPath)
	if err != nil {
		return
	}

	var events []nu.Value
	for _, calobj := range res {
		for _, ev := range calobj.Data.Events() {
			parsed, err := parseEvent(ev)
			if err != nil {
				slog.Warn("skip corrupted event", "err", err)
				continue
			}
			events = append(events, parsed.ToValue())
		}
	}

	err = call.ReturnValue(ctx, nu.ToValue(events))
	return
}

