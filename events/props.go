package events

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/teambition/rrule-go"
)

// Uid is a globally unique identifier for this event.
//
// VEVENT Property: UID
func (e Event) GetUID() (string, bool) {
	return e.getString(ical.PropUID)
}
func (e Event) SetUID(uid string) {
	e.setString(ical.PropUID, uid)
}

// Summary is the human-friendly title for this event.
//
// VEVENT Property: SUMMARY
func (e Event) GetSummary() (string, bool) {
	return e.getString(ical.PropSummary)
}
func (e Event) SetSummary(summary string) {
	e.setString(ical.PropSummary, summary)
}

// Location is a string that represents the location of this event, can be in
// any format.
//
// VEVENT Property: LOCATION
func (e Event) GetLocation() (string, bool) {
	return e.getString(ical.PropLocation)
}
func (e Event) SetLocation(location *string) {
	if location == nil {
		e.Props.Del(ical.PropLocation)
		return
	}
	e.setString(ical.PropLocation, *location)
}

// Description is a human-friendly description for this event.
//
// VEVENT Property: DESCRIPTION
func (e Event) GetDescription() (string, bool) {
	return e.getString(ical.PropDescription)
}
func (e Event) SetDescription(description *string) {
	if description == nil {
		e.Props.Del(ical.PropDescription)
		return
	}
	e.setString(ical.PropDescription, *description)
}

// Categories represents tags or categories this event belongs to, strings do
// not need to be in any particular format.
//
// VEVENT Property: CATEGORIES
func (e Event) GetCategories() ([]string, bool) {
	return e.getStringList(ical.PropCategories)
}
func (e Event) SetCategories(categories []string) {
	e.Props.Del(ical.PropCategories)
	if len(categories) > 0 {
		e.setStringList(ical.PropCategories, categories)
	}
}

// Created defines when the event was created.
//
// VEVENT Property: CREATED
func (e Event) GetCreated() (Datetime, bool) {
	return e.getDatetime(ical.PropCreated)
}
func (e Event) SetCreated(createdAt *Datetime) {
	if createdAt == nil {
		e.Props.Del(ical.PropCreated)
		return
	}
	e.setDatetime(ical.PropCreated, *createdAt)
}

// LastModified defines when the event was last modified.
//
// VEVENT Property: LAST-MOD
func (e Event) GetLastModified() (Datetime, bool) {
	return e.getDatetime(ical.PropLastModified)
}
func (e Event) SetLastModified(datetime *Datetime) {
	if datetime == nil {
		e.Props.Del(ical.PropLastModified)
		return
	}
	e.setDatetime(ical.PropLastModified, *datetime)
}

// Class is the classification of the event (default: PUBLIC)
//
// VEVENT Property: CLASS
func (e Event) GetClass() (EventClass, bool) {
	str, ok := e.getString(ical.PropClass)
	return EventClass(str), ok
}
func (e Event) SetClass(class *EventClass) {
	if class == nil {
		e.Props.Del(ical.PropClass)
		return
	}
	e.setString(ical.PropClass, string(*class))
}

// Geo defines latitude and longitude for an event.
//
// VEVENT Property: GEO
func (e Event) GetGeo() (EventGeo, bool) {
	str, ok := e.getString(ical.PropGeo)
	if !ok {
		return EventGeo{}, false
	}
	segments := strings.Split(str, ";")
	if len(segments) != 2 {
		panic("invalid GEO property format")
	}
	lat, err := strconv.ParseFloat(segments[0], 64)
	if err != nil {
		panic(err)
	}
	long, err := strconv.ParseFloat(segments[1], 64)
	if err != nil {
		panic(err)
	}
	return EventGeo{
		Latitude:  lat,
		Longitude: long,
	}, ok
}
func (e Event) SetGeo(geo *EventGeo) {
	if geo == nil {
		e.Props.Del(ical.PropGeo)
		return
	}
	e.setString(
		ical.PropGeo,
		fmt.Sprintf("%f;%f", geo.Latitude, geo.Longitude),
	)
}

