package dto

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/LQR471814/nu_plugin_caldav/events"
	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
	"github.com/teambition/rrule-go"
)

type PropValueDto struct {
	Value  string
	Params map[string][]string
}

type Event struct {
	Uid           *string
	Summary       *string
	Location      *string
	Description   *string
	Categories    []string
	DatetimeStamp *events.Datetime
	Created       *events.Datetime
	LastModified  *events.Datetime
	Class         *events.EventClass
	Geo           *events.EventGeo
	Priority      *int
	Sequence      *int
	Status        *events.EventStatus
	Transparency  *events.EventTransparency
	URL           *url.URL `name:"url"`
	Comment       *string
	Attach        *url.URL
	// TODO: implement attendees
	Contact   *string
	Organizer *url.URL
	Start     events.Datetime
	End       events.Datetime
	// TODO: implement duration
	RecurrenceRule           *rrule.RRule
	RecurrenceDates          []events.Datetime
	RecurrenceExceptionDates []events.Datetime
	RecurrenceInstance       *events.Datetime
	Trigger                  *events.EventTrigger
	Other                    map[string][]PropValueDto
}

func (e Event) String() string {
	var sb strings.Builder
	sb.WriteString("{")

	sb.WriteString("Start:")
	fmt.Fprint(&sb, e.Start.Stamp)
	if e.Start.AllDay {
		fmt.Fprint(&sb, ",AllDay")
	}
	if e.Start.Floating {
		fmt.Fprint(&sb, ",Floating")
	}
	sb.WriteString(" End:")
	fmt.Fprint(&sb, e.End.Stamp)
	if e.End.AllDay {
		fmt.Fprint(&sb, ",AllDay")
	}
	if e.End.Floating {
		fmt.Fprint(&sb, ",Floating")
	}

	if e.Uid != nil {
		sb.WriteString("Uid:")
		fmt.Fprint(&sb, *e.Uid)
	}
	if e.Summary != nil {
		sb.WriteString(" ")
		sb.WriteString("Summary:")
		fmt.Fprint(&sb, *e.Summary)
	}
	if e.Location != nil {
		sb.WriteString(" ")
		sb.WriteString("Location:")
		fmt.Fprint(&sb, *e.Location)
	}
	if e.Description != nil {
		sb.WriteString(" ")
		sb.WriteString("Description:")
		fmt.Fprint(&sb, *e.Description)
	}
	if e.Categories != nil {
		sb.WriteString(" ")
		sb.WriteString("Categories:")
		fmt.Fprint(&sb, e.Categories)
	}
	if e.DatetimeStamp != nil {
		sb.WriteString(" ")
		sb.WriteString("DatetimeStamp:")
		fmt.Fprint(&sb, *e.DatetimeStamp)
	}
	if e.Created != nil {
		sb.WriteString(" ")
		sb.WriteString("Created:")
		fmt.Fprint(&sb, *e.Created)
	}
	if e.LastModified != nil {
		sb.WriteString(" ")
		sb.WriteString("LastModified:")
		fmt.Fprint(&sb, *e.LastModified)
	}
	if e.Class != nil {
		sb.WriteString(" ")
		sb.WriteString("Class:")
		fmt.Fprint(&sb, *e.Class)
	}
	if e.Geo != nil {
		sb.WriteString(" ")
		sb.WriteString("Geo:")
		fmt.Fprint(&sb, *e.Geo)
	}
	if e.Priority != nil {
		sb.WriteString(" ")
		sb.WriteString("Priority:")
		fmt.Fprint(&sb, *e.Priority)
	}
	if e.Sequence != nil {
		sb.WriteString(" ")
		sb.WriteString("Sequence:")
		fmt.Fprint(&sb, *e.Sequence)
	}
	if e.Status != nil {
		sb.WriteString(" ")
		sb.WriteString("Status:")
		fmt.Fprint(&sb, *e.Status)
	}
	if e.Transparency != nil {
		sb.WriteString(" ")
		sb.WriteString("Transparency:")
		fmt.Fprint(&sb, *e.Transparency)
	}
	if e.URL != nil {
		sb.WriteString(" ")
		sb.WriteString("URL:")
		fmt.Fprint(&sb, *e.URL)
	}
	if e.Comment != nil {
		sb.WriteString(" ")
		sb.WriteString("Comment:")
		fmt.Fprint(&sb, *e.Comment)
	}
	if e.Attach != nil {
		sb.WriteString(" ")
		sb.WriteString("Attach:")
		fmt.Fprint(&sb, *e.Attach)
	}
	if e.Contact != nil {
		sb.WriteString(" ")
		sb.WriteString("Contact:")
		fmt.Fprint(&sb, *e.Contact)
	}
	if e.Organizer != nil {
		sb.WriteString(" ")
		sb.WriteString("Organizer:")
		fmt.Fprint(&sb, *e.Organizer)
	}
	if e.RecurrenceRule != nil {
		sb.WriteString(" ")
		sb.WriteString("RecurrenceRule:")
		fmt.Fprint(&sb, *e.RecurrenceRule)
	}
	if e.RecurrenceDates != nil {
		sb.WriteString(" ")
		sb.WriteString("RecurrenceDates:")
		fmt.Fprint(&sb, e.RecurrenceDates)
	}
	if e.RecurrenceExceptionDates != nil {
		sb.WriteString(" ")
		sb.WriteString("RecurrenceExceptionDates:")
		fmt.Fprint(&sb, e.RecurrenceExceptionDates)
	}
	if e.RecurrenceInstance != nil {
		sb.WriteString(" ")
		sb.WriteString("RecurrenceInstance:")
		fmt.Fprint(&sb, *e.RecurrenceInstance)
	}
	if e.Trigger != nil {
		sb.WriteString(" ")
		sb.WriteString("Trigger:")
		fmt.Fprint(&sb, *e.Trigger)
	}
	if e.Other != nil {
		sb.WriteString(" ")
		sb.WriteString("Other:")
		fmt.Fprint(&sb, e.Other)
	}
	sb.WriteString("}")

	return sb.String()
}

