// All "*Prop" types will have nuconv conversions generated for them
package props

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/teambition/rrule-go"
)

// Transform[T] provides methods to read/write T from a generic ical.Prop
type Transform[T any] interface {
	Unmarshal(prop *ical.Prop) (out T, err error)
	// Marshal should return a new prop and not mutate prev.
	Marshal(prev *ical.Prop, value T) (prop *ical.Prop, err error)
}

// WrapAny[T] casts Transform[T] -> Transform[any]
type WrapAny[T any] struct {
	Inner Transform[T]
}

// AsAny casts Transform[T] -> Transform[any]
func AsAny[T any](transform Transform[T]) Transform[any] {
	return WrapAny[T]{Inner: transform}
}

func (a WrapAny[T]) Unmarshal(prop *ical.Prop) (out any, err error) {
	val, err := a.Inner.Unmarshal(prop)
	if err != nil {
		return
	}
	out = any(val)
	return
}

func (a WrapAny[T]) Marshal(prev *ical.Prop, value any) (prop *ical.Prop, err error) {
	return a.Inner.Marshal(prev, value.(T))
}

type TypedTextProp[T ~string] struct{}

func (text TypedTextProp[T]) Unmarshal(prop *ical.Prop) (out T, err error) {
	out = T(prop.Value)
	return
}

func (text TypedTextProp[T]) Marshal(prev *ical.Prop, value T) (prop *ical.Prop, err error) {
	prop = ical.NewProp(prev.Name)
	prop.Value = string(value)
	return
}

type TextProp struct{}

func (text TextProp) Unmarshal(prop *ical.Prop) (out string, err error) {
	out = prop.Value
	return
}

func (text TextProp) Marshal(prev *ical.Prop, value string) (prop *ical.Prop, err error) {
	prop = ical.NewProp(prev.Name)
	prop.Value = value
	return
}

type IntProp struct{}

func (int IntProp) Unmarshal(prop *ical.Prop) (out int, err error) {
	out, err = prop.Int()
	return
}

func (int IntProp) Marshal(prev *ical.Prop, value int) (prop *ical.Prop, err error) {
	prop = ical.NewProp(prev.Name)
	prop.Value = fmt.Sprint(value)
	return
}

type URLProp struct{}

func (URLProp) Unmarshal(prop *ical.Prop) (out *url.URL, err error) {
	return url.Parse(prop.Value)
}

func (URLProp) Marshal(prev *ical.Prop, value *url.URL) (prop *ical.Prop, err error) {
	prop = ical.NewProp(prev.Name)
	prop.Value = value.String()
	return
}

type TextListProp struct{}

func (TextListProp) Unmarshal(prop *ical.Prop) (out []string, err error) {
	return prop.TextList()
}

func (TextListProp) Marshal(prev *ical.Prop, value []string) (prop *ical.Prop, err error) {
	prop = ical.NewProp(prev.Name)
	prop.SetTextList(value)
	return
}

type DurationProp struct{}

func (DurationProp) Unmarshal(prop *ical.Prop) (out time.Duration, err error) {
	return prop.Duration()
}

func (DurationProp) Marshal(prev *ical.Prop, value time.Duration) (prop *ical.Prop, err error) {
	prop = ical.NewProp(prev.Name)
	prop.SetDuration(value)
	return
}

const (
	date_format         = "20060102"
	datetime_format     = "20060102T150405"
	datetime_utc_format = "20060102T150405Z"
)

type DatetimeProp struct {
	Timezone *time.Location
}

func (dt DatetimeProp) Unmarshal(prop *ical.Prop) (out Datetime, err error) {
	tz, err := getTzidParam(prop, dt.Timezone)
	if err != nil {
		return
	}

	return parseDateText(prop.Value, tz)
}

func (DatetimeProp) Marshal(prev *ical.Prop, datetime Datetime) (prop *ical.Prop, err error) {
	prop = ical.NewProp(prev.Name)
	prop.Value = serializeDateText(datetime)

	if !datetime.Floating {
		setTzidParam(prop, datetime.Stamp.Location())
	}

	return
}

type DatetimeListProp struct {
	Timezone *time.Location
}