// Priority defines the priority of the calendar component.
//
//   - default: 0
//   - acceptable range: [0, 9]
//   - 1 is highest priority
//   - 9 is lowest priority
//
// If you are using an (HIGH, MEDIUM, LOW) priority system, you should follow
// the standard:
//
//   - [1, 4] specifies "HIGH" priority
//   - 5 specifies "MEDIUM" priority
//   - [6, 9] specifies "LOW" priority
//
// If you are using an (A1, A2 ... C3) priority system, you should follow the
// standard:
//
//   - A1 -> 1
//   - A2 -> 2
//   - ...
//   - C3 -> 9
//
// VEVENT Property: PRIORITY
func (e Event) GetPriority() (int, bool) {
	return e.getInt(ical.PropPriority)
}
func (e Event) SetPriority(priority *int) {
	if priority == nil {
		e.Props.Del(ical.PropPriority)
		return
	}
	e.setInt(ical.PropPriority, *priority)
}

// Sequence is a number that is incremented every time the organizer of the
// event makes a significant revision to the calendar component, it is
// effectively a "version".
//
// The attendee of the event includes this number when sending the event they
// are deciding on attending to the organizer to make it clear what version of
// the event they are okay with attending.
//
// VEVENT Property: SEQUENCE
func (e Event) GetSequence() (int, bool) {
	return e.getInt(ical.PropSequence)
}
func (e Event) SetSequence(sequence *int) {
	if sequence == nil {
		e.Props.Del(ical.PropSequence)
		return
	}
	e.setInt(ical.PropSequence, *sequence)
}

// Status defines the overall status or confirmation of the event.
//
// VEVENT Property: STATUS
func (e Event) GetStatus() (EventStatus, bool) {
	str, ok := e.getString(ical.PropStatus)
	return EventStatus(str), ok
}
func (e Event) SetStatus(class *EventStatus) {
	if class == nil {
		e.Props.Del(ical.PropStatus)
		return
	}
	e.setString(ical.PropStatus, string(*class))
}

// Transparency defines whether or not an event is transparent to busy time
// searches.
//
// VEVENT Property: TRANSP
func (e Event) GetTransparency() (EventTransparency, bool) {
	str, ok := e.getString(ical.PropTransparency)
	return EventTransparency(str), ok
}
func (e Event) SetTransparency(transparency *EventTransparency) {
	if transparency == nil {
		e.Props.Del(ical.PropTransparency)
		return
	}
	e.setString(ical.PropTransparency, string(*transparency))
}

// URL defines a URL associated with the event.
//
// VEVENT Property: URL
func (e Event) GetURL() (*url.URL, bool) {
	return e.getURL(ical.PropURL)
}
func (e Event) SetURL(url *url.URL) {
	if url == nil {
		e.Props.Del(ical.PropURL)
		return
	}
	e.setURL(ical.PropURL, url)
}

// Comment is a comment intended for the calendar user.
//
// VEVENT Property: COMMENT
func (e Event) GetComment() (string, bool) {
	return e.getString(ical.PropComment)
}
func (e Event) SetComment(comment *string) {
	if comment == nil {
		e.Props.Del(ical.PropComment)
		return
	}
	e.setString(ical.PropComment, *comment)
}

// Attach is an attachment, it is a URL by default.
//
// It can also be binary, but that functionality is rarely implemented.
//
// VEVENT Property: ATTACH
func (e Event) GetAttach() (*url.URL, bool) {
	return e.getURL(ical.PropAttach)
}
func (e Event) SetAttach(attachment *url.URL) {
	if attachment == nil {
		e.Props.Del(ical.PropAttach)
		return
	}
	e.setURL(ical.PropAttach, attachment)
}

