package dto

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/LQR471814/nu_plugin_caldav/internal/events"
	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
	"github.com/teambition/rrule-go"
)

// rrule.RRule but with gob encoding support
type RRule struct {
	*rrule.RRule
}

func (r RRule) GobEncode() ([]byte, error) {
	return []byte(r.String()), nil
}

func (r *RRule) GobDecode(s []byte) error {
	rule, err := rrule.StrToRRule(string(s))
	if err != nil {
		return err
	}
	*r = RRule{RRule: rule}
	return nil
}

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
	RecurrenceRule           RRule
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
	if e.RecurrenceRule.RRule != nil {
		sb.WriteString(" ")
		sb.WriteString("RecurrenceRule:")
		fmt.Fprint(&sb, *e.RecurrenceRule.RRule)
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

func optionalEventProp[T any](value T, err error) (T, bool, error) {
	if err == nil {
		return value, true, nil
	}
	if errors.Is(err, events.ErrPropertyNotFound) {
		var zero T
		return zero, false, nil
	}
	var zero T
	return zero, false, err
}

func requireEventProp[T any](value T, err error) (T, error) {
	if err != nil {
		var zero T
		return zero, fmt.Errorf("required event property: %w", err)
	}
	return value, nil
}

func NewEvent(e events.Event) (out Event, err error) {
	uid, err := requireEventProp(e.GetUID())
	if err != nil {
		return
	}
	out.Uid = &uid

	if res, ok, err := optionalEventProp(e.GetSummary()); err != nil {
		return out, err
	} else if ok {
		out.Summary = &res
	}
	if res, ok, err := optionalEventProp(e.GetLocation()); err != nil {
		return out, err
	} else if ok {
		out.Location = &res
	}
	if res, ok, err := optionalEventProp(e.GetDescription()); err != nil {
		return out, err
	} else if ok {
		out.Description = &res
	}
	if res, ok, err := optionalEventProp(e.GetCategories()); err != nil {
		return out, err
	} else if ok {
		out.Categories = res
	}
	if res, ok, err := optionalEventProp(e.GetDatetimeStamp()); err != nil {
		return out, err
	} else if ok {
		out.DatetimeStamp = &res
	}
	if res, ok, err := optionalEventProp(e.GetCreated()); err != nil {
		return out, err
	} else if ok {
		out.Created = &res
	}
	if res, ok, err := optionalEventProp(e.GetLastModified()); err != nil {
		return out, err
	} else if ok {
		out.LastModified = &res
	}
	if res, ok, err := optionalEventProp(e.GetClass()); err != nil {
		return out, err
	} else if ok {
		out.Class = &res
	}
	if res, ok, err := optionalEventProp(e.GetGeo()); err != nil {
		return out, err
	} else if ok {
		out.Geo = &res
	}
	if res, ok, err := optionalEventProp(e.GetPriority()); err != nil {
		return out, err
	} else if ok {
		out.Priority = &res
	}
	if res, ok, err := optionalEventProp(e.GetSequence()); err != nil {
		return out, err
	} else if ok {
		out.Sequence = &res
	}
	if res, ok, err := optionalEventProp(e.GetStatus()); err != nil {
		return out, err
	} else if ok {
		out.Status = &res
	}
	if res, ok, err := optionalEventProp(e.GetTransparency()); err != nil {
		return out, err
	} else if ok {
		out.Transparency = &res
	}
	if res, ok, err := optionalEventProp(e.GetURL()); err != nil {
		return out, err
	} else if ok {
		out.URL = res
	}
	if res, ok, err := optionalEventProp(e.GetComment()); err != nil {
		return out, err
	} else if ok {
		out.Comment = &res
	}
	if res, ok, err := optionalEventProp(e.GetAttach()); err != nil {
		return out, err
	} else if ok {
		out.Attach = res
	}
	if res, ok, err := optionalEventProp(e.GetContact()); err != nil {
		return out, err
	} else if ok {
		out.Contact = &res
	}
	if res, ok, err := optionalEventProp(e.GetOrganizer()); err != nil {
		return out, err
	} else if ok {
		out.Organizer = res
	}

	start, err := requireEventProp(e.GetStart())
	if err != nil {
		return
	}
	out.Start = start

	end, err := requireEventProp(e.GetEnd())
	if err != nil {
		return
	}
	out.End = end

	if res, ok, err := optionalEventProp(e.GetRecurrenceRule()); err != nil {
		return out, err
	} else if ok {
		out.RecurrenceRule.RRule = res
	}
	if res, ok, err := optionalEventProp(e.GetRecurrenceDates()); err != nil {
		return out, err
	} else if ok {
		out.RecurrenceDates = res
	}
	if res, ok, err := optionalEventProp(e.GetRecurrenceExceptionDates()); err != nil {
		return out, err
	} else if ok {
		out.RecurrenceExceptionDates = res
	}
	if res, ok, err := optionalEventProp(e.GetRecurrenceInstance()); err != nil {
		return out, err
	} else if ok {
		out.RecurrenceInstance = &res
	}
	if res, ok, err := optionalEventProp(e.GetTrigger()); err != nil {
		return out, err
	} else if ok {
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

func (o Event) Apply(e events.Event) error {
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
	if o.RecurrenceRule.RRule != nil {
		e.SetRecurrenceRule(o.RecurrenceRule.RRule)
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
		if o.Trigger.Relative == nil && o.Trigger.Absolute == nil {
			return fmt.Errorf("event trigger must be set to either relative or absolute")
		}
		if o.Trigger.Relative != nil {
			switch o.Trigger.RelativeTo {
			case events.EVENT_TRIGGER_REL_START, events.EVENT_TRIGGER_REL_END:
			default:
				return fmt.Errorf("unsupported event trigger relative target: %d", o.Trigger.RelativeTo)
			}
		}
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
	return nil
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

func NewEventObject(obj caldav.CalendarObject) (EventObject, error) {
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
			dtoEvent, err := NewEvent(event)
			if err != nil {
				return dtoObj, fmt.Errorf("convert recurrence override %q: %w", obj.Path, err)
			}
			dtoObj.Overrides = append(dtoObj.Overrides, dtoEvent)
			continue
		}
		dtoEvent, err := NewEvent(event)
		if err != nil {
			return dtoObj, fmt.Errorf("convert main event %q: %w", obj.Path, err)
		}
		dtoObj.Main = dtoEvent
	}
	return dtoObj, nil
}

type EventObjectList []EventObject

func NewEventObjectList(objects []caldav.CalendarObject) (EventObjectList, error) {
	dtoObjects := make([]EventObject, len(objects))
	// each calendar object only ever stores one unique VEVENT object.
	//
	// exception:
	// if the VEVENT has recurrence overrides, the recurrence overrides will
	// come with the original VEVENT as separate VEVENT components.
	for i, obj := range objects {
		dtoObj, err := NewEventObject(obj)
		if err != nil {
			return dtoObjects, err
		}
		dtoObjects[i] = dtoObj
	}
	return dtoObjects, nil
}

type CalendarList []caldav.Calendar