func (dt DatetimeListProp) Unmarshal(prop *ical.Prop) (out []Datetime, err error) {
	tz, err := getTzidParam(prop, dt.Timezone)
	if err != nil {
		return
	}

	segments := strings.Split(prop.Value, ",")
	out = make([]Datetime, len(segments))

	for i, s := range segments {
		var parsed Datetime
		parsed, err = parseDateText(s, tz)
		if err != nil {
			return
		}

		if parsed.Floating && tz != nil {
			parsed.Floating = false
		}

		out[i] = parsed
	}

	return
}

func (DatetimeListProp) Marshal(prev *ical.Prop, datetimes []Datetime) (prop *ical.Prop, err error) {
	prop = ical.NewProp(prev.Name)

	var specifytz *time.Location
	for _, d := range datetimes {
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
	datestr := make([]string, len(datetimes))

	for i, d := range datetimes {
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
	return
}

type GeoProp struct{}

func (g GeoProp) Unmarshal(prop *ical.Prop) (out EventGeo, err error) {
	str, err := TextProp{}.Unmarshal(prop)
	if err != nil {
		err = fmt.Errorf("text prop unmarshal: %w", err)
		return
	}
	segments := strings.Split(str, "\\;")
	if len(segments) != 2 {
		err = fmt.Errorf("not exactly 2 segments separated by '\\;'")
		return
	}
	lat, err := strconv.ParseFloat(segments[0], 64)
	if err != nil {
		return
	}
	long, err := strconv.ParseFloat(segments[1], 64)
	if err != nil {
		return
	}
	out = EventGeo{
		Latitude:  lat,
		Longitude: long,
	}
	return
}

func (g GeoProp) Marshal(prev *ical.Prop, geo EventGeo) (prop *ical.Prop, err error) {
	prop = ical.NewProp(prev.Name)
	prop.Value = fmt.Sprintf("%f;%f", geo.Latitude, geo.Longitude)
	return
}

type RRuleProp struct {
	Start Datetime
}

func (g RRuleProp) Unmarshal(prop *ical.Prop) (out *rrule.RRule, err error) {
	var ropts *rrule.ROption

	ropts, err = rrule.StrToROption(prop.Value)
	if err != nil {
		err = fmt.Errorf("str to ropts: %w", err)
		return
	}

	if ropts == nil {
		err = fmt.Errorf("ropts is nil")
		return
	}

	if ropts.Dtstart.Equal(time.Time{}) {
		// set default dtstart to original event's starting time
		ropts.Dtstart = g.Start.Stamp
	}

	out, err = rrule.NewRRule(*ropts)
	if err != nil {
		err = fmt.Errorf("new rrule: %w", err)
		return
	}

	return
}

func (g RRuleProp) Marshal(prev *ical.Prop, value *rrule.RRule) (prop *ical.Prop, err error) {
	prop = ical.NewProp(prev.Name)
	prop.Value = value.String()
	return
}

type TriggerProp struct {
	Timezone *time.Location
}

func (t TriggerProp) Unmarshal(prop *ical.Prop) (out EventTrigger, err error) {
	valueType := prop.Params.Get(ical.ParamValue)
	switch valueType {
	case "", "DURATION": // duration by default
		var dur time.Duration

		dur, err = prop.Duration()
		if err != nil {
			return
		}

		out.Relative = &dur

		switch prop.Params.Get(ical.ParamRelated) {
		case "", "START": // start by default
			out.RelativeTo = EVENT_TRIGGER_REL_START
		case "END":
			out.RelativeTo = EVENT_TRIGGER_REL_END
		}
	case "DATE-TIME":
		var dt time.Time
		dt, err = prop.DateTime(t.Timezone)
		if err != nil {
			return
		}
		out.Absolute = &dt
	}
	return
}

func (t TriggerProp) Marshal(prev *ical.Prop, trigger EventTrigger) (prop *ical.Prop, err error) {
	prop = ical.NewProp(prev.Name)

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
		err = fmt.Errorf("event trigger must be set to either relative or absolute")
		return
	}

	prop.SetDateTime(*trigger.Absolute)
	return
}
