package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/teambition/rrule-go"
)

type eventTrigger struct {
	Relative *time.Duration
	Absolute *Datetime
}

type eventRecurrence struct {
	Rule    *rrule.RRule
	ExDates []Datetime
	RDates  []Datetime
}

type caldavEvent struct {
	Uid           string
	Name          string
	Location      string
	Description   string
	Categories    []string
	Start, End    Datetime
	RecurrenceSet *eventRecurrence
	RecurrenceId  *Datetime
	Trigger       *eventTrigger
}

func parseEvent(e ical.Event) (event caldavEvent, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parse event: %w", err)
		}
	}()

	(&event).parseUid(e)
	(&event).parseName(e)
	(&event).parseLocation(e)
	(&event).parseDescription(e)

	err = (&event).parseStart(e)
	if err != nil {
		return
	}
	err = (&event).parseEnd(e)
	if err != nil {
		return
	}

	err = (&event).parseCategories(e)
	if err != nil {
		return
	}
	err = (&event).parseRecurrence(e, event.Start.Stamp)
	if err != nil {
		return
	}
	err = (&event).parseTrigger(e)
	if err != nil {
		return
	}

	if event.Uid == "" {
		err = fmt.Errorf("uid is nil")
		return
	}

	return
}

func (out *caldavEvent) parseUid(e ical.Event) {
	uidProp := e.Props.Get(ical.PropUID)
	if uidProp == nil {
		return
	}
	out.Uid = uidProp.Value
}

func (out *caldavEvent) parseName(e ical.Event) {
	nameProp := e.Props.Get(ical.PropSummary)
	if nameProp == nil {
		return
	}
	out.Name = nameProp.Value
}

func (out *caldavEvent) parseDescription(e ical.Event) {
	descProp := e.Props.Get(ical.PropDescription)
	if descProp == nil {
		return
	}
	out.Description = descProp.Value
}

func (out *caldavEvent) parseLocation(e ical.Event) {
	locProp := e.Props.Get(ical.PropLocation)
	if locProp == nil {
		return
	}
	out.Location = locProp.Value
}

func (out *caldavEvent) parseCategories(e ical.Event) error {
	catProp := e.Props.Get(ical.PropCategories)
	if catProp == nil {
		return nil
	}
	categories, err := catProp.TextList()
	if err != nil {
		return err
	}
	out.Categories = categories
	return nil
}

func (out *caldavEvent) parseStart(e ical.Event) error {
	start, err := getPropDate(e.Props.Get(ical.PropDateTimeStart))
	if err != nil {
		return err
	}
	out.Start = start
	return nil
}

func (out *caldavEvent) parseEnd(e ical.Event) error {
	end, err := getPropDate(e.Props.Get(ical.PropDateTimeEnd))
	if err != nil {
		return err
	}
	out.End = end
	return nil
}

func (out *caldavEvent) parseRecurrence(e ical.Event, start time.Time) (err error) {
	// parse Recurrence-ID
	recurIdProp := e.Props.Get(ical.PropRecurrenceID)
	if recurIdProp != nil && recurIdProp.Value != "" {
		var id Datetime
		id, err = getPropDate(recurIdProp)
		if err != nil {
			return
		}
		out.RecurrenceId = &id
		return
	}

	// parse RRULE (does not support tzid)
	rruleProp := e.Props.Get(ical.PropRecurrenceRule)
	if rruleProp == nil {
		return
	}
	var ropts *rrule.ROption
	ropts, err = rrule.StrToROption(rruleProp.Value)
	if err != nil {
		return
	}
	if ropts == nil {
		err = fmt.Errorf("ropts is nil")
		return
	}
	if ropts.Dtstart.Equal(time.Time{}) {
		// set default dtstart to original event's starting time
		ropts.Dtstart = start
	}
	var rule *rrule.RRule
	rule, err = rrule.NewRRule(*ropts)
	if err != nil {
		return
	}
	out.RecurrenceSet = &eventRecurrence{
		Rule: rule,
	}

	// parse RDATE
	rdateProp := e.Props.Get(ical.PropRecurrenceDates)
	if rdateProp != nil {
		var dates []Datetime
		dates, err = getPropDateList(rdateProp)
		if err != nil {
			return
		}
		out.RecurrenceSet.RDates = dates
	}

	// PARSE EXDATE
	exdateProp := e.Props.Get(ical.PropExceptionDates)
	if exdateProp != nil {
		var dates []Datetime
		dates, err = getPropDateList(exdateProp)
		if err != nil {
			return
		}
		out.RecurrenceSet.ExDates = dates
	}

	return
}

func (out *caldavEvent) parseTrigger(e ical.Event) (err error) {
	triggerProp := e.Props.Get(ical.PropTrigger)
	if triggerProp == nil {
		return
	}
	dur, err := triggerProp.Duration()
	if err == nil {
		out.Trigger.Relative = &dur
		return
	}
	stamp, err := getPropDate(triggerProp)
	if err == nil {
		out.Trigger.Absolute = &stamp
		return
	}
	return
}