func NewEvent(e events.Event) (out Event) {
	uid, ok := e.GetUID()
	if !ok {
		panic("UID not defined in event")
	}
	out.Uid = &uid

	if res, ok := e.GetSummary(); ok {
		out.Summary = &res
	}
	if res, ok := e.GetLocation(); ok {
		out.Location = &res
	}
	if res, ok := e.GetDescription(); ok {
		out.Description = &res
	}
	if res, ok := e.GetCategories(); ok {
		out.Categories = res
	}
	if res, ok := e.GetDatetimeStamp(); ok {
		out.DatetimeStamp = &res
	}
	if res, ok := e.GetCreated(); ok {
		out.Created = &res
	}
	if res, ok := e.GetLastModified(); ok {
		out.LastModified = &res
	}
	if res, ok := e.GetClass(); ok {
		out.Class = &res
	}
	if res, ok := e.GetGeo(); ok {
		out.Geo = &res
	}
	if res, ok := e.GetPriority(); ok {
		out.Priority = &res
	}
	if res, ok := e.GetSequence(); ok {
		out.Sequence = &res
	}
	if res, ok := e.GetStatus(); ok {
		out.Status = &res
	}
	if res, ok := e.GetTransparency(); ok {
		out.Transparency = &res
	}
	if res, ok := e.GetURL(); ok {
		out.URL = res
	}
	if res, ok := e.GetComment(); ok {
		out.Comment = &res
	}
	if res, ok := e.GetAttach(); ok {
		out.Attach = res
	}
	if res, ok := e.GetContact(); ok {
		out.Contact = &res
	}
	if res, ok := e.GetOrganizer(); ok {
		out.Organizer = res
	}

	start, ok := e.GetStart()
	if !ok {
		panic("START not defined in event")
	}
	out.Start = start

	end, ok := e.GetEnd()
	if !ok {
		panic("END not defined in event")
	}
	out.End = end

	if res, ok := e.GetRecurrenceRule(); ok {
		out.RecurrenceRule = res
	}
	if res, ok := e.GetRecurrenceDates(); ok {
		out.RecurrenceDates = res
	}
	if res, ok := e.GetRecurrenceExceptionDates(); ok {
		out.RecurrenceExceptionDates = res
	}
	if res, ok := e.GetRecurrenceInstance(); ok {
		out.RecurrenceInstance = &res
	}
	if res, ok := e.GetTrigger(); ok {
		out.Trigger = &res
	}

	out.Other = make(map[string][]PropValueDto)
	for _, p := range e.GetOtherProps() {
		values := make([]PropValueDto, len(p.Values))
		for i, v := range p.Values {
			values[i] = PropValueDto{Value: v.Value, Params: v.Params}
		}
		out.Other[p.Key] = values
	}

	return
}

