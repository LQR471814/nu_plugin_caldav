package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
	"github.com/teambition/rrule-go"
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

	var events map[string]caldavEvent
	switch in := call.Input.(type) {
	case nu.Value:
		switch vt := in.Value.(type) {
		case []nu.Record:
			events = make(map[string]caldavEvent, len(vt))
			for _, ev := range vt {
				out := caldavEvent{}

				if _, ok := ev["uid"]; !ok {
					err = fmt.Errorf("insert Event record must have `uid` property defined and not nothing")
					return
				}
				out.Uid, err = tryCast[string](ev["uid"])
				if err != nil {
					return
				}

				if v, ok := ev["name"]; ok {
					out.Name, err = tryCast[string](v)
					if err != nil {
						return
					}
				}
				if v, ok := ev["location"]; ok {
					out.Location, err = tryCast[string](v)
					if err != nil {
						return
					}
				}
				if v, ok := ev["categories"]; ok {
					out.Categories, err = tryCast[[]string](v)
					if err != nil {
						return
					}
				}
				if v, ok := ev["start"]; ok {
					out.Start, err = tryCast[time.Time](v)
					if err != nil {
						return
					}
				}
				if v, ok := ev["end"]; ok {
					out.End, err = tryCast[time.Time](v)
					if err != nil {
						return
					}
				}
				if v, ok := ev["recurrence_id"]; ok {
					out.RId, err = tryCast[time.Time](v)
					if err != nil {
						return
					}
				}
				if v, ok := ev["recurrence_rule"]; ok {
					var rruleStr string
					rruleStr, err = tryCast[string](v)
					if err != nil {
						return
					}
					var parsed *rrule.RRule
					parsed, err = rrule.StrToRRule(rruleStr)
					out.RRule = parsed
				}
				if v, ok := ev["recurrence_exceptions"]; ok {
					out.ExDates, err = tryCast[[]time.Time](v)
					if err != nil {
						return
					}
				}
				if v, ok := ev["trigger"]; ok {
					var triggerObj nu.Record
					triggerObj, err = tryCast[nu.Record](v)
					if err != nil {
						return
					}
					if relv, ok := triggerObj["relative"]; ok {
						out.Trigger.Relative, err = tryCast[time.Duration](relv)
						if err != nil {
							return
						}
						out.Trigger.NotNone = true
					} else if absv, ok := triggerObj["absolute"]; ok {
						out.Trigger.Absolute, err = tryCast[time.Time](absv)
						if err != nil {
							return
						}
						out.Trigger.NotNone = true
					}
				}

				events[out.Uid] = out
			}
		default:
			err = fmt.Errorf("unsupported input type %T", call.Input)
			return
		}
	default:
		err = fmt.Errorf("unsupported input type %T", call.Input)
		return
	}

	res, err := fetchCalendarObjects(ctx, client, calendarPath)
	if err != nil {
		return
	}
	for _, calobj := range res {
		cal := calobj.Data
		for _, ev := range calobj.Data.Events() {
		}
	}

	_, err = client.PutCalendarObject(ctx, calendarPath)
	if err != nil {
		return
	}
}
