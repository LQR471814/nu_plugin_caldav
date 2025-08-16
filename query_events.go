package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
	"github.com/teambition/rrule-go"
)

var eventsType = types.Table(types.RecordDef{
	"object_path":      types.String(),
	"uid":              types.String(),
	"name":             types.String(),
	"location":         types.String(),
	"description":      types.String(),
	"categories":       types.List(types.String()),
	"start":            types.Date(),
	"end":              types.Date(),
	"recurrence_id":    types.Date(),
	"recurrence_rule?": types.String(),
	"trigger": types.Record(types.RecordDef{
		"relative": types.Duration(),
		"absolute": types.Date(),
	}),
})

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

var max_time = time.Unix(1<<63-62135596801, 999999999)

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

	var events []map[string]any
	for _, calobj := range res {
		for _, ev := range calobj.Data.Events() {
			parsed, err := parseEvent(ev)
			if err != nil {
				slog.Warn("skip corrupted event", "err", err)
				continue
			}
			convert := parsed.ToMap()
			events = append(events, convert)
		}
	}

	err = call.ReturnValue(ctx, nu.ToValue(events))
	return
}

type eventTrigger struct {
	Relative time.Duration
	Absolute time.Time
	NotNone  bool
}

type caldavEvent struct {
	Uid         string
	Name        string
	Location    string
	Description string
	Categories  []string
	Start, End  time.Time
	ExDates     []time.Time
	RRule       *rrule.RRule
	RDates      string
	RId         time.Time
	Trigger     eventTrigger
}

func (ce caldavEvent) ToMap() (out map[string]any) {
	out = map[string]any{
		"uid":         ce.Uid,
		"name":        ce.Name,
		"location":    ce.Location,
		"description": ce.Description,
		"categories":  ce.Categories,
		"start":       ce.Start,
		"end":         ce.End,
	}
	if ce.RRule != nil {
		out["recurrence_rule"] = ce.RRule
	}
	if len(ce.ExDates) > 0 {
		out["recurrence_exceptions"] = ce.ExDates
	}
	if ce.RId != (time.Time{}) {
		out["recurrence_id"] = ce.RId
	}
	if ce.Trigger.NotNone {
		trigger := map[string]any{}
		if ce.Trigger.Relative > 0 {
			trigger["relative"] = ce.Trigger.Relative
		} else {
			trigger["absolute"] = ce.Trigger.Absolute
		}
		out["trigger"] = trigger
	}
	return
}

func parseEvent(e ical.Event) (event caldavEvent, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parse event: %w", err)
		}
	}()

	(&event).ParseUID(e)
	(&event).ParseName(e)
	(&event).ParseLocation(e)
	(&event).ParseDescription(e)

	err = (&event).ParseStart(e)
	if err != nil {
		return
	}
	err = (&event).ParseEnd(e)
	if err != nil {
		return
	}

	err = (&event).ParseCategories(e)
	if err != nil {
		return
	}
	err = (&event).ParseExceptions(e)
	if err != nil {
		return
	}
	err = (&event).ParseRecurrence(e, event.Start)
	if err != nil {
		return
	}
	err = (&event).ParseTrigger(e)
	if err != nil {
		return
	}

	if event.Uid == "" {
		err = fmt.Errorf("uid is nil")
		return
	}
	if event.Name == "" {
		err = fmt.Errorf("name is nil")
		return
	}

	return
}

func (ce *caldavEvent) ParseUID(e ical.Event) {
	uidProp := e.Props.Get(ical.PropUID)
	if uidProp == nil {
		return
	}
	ce.Uid = uidProp.Value
}

func (ce *caldavEvent) ParseName(e ical.Event) {
	nameProp := e.Props.Get(ical.PropSummary)
	if nameProp == nil {
		return
	}
	ce.Name = nameProp.Value
}

func (ce *caldavEvent) ParseDescription(e ical.Event) {
	descProp := e.Props.Get(ical.PropDescription)
	if descProp == nil {
		return
	}
	ce.Description = descProp.Value
}

func (ce *caldavEvent) ParseLocation(e ical.Event) {
	locProp := e.Props.Get(ical.PropLocation)
	if locProp == nil {
		return
	}
	ce.Location = locProp.Value
}

func (ce *caldavEvent) ParseCategories(e ical.Event) error {
	catProp := e.Props.Get(ical.PropCategories)
	if catProp == nil {
		return nil
	}
	categories, err := catProp.TextList()
	if err != nil {
		return err
	}
	ce.Categories = categories
	return nil
}

func (ce *caldavEvent) ParseStart(e ical.Event) error {
	start, err := e.DateTimeStart(time.Local)
	if err != nil {
		return err
	}
	ce.Start = start
	return nil
}

func (ce *caldavEvent) ParseEnd(e ical.Event) error {
	end, err := e.DateTimeEnd(time.Local)
	if err != nil {
		return err
	}
	ce.End = end
	return nil
}

func (ce *caldavEvent) ParseExceptions(e ical.Event) error {
	exProp := e.Props.Get(ical.PropExceptionDates)
	if exProp == nil {
		return nil
	}

	tzId := exProp.Params.Get(ical.PropTimezoneID)
	var err error
	var tz *time.Location
	if tzId != "" {
		tz, err = time.LoadLocation(tzId)
		if err != nil {
			return err
		}
	}

	var datetime time.Time
	datetime, err = exProp.DateTime(tz)
	if err != nil {
		return err
	}

	ce.ExDates = append(ce.ExDates, datetime)
	return nil
}

func (ce *caldavEvent) ParseRecurrence(e ical.Event, start time.Time) (err error) {
	recurIdProp := e.Props.Get(ical.PropRecurrenceID)
	if recurIdProp != nil && recurIdProp.Value != "" {
		ce.RId, err = recurIdProp.DateTime(time.Local)
		if err != nil {
			return
		}
	}

	rdateProp := e.Props.Get(ical.PropRecurrenceDates)
	if rdateProp != nil {
		ce.RDates = rdateProp.Value
	}

	rruleProp := e.Props.Get(ical.PropRecurrenceRule)
	if rruleProp != nil {
		var ropts *rrule.ROption
		ropts, err = rrule.StrToROptionInLocation(rruleProp.Value, time.Local)
		if err != nil {
			return
		}
		if ropts == nil {
			err = fmt.Errorf("ropts is nil")
			return
		}

		// set default dtstart to original event's starting time
		if ropts.Dtstart.Equal(time.Time{}) {
			ropts.Dtstart = start
		}

		ce.RRule, err = rrule.NewRRule(*ropts)
		if err != nil {
			return
		}
	}

	return
}

func (ce *caldavEvent) ParseTrigger(e ical.Event) (err error) {
	triggerProp := e.Props.Get(ical.PropTrigger)
	if triggerProp == nil {
		return
	}
	ce.Trigger.Relative, err = triggerProp.Duration()
	if err == nil {
		return
	}
	ce.Trigger.Absolute, err = triggerProp.DateTime(time.Local)
	if err != nil {
		return err
	}
	return
}