// Attendee is a list of attendees to the event, each identified with a
// CAL-ADDRESS URL.
//
// VEVENT Property: ATTENDEE
func (e Event) GetAttendees() []*url.URL {
	// TODO: implement later
	return nil
}
func (e Event) SetAttendees(attendees []*url.URL) {
	// TODO: implement later
}

// Contact is some contact information associated with the event.
//
// VEVENT Property: CONTACT
func (e Event) GetContact() (string, bool) {
	return e.getString(ical.PropContact)
}
func (e Event) SetContact(contact *string) {
	if contact == nil {
		e.Props.Del(ical.PropContact)
		return
	}
	e.setString(ical.PropContact, *contact)
}

// Organizer is the organizer of the event, identified with a CAL-ADDRESS URL.
//
// VEVENT Property: ORGANIZER
func (e Event) GetOrganizer() (*url.URL, bool) {
	return e.getURL(ical.PropOrganizer)
}
func (e Event) SetOrganizer(organizer *url.URL) {
	if organizer == nil {
		e.Props.Del(ical.PropOrganizer)
		return
	}
	e.setURL(ical.PropOrganizer, organizer)
}

// TODO: implement
// rstatus
// resources

// Start defines when the event begins.
func (e Event) GetStart() (Datetime, bool) {
	return e.getDatetime(ical.PropDateTimeStart)
}
func (e Event) SetStart(start Datetime) {
	e.setDatetime(ical.PropDateTimeStart, start)
}

// End defines when the event ends.
func (e Event) GetEnd() (Datetime, bool) {
	return e.getDatetime(ical.PropDateTimeEnd)
}
func (e Event) SetEnd(start Datetime) {
	e.setDatetime(ical.PropDateTimeEnd, start)
}

// Duration defines the event's duration.
func (e Event) GetDuration() (time.Duration, bool) {
	// TODO: implement later
	return 0, false
}
func (e Event) SetDuration(duration time.Duration) {
	// TODO: implement later
}

func (e Event) GetRecurrenceRule() (out *rrule.RRule, ok bool) {
	// parse RRULE (does not support tzid)
	rruleProp := e.Props.Get(ical.PropRecurrenceRule)
	if rruleProp == nil {
		return
	}
	var ropts *rrule.ROption
	ropts, err := rrule.StrToROption(rruleProp.Value)
	if err != nil {
		panic(err)
	}
	if ropts == nil {
		panic("ropts is nil")
	}
	if ropts.Dtstart.Equal(time.Time{}) {
		dt, ok := e.GetStart()
		if !ok {
			panic("start time not defined")
		}
		// set default dtstart to original event's starting time
		ropts.Dtstart = dt.Stamp
	}
	var rule *rrule.RRule
	rule, err = rrule.NewRRule(*ropts)
	if err != nil {
		panic(err)
	}
	return rule, true
}
func (e Event) SetRecurrenceRule(rule *rrule.RRule) {
	if rule == nil {
		e.Props.Del(ical.PropRecurrenceRule)
		return
	}
	prop := ical.NewProp(ical.PropRecurrenceRule)
	prop.Value = rule.String()
	e.Props.Set(prop)
}

func (e Event) GetRecurrenceDates() ([]Datetime, bool) {
	return e.getDatetimeList(ical.PropRecurrenceDates)
}
func (e Event) SetRecurrenceDates(dates []Datetime) {
	e.Props.Del(ical.PropRecurrenceDates)
	if len(dates) > 0 {
		e.setDatetimeList(
			ical.PropRecurrenceDates,
			dates,
		)
	}
}

func (e Event) GetRecurrenceExceptionDates() ([]Datetime, bool) {
	return e.getDatetimeList(ical.PropExceptionDates)
}

func (e Event) SetRecurrenceExceptionDates(exceptions []Datetime) {
	e.Props.Del(ical.PropExceptionDates)
	if len(exceptions) > 0 {
		e.setDatetimeList(
			ical.PropExceptionDates,
			exceptions,
		)
	}
}

