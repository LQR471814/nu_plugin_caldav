package enrich

import (
	"fmt"
	"net/url"
	"time"

	"github.com/LQR471814/nu_plugin_caldav/internal/enrich/dto"
	"github.com/LQR471814/nu_plugin_caldav/internal/enrich/props"
	"github.com/LQR471814/nu_plugin_caldav/internal/nuconv"
	"github.com/ainvaltin/nu-plugin"
	"github.com/emersion/go-ical"
	"github.com/teambition/rrule-go"
)

// Event is an enriched wrapper around ical.Event that allows for:
//
//   - Type-aware interpretation of ical.Event values
//   - Conversion to/from nu types
type Event struct {
	Underlying ical.Event
	Timezone   *time.Location
}

// TransformConstructor is a function that creates a props.Transform for a
// given event
type TransformConstructor[T any] = func(e Event) props.Transform[T]

// Field is a generic constructor of type-enriched manipulators
type Field[T any] struct {
	name      string
	transform TransformConstructor[T]
	isEmpty   func(v T) bool
}

func NewField[T any](
	name string,
	transform TransformConstructor[T],
	isEmpty func(v T) bool,
) Field[T] {
	return Field[T]{
		name:      name,
		transform: transform,
		isEmpty:   isEmpty,
	}
}

func (f Field[T]) Name() string {
	return f.name
}

func (f Field[T]) Transformer() TransformConstructor[T] {
	return f.transform
}

func (f Field[T]) IsEmpty() func(v T) bool {
	return f.isEmpty
}

func (f Field[T]) Get(e Event) (T, error) {
	prop := e.Underlying.Props.Get(f.name)
	transform := f.transform(e)
	return transform.Unmarshal(prop)
}

func (f Field[T]) Set(e Event, value T) error {
	if f.isEmpty != nil && f.isEmpty(value) {
		e.Underlying.Props.Del(f.name)
		return nil
	}
	prev := e.Underlying.Props.Get(f.name)
	if prev == nil {
		prev = ical.NewProp(f.name)
	}
	transform := f.transform(e)
	prop, err := transform.Marshal(prev, value)
	if err != nil {
		return err
	}
	e.Underlying.Props.Set(prop)
	return nil
}

func isZero[T comparable](v T) bool {
	var zero T
	return v == zero
}

func isNil[T any](v *T) bool {
	return v == nil
}

func isEmptySlice[T any](v []T) bool {
	return len(v) == 0
}

// wrapTransform creates a TransformConstructor which always returns a
// particular transformer.
func wrapTransform[T any](transform props.Transform[T]) TransformConstructor[T] {
	return func(e Event) props.Transform[T] { return transform }
}

func datetimeTransform(e Event) props.Transform[props.Datetime] {
	return props.DatetimeProp{Timezone: e.Timezone}
}

func datetimeListTransform(e Event) props.Transform[[]props.Datetime] {
	return props.DatetimeListProp{Timezone: e.Timezone}
}

var (
	textTransform = wrapTransform(props.TextProp{})
	intTransform  = wrapTransform(props.IntProp{})
	urlTransform  = wrapTransform(props.URLProp{})
)

var fieldRegistry = make(map[string]GenericField)

func register[T any](nuKey string, field Field[T]) Field[T] {
	var zero T
	fieldRegistry[nuKey] = GenericField{
		Zero:  zero,
		Inner: AsAny(field),
	}
	return field
}

func AllFields() map[string]GenericField {
	return fieldRegistry
}

