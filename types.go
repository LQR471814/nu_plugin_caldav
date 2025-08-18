package main

import (
	"time"

	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/types"
	"github.com/teambition/rrule-go"
)

var max_time = time.Date(2100, time.January, 1, 0, 0, 0, 0, time.UTC)

var dateType = types.Record(types.RecordDef{
	// date represents the underlying datetime information of the
	// date & time
	"date": types.Date(),
	// if all_day is true, just the date should be considered and the time
	// ignored
	"all_day": types.Bool(),
	// if floating is true, the timezone of the date is always at the local
	// time of whoever is using the calendar and not "official" timezone should
	// be given to it
	"floating": types.Bool(),
})

type Datetime struct {
	Stamp    time.Time
	AllDay   bool
	Floating bool
}

func (d *Datetime) FromValue(v nu.Value) (err error) {
	m, err := tryCast[nu.Record](v)
	if err != nil {
		return
	}
	stamp, err := tryCast[time.Time](m["date"])
	if err != nil {
		return
	}
	allDay, err := tryCast[bool](m["all_day"])
	if err != nil {
		return
	}
	floating, err := tryCast[bool](m["floating"])
	if err != nil {
		return
	}
	d.Stamp = stamp
	d.AllDay = allDay
	d.Floating = floating
	return
}

func (d Datetime) ToValue() (out nu.Value) {
	return nu.ToValue(nu.Record{
		"date":     nu.ToValue(d.Stamp),
		"all_day":  nu.ToValue(d.AllDay),
		"floating": nu.ToValue(d.Floating),
	})
}

var eventsType = types.Table(types.RecordDef{
	// universal id for the event, unique across calendars
	"uid": types.String(),
	// name of the event
	"name": types.String(),
	// location the event occurs at
	"location": types.String(),
	// description of the event (possibly multiline)
	"description": types.String(),
	// list of "tags" for the event
	"categories": types.List(types.String()),
	// event start
	"start": dateType,
	// event end
	"end": dateType,
	// recurrence_id (optional) designates the event as an override of another
	// recurring event (if set, recurrence_set must not be set)
	"recurrence_id": dateType,
	// recurrence_set (optional) designates the event as the originator of a
	// recurring event (if set, recurrence_id must not be set)
	"recurrence_set": types.Record(types.RecordDef{
		// recurrence rule (required)
		"rule": types.String(),
		// defines which recurrences should not occur (optional)
		//
		// note: the timezone of these dates can only be one of:
		// 1. UTC
		// 2. Floating time (no timezone)
		// 3. A single other explicit time zone
		//
		// Ex.
		// 	OK:
		// 		(UTC, floating, UTC, America/Los_Angeles, America/Los_Angeles)
		// 	BAD:
		// 		(UTC, floating, UTC, Asia/Shanghai, America/Los_Angeles)
		// 	There can only be one other explicit time zone outside of UTC.
		"exceptions": types.List(dateType),
		// defines which additional dates recurrences should occur (optional)
		//
		// note: all recurrences' timezone will be:
		// 1. UTC
		// 2. Floating time (no timezone)
		// 3. Specific time zone
		"additional": types.List(dateType),
	}),
	// trigger (optional) defines a notification trigger for the event
	"trigger": types.Record(types.RecordDef{
		// if set, notification will be triggered this duration before the
		// event (if set absolute should not be set)
		"relative": types.Duration(),
		// if set, notification will be triggered at a given absolute time (if
		// set, relative should not be set)
		"absolute": dateType,
	}),
})

func (ce caldavEvent) ToValue() (out nu.Value) {
	rec := nu.Record{
		"uid":         nu.ToValue(ce.Uid),
		"name":        nu.ToValue(ce.Name),
		"location":    nu.ToValue(ce.Location),
		"description": nu.ToValue(ce.Description),
		"categories":  nu.ToValue(ce.Categories),
		"start":       nu.ToValue(ce.Start.ToValue()),
		"end":         nu.ToValue(ce.End.ToValue()),
	}

	rec["recurrence_set"] = nu.ToValue(nil)
	if ce.RecurrenceSet != nil {
		set := make(nu.Record)
		set["rule"] = nu.ToValue(ce.RecurrenceSet.Rule.String())

		exdates := make([]nu.Value, len(ce.RecurrenceSet.ExDates))
		for i, d := range ce.RecurrenceSet.ExDates {
			exdates[i] = d.ToValue()
		}
		set["exceptions"] = nu.ToValue(exdates)

		rdates := make([]nu.Value, len(ce.RecurrenceSet.RDates))
		for i, d := range ce.RecurrenceSet.RDates {
			rdates[i] = d.ToValue()
		}
		set["additional"] = nu.ToValue(rdates)

		rec["recurrence_set"] = nu.ToValue(set)
	}

	rec["recurrence_id"] = nu.ToValue(nil)
	if ce.RecurrenceId != nil {
		rec["recurrence_id"] = nu.ToValue(*ce.RecurrenceId)
	}

	rec["trigger"] = nu.ToValue(nil)
	if ce.Trigger != nil {
		trig := make(nu.Record)
		trig["relative"] = nu.ToValue(nil)
		trig["absolute"] = nu.ToValue(nil)
		if ce.Trigger.Relative != nil {
			trig["relative"] = nu.ToValue(*ce.Trigger.Relative)
		} else {
			trig["absolute"] = (*ce.Trigger.Absolute).ToValue()
		}
		rec["trigger"] = nu.ToValue(trig)
	}

	out = nu.ToValue(rec)
	return
}