// Recurrence instance if set, defines this event as an override for a
// particular recurrence instance.
//
// The original event that it is being overriden is given by the event's uid.
//
// VEVENT Property: RECURID
func (e Event) GetRecurrenceInstance() (Datetime, bool) {
	return e.getDatetime(ical.PropRecurrenceID)
}
func (e Event) SetRecurrenceInstance(instance *Datetime) {
	if instance == nil {
		e.Props.Del(ical.PropRecurrenceID)
		return
	}
	e.setDatetime(ical.PropRecurrenceID, *instance)
}

// Trigger defines the notification trigger time (if any) for the event.
//
// VEVENT -> VALARM Property: TRIGGER
func (e Event) GetTrigger() (out EventTrigger, ok bool) {
	prop := e.Props.Get(ical.PropTrigger)
	if prop == nil {
		return
	}
	valueType := prop.Params.Get(ical.ParamValue)
	switch valueType {
	case "", "DURATION": // duration by default
		dur, err := prop.Duration()
		if err != nil {
			panic(err)
		}
		out.Relative = &dur
		switch prop.Params.Get(ical.ParamRelated) {
		case "", "START": // start by default
			out.RelativeTo = EVENT_TRIGGER_REL_START
		case "END":
			out.RelativeTo = EVENT_TRIGGER_REL_END
		}
	case "DATE-TIME":
		dt, err := prop.DateTime(e.Timezone)
		if err != nil {
			panic(err)
		}
		out.Absolute = &dt
	}
	ok = true
	return
}
func (e Event) SetTrigger(trigger *EventTrigger) {
	if trigger == nil {
		e.Props.Del(ical.PropTrigger)
		return
	}
	prop := &ical.Prop{Name: ical.PropTrigger}
	defer e.Props.Set(prop)
	if trigger.Relative != nil {
		prop.SetDuration(*trigger.Relative)
		prop.Params.Set(ical.ParamValue, "DURATION")
		switch trigger.RelativeTo {
		case EVENT_TRIGGER_REL_START:
			prop.Params.Set(ical.ParamRelated, "START")
		case EVENT_TRIGGER_REL_END:
			prop.Params.Set(ical.ParamRelated, "END")
		}
		return
	}
	if trigger.Absolute == nil {
		panic("event trigger must be set to either relative or absolute")
	}
	prop.SetDateTime(*trigger.Absolute)
}

type KeyValues struct {
	Key    string
	Values []ical.Prop
}

// GetOtherProps returns the remaining props on the event not covered by the
// standard ical spec
func (e Event) GetOtherProps() (out []KeyValues) {
	for k, v := range e.Props {
		switch k {
		case ical.PropUID,
			ical.PropSummary,
			ical.PropLocation,
			ical.PropDescription,
			ical.PropCategories,
			ical.PropCreated,
			ical.PropLastModified,
			ical.PropClass,
			ical.PropGeo,
			ical.PropPriority,
			ical.PropSequence,
			ical.PropStatus,
			ical.PropTransparency,
			ical.PropURL,
			ical.PropComment,
			ical.PropAttach,
			// TODO: implement
			// ical.PropAttendee,
			ical.PropContact,
			ical.PropOrganizer,
			// TODO: implement
			// rstatus
			// resources
			ical.PropDateTimeStart,
			ical.PropDateTimeEnd,
			// TODO: implement
			// ical.PropDuration,
			ical.PropRecurrenceRule,
			ical.PropRecurrenceDates,
			ical.PropExceptionDates,
			ical.PropRecurrenceID:
			continue
		default:
			out = append(out, KeyValues{
				Key:    k,
				Values: v,
			})
		}
	}
	return
}

func (e Event) SetOtherProp(prop *ical.Prop) {
	e.Props.Set(prop)
}

func (e Event) AddOtherProp(prop KeyValues) {
	for _, v := range prop.Values {
		e.Props.Add(&v)
	}
}