var (
	// Uid is a globally unique identifier for this event.
	//
	// VEVENT Property: Uid
	Uid = register("uid", NewField(ical.PropUID, textTransform, nil))

	// Summary is the human-friendly title for this event.
	//
	// VEVENT Property: SUMMARY
	Summary = register("summary", NewField(ical.PropSummary, textTransform, nil))

	// Location is a string that represents the location of this event, can be in
	// any format.
	//
	// VEVENT Property: LOCATION
	Location = register("location", NewField(ical.PropLocation, textTransform, isZero[string]))

	// Description is a human-friendly description for this event.
	//
	// VEVENT Property: DESCRIPTION
	Description = register("description", NewField(ical.PropDescription, textTransform, isZero[string]))

	// Categories represents tags or categories this event belongs to, strings do
	// not need to be in any particular format.
	//
	// VEVENT Property: CATEGORIES
	Categories = register("categories", NewField(
		ical.PropCategories,
		wrapTransform(props.TextListProp{}),
		isEmptySlice[string],
	))

	// DatetimeStamp defines when the event is initially created (not in the store,
	// but on the client).
	//
	// VEVENT Property: DTSTAMP
	DatetimeStamp = register("datetime_stamp", NewField(ical.PropDateTimeStamp, datetimeTransform, isZero[props.Datetime]))

	// Created defines when the event was created in the store.
	//
	// VEVENT Property: CREATED
	Created = register("created", NewField(ical.PropCreated, datetimeTransform, isZero[props.Datetime]))

	// LastModified defines when the event was last modified in the store.
	//
	// VEVENT Property: LAST-MOD
	LastModified = register("last_modified", NewField(ical.PropLastModified, datetimeTransform, isZero[props.Datetime]))

	// Class is the classification of the event (default: PUBLIC)
	//
	// VEVENT Property: CLASS
	Class = register("class", NewField(
		ical.PropClass,
		wrapTransform(props.TypedTextProp[props.EventClass]{}),
		isZero[props.EventClass],
	))

	// Geo defines latitude and longitude for an event.
	//
	// VEVENT Property: GEO
	Geo = register("geo", NewField(
		ical.PropGeo,
		wrapTransform(props.GeoProp{}),
		isZero[props.EventGeo],
	))

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
	Priority = register("priority", NewField(ical.PropPriority, intTransform, nil))

	// Sequence is a number that is incremented every time the organizer of the
	// event makes a significant revision to the calendar component, it is
	// effectively a "version".
	//
	// The attendee of the event includes this number when sending the event they
	// are deciding on attending to the organizer to make it clear what version of
	// the event they are okay with attending.
	//
	// VEVENT Property: SEQUENCE
	Sequence = register("sequence", NewField(ical.PropSequence, intTransform, nil))

	// Status defines the overall status or confirmation of the event.
	//
	// VEVENT Property: STATUS
	Status = register("status", NewField(
		ical.PropStatus,
		wrapTransform(props.TypedTextProp[props.EventStatus]{}),
		isZero[props.EventStatus],
	))

	// Transparency defines whether or not an event is transparent to busy time
	// searches.
	//
	// VEVENT Property: TRANSP
	Transparency = register("transparency", NewField(
		ical.PropTransparency,
		wrapTransform(props.TypedTextProp[props.EventTransparency]{}),
		isZero[props.EventTransparency],
	))

	// URL defines a URL associated with the event.
	//
	// VEVENT Property: URL
	URL = register("url", NewField(ical.PropURL, urlTransform, isNil[url.URL]))

	// Comment is a comment intended for the calendar user.
	//
	// VEVENT Property: COMMENT
	Comment = register("comment", NewField(ical.PropComment, textTransform, isZero[string]))

	// Attach is an attachment, it is a URL by default.
	//
	// It can also be binary, but that functionality is rarely implemented.
	//
	// VEVENT Property: ATTACH
	Attach = register("attach", NewField(ical.PropAttach, urlTransform, isNil[url.URL]))

	// Attendee is a list of attendees to the event, each identified with a
	// CAL-ADDRESS URL.
	//
	// VEVENT Property: ATTENDEE
	// Attendee = TODO

	// Contact is some contact information associated with the event.
	//
	// VEVENT Property: CONTACT
	Contact = register("contact", NewField(ical.PropContact, textTransform, isZero[string]))

	// Organizer is the organizer of the event, identified with a CAL-ADDRESS URL.
	//
	// VEVENT Property: ORGANIZER
	Organizer = register("organizer", NewField(ical.PropOrganizer, urlTransform, isNil[url.URL]))

	// TODO: implement
	// rstatus
	// resources

	// Start defines when the event begins.
	Start = register("start", NewField(ical.PropDateTimeStart, datetimeTransform, nil))

	// End defines when the event ends.
	End = register("end", NewField(ical.PropDateTimeEnd, datetimeTransform, nil))

	// Duration defines the event's duration.
	Duration = register("duration", NewField(
		ical.PropDuration,
		wrapTransform(props.DurationProp{}),
		nil,
	))

	// Recurrence instance if set, defines this event as an override for a
	// particular recurrence instance.
	//
	// The original event that it is being overriden is given by the event's uid.
	//
	// VEVENT Property: RECURID
	RecurrenceRule = register("recurrence_rule", NewField(
		ical.PropRecurrenceRule,
		func(e Event) props.Transform[*rrule.RRule] {
			start, _ := Start.Get(e)
			return props.RRuleProp{Start: start}
		},
		isNil[rrule.RRule],
	))

	RecurrenceDates = register(
		"recurrence_dates",
		NewField(ical.PropRecurrenceDates, datetimeListTransform, isEmptySlice[props.Datetime]),
	)

	RecurrenceExceptionDates = register(
		"recurrence_exception_dates",
		NewField(ical.PropExceptionDates, datetimeListTransform, isEmptySlice[props.Datetime]),
	)

	RecurrenceInstance = register(
		"recurrence_instance",
		NewField(ical.PropRecurrenceID, datetimeTransform, isZero[props.Datetime]),
	)

	// Trigger defines the notification trigger time (if any) for the event.
	//
	// VEVENT -> VALARM Property: TRIGGER
	Trigger = register("trigger", NewField(
		ical.PropTrigger,
		func(e Event) props.Transform[props.EventTrigger] {
			return props.TriggerProp{Timezone: e.Timezone}
		},
		isZero[props.EventTrigger],
	))
)