func (o Event) Apply(e events.Event) {
	if o.Uid != nil {
		e.SetUID(*o.Uid)
	}
	if o.Summary != nil {
		e.SetSummary(*o.Summary)
	}
	if o.Location != nil {
		e.SetLocation(o.Location)
	}
	if o.Description != nil {
		e.SetDescription(o.Description)
	}
	if o.Categories != nil {
		e.SetCategories(o.Categories)
	}
	if o.DatetimeStamp != nil {
		e.SetDatetimeStamp(o.DatetimeStamp)
	}
	if o.Created != nil {
		e.SetCreated(o.Created)
	}
	if o.LastModified != nil {
		e.SetLastModified(o.LastModified)
	}
	if o.Class != nil {
		e.SetClass(o.Class)
	}
	if o.Geo != nil {
		e.SetGeo(o.Geo)
	}
	if o.Priority != nil {
		e.SetPriority(o.Priority)
	}
	if o.Sequence != nil {
		e.SetSequence(o.Sequence)
	}
	if o.Status != nil {
		e.SetStatus(o.Status)
	}
	if o.Transparency != nil {
		e.SetTransparency(o.Transparency)
	}
	if o.URL != nil {
		e.SetURL(o.URL)
	}
	if o.Comment != nil {
		e.SetComment(o.Comment)
	}
	if o.Attach != nil {
		e.SetAttach(o.Attach)
	}
	if o.Contact != nil {
		e.SetContact(o.Contact)
	}
	if o.Organizer != nil {
		e.SetOrganizer(o.Organizer)
	}
	e.SetStart(o.Start)
	e.SetEnd(o.End)
	if o.RecurrenceRule != nil {
		e.SetRecurrenceRule(o.RecurrenceRule)
	}
	if o.RecurrenceDates != nil {
		e.SetRecurrenceDates(o.RecurrenceDates)
	}
	if o.RecurrenceExceptionDates != nil {
		e.SetRecurrenceExceptionDates(o.RecurrenceExceptionDates)
	}
	if o.RecurrenceInstance != nil {
		e.SetRecurrenceInstance(o.RecurrenceInstance)
	}
	if o.Trigger != nil {
		e.SetTrigger(o.Trigger)
	}
	for key, values := range o.Other {
		props := make([]ical.Prop, len(values))
		for i, v := range values {
			props[i] = ical.Prop{
				Name:   key,
				Value:  v.Value,
				Params: v.Params,
			}
		}
		e.AddOtherProp(events.KeyValues{
			Key:    key,
			Values: props,
		})
	}
}

// EventObject contains a VEVENT and fields related to it.
type EventObject struct {
	// ObjectPath is the event's calendar object path.
	ObjectPath *string
	// Main contains the main event for which the Overrides override.
	Main Event
	// Overrides contains all the recurrence overrides of the recurring event,
	// if the event is not recurring or there are no overrides, this list will
	// be empty/nil.
	Overrides []Event
}

func NewEventObject(obj caldav.CalendarObject) EventObject {
	dtoObj := EventObject{ObjectPath: &obj.Path}
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
			dtoObj.Overrides = append(dtoObj.Overrides, NewEvent(event))
			continue
		}
		dtoObj.Main = NewEvent(event)
	}
	return dtoObj
}

type EventObjectList []EventObject

func NewEventObjectList(objects []caldav.CalendarObject) EventObjectList {
	dtoObjects := make([]EventObject, len(objects))
	// each calendar object only ever stores one unique VEVENT object.
	//
	// exception:
	// if the VEVENT has recurrence overrides, the recurrence overrides will
	// come with the original VEVENT as separate VEVENT components.
	for i, obj := range objects {
		dtoObjects[i] = NewEventObject(obj)
	}
	return dtoObjects
}

type CalendarList []caldav.Calendar