func (out *caldavEvent) FromValue(in nu.Value) (err error) {
	rec, err := tryCast[nu.Record](in)
	if err != nil {
		return
	}

	if v, ok := rec["uid"]; ok {
		out.Uid, err = tryCast[string](v)
		if err != nil {
			return
		}
	}
	if v, ok := rec["name"]; ok {
		out.Name, err = tryCast[string](v)
		if err != nil {
			return
		}
	}
	if v, ok := rec["location"]; ok {
		out.Location, err = tryCast[string](v)
		if err != nil {
			return
		}
	}
	if v, ok := rec["categories"]; ok {
		out.Categories, err = tryCast[[]string](v)
		if err != nil {
			return
		}
	}
	if v, ok := rec["start"]; ok {
		dtstart := &out.Start
		err = dtstart.FromValue(v)
		if err != nil {
			return
		}
	}
	if v, ok := rec["end"]; ok {
		dtend := &out.End
		err = dtend.FromValue(v)
		if err != nil {
			return
		}
	}

	if v, ok := rec["recurrence_id"]; ok {
		out.RecurrenceId = &Datetime{}
		err = out.RecurrenceId.FromValue(v)
		if err != nil {
			return
		}
	} else if v, ok := rec["recurrence_set"]; ok {
		var set nu.Record
		set, err = tryCast[nu.Record](v)
		if err != nil {
			return
		}
		out.RecurrenceSet = &eventRecurrence{}

		var rule string
		rule, err = tryCast[string](set["rule"])
		if err != nil {
			return
		}
		var ropts *rrule.ROption
		ropts, err = rrule.StrToROption(rule)
		if err != nil {
			return
		}
		out.RecurrenceSet.Rule, err = rrule.NewRRule(*ropts)
		if err != nil {
			return
		}

		var exceptions []nu.Value
		exceptions, err = tryCast[[]nu.Value](set["exceptions"])
		if err != nil {
			return
		}
		out.RecurrenceSet.ExDates = make([]Datetime, len(exceptions))
		for i, ex := range exceptions {
			d := Datetime{}
			err = (&d).FromValue(ex)
			if err != nil {
				return
			}
			out.RecurrenceSet.ExDates[i] = d
		}

		var additional []nu.Value
		additional, err = tryCast[[]nu.Value](set["additional"])
		if err != nil {
			return
		}
		out.RecurrenceSet.RDates = make([]Datetime, len(additional))
		for i, ad := range additional {
			d := Datetime{}
			err = (&d).FromValue(ad)
			if err != nil {
				return
			}
			out.RecurrenceSet.RDates[i] = d
		}
	}

	if v, ok := rec["trigger"]; ok {
		var trig nu.Record
		trig, err = tryCast[nu.Record](v)
		if err != nil {
			return
		}

		out.Trigger = &eventTrigger{}

		if rel, ok := trig["relative"]; ok {
			var dur time.Duration
			dur, err = tryCast[time.Duration](rel)
			if err != nil {
				return
			}
			out.Trigger.Relative = &dur
		} else if abs, ok := trig["absolute"]; ok {
			parsed := &Datetime{}
			err = parsed.FromValue(abs)
			if err != nil {
				return
			}
			out.Trigger.Absolute = parsed
		}
	}
	return
}

var calendarType = types.Table(types.RecordDef{
	// path to the calendar
	"path": types.String(),
	// name of the calendar
	"name": types.String(),
	// description of the calendar
	"description": types.String(),
	// quota information
	"max_resource_size": types.Int(),
	// support information
	"supported_component_set": types.List(types.String()),
})