type GenericField struct {
	Zero  any
	Inner Field[any]
}

// AsAny casts a Field[T] -> Field[any]
func AsAny[T any](field Field[T]) Field[any] {
	return NewField(
		field.name,
		func(e Event) props.Transform[any] {
			transform := field.transform(e)
			return props.AsAny(transform)
		},
		func(v any) bool {
			return field.isEmpty(v.(T))
		},
	)
}

func (e Event) setOtherFields(out nu.Record) (err error) {
	for key, value := range e.Underlying.Props {
		_, ok := fieldRegistry[key]
		if ok {
			continue
		}
		values := make([]dto.PropValue, len(value))
		for i, v := range value {
			values[i] = dto.PropValue{Value: v.Value, Params: v.Params}
		}
		var nuValue nu.Value
		nuValue, err = nuconv.Marshal(values)
		if err != nil {
			return
		}
		out[key] = nuValue
	}
	return
}

func (e Event) ToNu() (out nu.Value, err error) {
	record := make(nu.Record)
	for nuKey, field := range fieldRegistry {
		var val any
		val, err = field.Inner.Get(e)
		if err != nil {
			return
		}
		record[nuKey], err = nuconv.Marshal(val)
		if err != nil {
			return
		}
	}
	err = e.setOtherFields(record)
	if err != nil {
		return
	}
	out = nu.ToValue(record)
	return
}

func (e Event) getOtherFields(rec nu.Record) (err error) {
	for nuKey, value := range rec {
		_, ok := fieldRegistry[nuKey]
		if !ok {
			continue
		}
		var propVal dto.PropValueList
		propVal, err = nuconv.PropValueListFromNu(value)
		if err != nil {
			return
		}
		for _, val := range propVal {
			prop := ical.NewProp(nuKey)
			prop.Value = val.Value
			prop.Params = ical.Params(val.Params)
			e.Underlying.Props.Add(prop)
		}
	}
	return
}

func EventFromNu(val nu.Value) (out Event, err error) {
	out = Event{
		Underlying: *ical.NewEvent(),
		Timezone:   time.Local,
	}
	record, ok := val.Value.(nu.Record)
	if !ok {
		err = fmt.Errorf("expected nu.Record got %T", val.Value)
		return
	}
	for nuKey, field := range fieldRegistry {
		value := field.Zero
		err = nuconv.Unmarshal(record[nuKey], &value)
		if err != nil {
			return
		}
		err = field.Inner.Set(out, value)
		if err != nil {
			return
		}
	}
	out.getOtherFields(record)
	return
}

func (e Event) Apply(other Event) (err error) {
	for key, value := range e.Underlying.Props {
		field, ok := fieldRegistry[key]

		if !ok {
			for i := range value {
				e.Underlying.Props.Add(&value[i])
			}
			continue
		}

		var value any
		value, err = field.Inner.Get(other)
		if err != nil {
			return
		}
		err = field.Inner.Set(e, value)
		if err != nil {
			return
		}
	}
	return
}

func (e Event) ToCaldav() ical.Event {
	return e.Underlying
}

func EventFromCaldav(ev ical.Event) Event {
	return Event{
		Timezone:   time.Local,
		Underlying: ev,
	}
}