func (out caldavEvent) Write(e ical.Event) {
	e.Props.Set(&ical.Prop{
		Name:  ical.PropUID,
		Value: out.Uid,
	})
	e.Props.Set(&ical.Prop{
		Name:  ical.PropSummary,
		Value: out.Name,
	})
	e.Props.Set(&ical.Prop{
		Name:  ical.PropDescription,
		Value: out.Description,
	})
	e.Props.Set(&ical.Prop{
		Name:  ical.PropLocation,
		Value: out.Location,
	})
	cats := &ical.Prop{
		Name: ical.PropCategories,
	}
	cats.SetTextList(out.Categories)
	e.Props.Set(cats)

	dtstart := &ical.Prop{
		Name: ical.PropDateTimeStart,
	}
	setPropDate(dtstart, out.Start)
	e.Props.Set(dtstart)

	dtend := &ical.Prop{
		Name: ical.PropDateTimeEnd,
	}
	setPropDate(dtend, out.End)
	e.Props.Set(dtend)

	if out.RecurrenceId != nil {
		prop := &ical.Prop{
			Name: ical.PropRecurrenceID,
		}
		setPropDate(prop, *out.RecurrenceId)
		e.Props.Set(prop)
	}

	if out.RecurrenceSet != nil {
		e.Props.Set(&ical.Prop{
			Name:  ical.PropRecurrenceRule,
			Value: out.RecurrenceSet.Rule.String(),
		})

		exdates := &ical.Prop{
			Name: ical.PropExceptionDates,
		}
		setPropDateList(exdates, out.RecurrenceSet.ExDates)
		e.Props.Set(exdates)

		rdates := &ical.Prop{
			Name: ical.PropRecurrenceDates,
		}
		setPropDateList(rdates, out.RecurrenceSet.RDates)
		e.Props.Set(rdates)
	}

	if out.Trigger != nil {
		trigProp := &ical.Prop{
			Name: ical.PropTrigger,
		}
		if out.Trigger.Relative != nil {
			trigProp.SetDuration(*out.Trigger.Relative)
		} else {
			setPropDate(trigProp, *out.Trigger.Absolute)
		}
		e.Props.Set(trigProp)
	}
}

const (
	date_format         = "20060102"
	datetime_format     = "20060102T150405"
	datetime_utc_format = "20060102T150405Z"
)

func setPropDate(prop *ical.Prop, date Datetime) {
	if !date.Floating {
		setTzidParam(prop, date.Stamp.Location())
	}
	prop.SetText(serializeDateText(date))
}

func setPropDateList(prop *ical.Prop, dates []Datetime) {
	var specifytz *time.Location
	for _, d := range dates {
		if d.Floating || d.Stamp.Location() == time.UTC {
			continue
		}
		specifytz = d.Stamp.Location()
		break
	}

	if specifytz != nil {
		setTzidParam(prop, specifytz)
	}

	allday := true
	datestr := make([]string, len(dates))
	for i, d := range dates {
		if !d.AllDay {
			allday = false
		}
		datestr[i] = serializeDateText(d)
	}
	if allday {
		prop.Params.Set(ical.ParamValue, "DATE")
	} else {
		prop.Params.Set(ical.ParamValue, "DATE-TIME")
	}
	prop.Value = strings.Join(datestr, ",")
}

func getPropDateList(prop *ical.Prop) (dates []Datetime, err error) {
	tz, err := getTzidParam(prop)
	if err != nil {
		return
	}
	datestr := strings.Split(prop.Value, ",")

	dates = make([]Datetime, len(datestr))
	for i, s := range datestr {
		var parsed Datetime
		parsed, err = parseDateText(s, tz)
		if err != nil {
			return
		}
		if parsed.Floating && tz != nil {
			parsed.Floating = false
		}
		dates[i] = parsed
	}
	return
}

func getPropDate(prop *ical.Prop) (date Datetime, err error) {
	tz, err := getTzidParam(prop)
	if err != nil {
		return
	}
	date, err = parseDateText(prop.Value, tz)
	if err != nil {
		return
	}
	if date.Floating && tz != nil {
		date.Floating = false
	}
	return
}

func getTzidParam(prop *ical.Prop) (tz *time.Location, err error) {
	tz = time.Local
	tzid := prop.Params.Get(ical.PropTimezoneID)
	if tzid == "" {
		return
	}
	tz, err = time.LoadLocation(tzid)
	if err != nil {
		return
	}
	return
}

func setTzidParam(prop *ical.Prop, tz *time.Location) {
	if tz == nil {
		prop.Params.Del(ical.PropTimezoneID)
		return
	}
	prop.Params.Set(ical.PropTimezoneID, tz.String())
}

func parseDateText(s string, tz *time.Location) (d Datetime, err error) {
	if tz == nil {
		tz = time.Local
	}
	var layout string
	switch len(s) {
	case len(date_format):
		layout = date_format
		d.AllDay = true
		d.Floating = true
	case len(datetime_format):
		layout = datetime_format
		d.Floating = true
	case len(datetime_utc_format):
		layout = datetime_utc_format
	}
	d.Stamp, err = time.ParseInLocation(layout, s, tz)
	if err != nil {
		return
	}
	return
}

func serializeDateText(d Datetime) string {
	if d.AllDay {
		return d.Stamp.Format(date_format)
	}
	if d.Stamp.Location() == time.UTC {
		return d.Stamp.Format(datetime_utc_format)
	}
	return d.Stamp.Format(datetime_format)
}
