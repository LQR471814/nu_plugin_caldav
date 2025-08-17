package main

import (
	"context"
	"fmt"

	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
)

var upsertEventsCmd = &nu.Command{
	Signature: nu.PluginSignature{
		Name:        "caldav upsert events",
		Category:    "Network",
		Desc:        "Updates or inserts events from the given input. Note: This does not require event records to have all fields defined (the only required field is `uid`), fields with value `nothing` will not be overwritten.",
		SearchTerms: []string{"caldav", "upsert", "events"},
		RequiredPositional: []nu.PositionalArg{
			{
				Name:  "calendar_path",
				Desc:  "The `path` attribute of the calendar record returned by `caldav query calendars`.",
				Shape: syntaxshape.String(),
			},
		},
		InputOutputTypes: []nu.InOutTypes{
			{
				In:  eventsType,
				Out: types.Nothing(),
			},
		},
	},
	OnRun: upsertEventsCmdExec,
}

func init() {
	commands = append(commands, upsertEventsCmd)
}

func upsertEventsCmdExec(ctx context.Context, call *nu.ExecCommand) (err error) {
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

	calendarObjects := map[string]caldav.CalendarObject{}
	type parsedEvent struct {
		parsed   *caldavEvent
		original ical.Event
	}
	events := map[string]parsedEvent{}
	for _, calobj := range res {
		for _, component := range calobj.Data.Children {
			if component.Name != ical.CompEvent {
				continue
			}
			ev := ical.Event{
				Component: component,
			}
			var out caldavEvent
			out, err = parseEvent(ev)
			if err != nil {
				return
			}
			events[out.Uid] = parsedEvent{
				parsed:   &out,
				original: ev,
			}
			calendarObjects[out.Uid] = calobj
		}
	}

	upserted := map[string]struct{}{}
	switch in := call.Input.(type) {
	case nu.Value:
		switch vt := in.Value.(type) {
		case []nu.Record:
			for _, rec := range vt {
				uidv, ok := rec["uid"]
				if !ok {
					err = fmt.Errorf("the `uid` property is required on an Event record")
					return
				}
				var uid string
				uid, err = tryCast[string](uidv)
				if err != nil {
					return
				}

				ev := events[uid]
				upserted[calendarObjects[uid].Path] = struct{}{}
				err = ev.parsed.FromValue(nu.ToValue(rec))
				if err != nil {
					return
				}
				ev.parsed.Write(ev.original)
			}
		default:
			err = fmt.Errorf("unsupported input type %T", call.Input)
			return
		}
	default:
		err = fmt.Errorf("unsupported input type %T", call.Input)
		return
	}

	for path := range upserted {
		obj := calendarObjects[path]
		_, err = client.PutCalendarObject(ctx, path, obj.Data)
		if err != nil {
			return
		}
	}

	return
}
