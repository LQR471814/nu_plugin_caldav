package events

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/emersion/go-ical"
)

var ErrPropertyNotFound = errors.New("event property not found")

func propertyNotFoundError(name string) error {
	return fmt.Errorf("%s: %w", name, ErrPropertyNotFound)
}

func (e Event) getString(name string) (string, error) {
	prop := e.Props.Get(name)
	if prop == nil {
		return "", propertyNotFoundError(name)
	}
	return prop.Value, nil
}
func (e Event) setString(name, value string) {
	prop := ical.NewProp(name)
	prop.Value = value
	e.Props.Set(prop)
}

func (e Event) getInt(name string) (int, error) {
	prop := e.Props.Get(name)
	if prop == nil {
		return 0, propertyNotFoundError(name)
	}
	v, err := prop.Int()
	if err != nil {
		return 0, fmt.Errorf("%s: parse integer %q: %w", name, prop.Value, err)
	}
	return v, nil
}
func (e Event) setInt(name string, value int) {
	prop := ical.NewProp(name)
	prop.Value = fmt.Sprint(value)
	e.Props.Set(prop)
}

func (e Event) getURL(name string) (*url.URL, error) {
	prop := e.Props.Get(name)
	if prop == nil {
		return nil, propertyNotFoundError(name)
	}
	u, err := url.Parse(prop.Value)
	if err != nil {
		return nil, fmt.Errorf("%s: parse URL %q: %w", name, prop.Value, err)
	}
	return u, nil
}
func (e Event) setURL(name string, value *url.URL) {
	prop := ical.NewProp(name)
	prop.Value = value.String()
	e.Props.Set(prop)
}

func (e Event) getStringList(name string) ([]string, error) {
	prop := e.Props.Get(name)
	if prop == nil {
		return nil, propertyNotFoundError(name)
	}
	list, err := prop.TextList()
	if err != nil {
		return nil, fmt.Errorf("%s: parse text list %q: %w", name, prop.Value, err)
	}
	return list, nil
}
func (e Event) setStringList(name string, value []string) {
	prop := &ical.Prop{Name: name}
	prop.SetTextList(value)
	e.Props.Set(prop)
}

func (e Event) getDuration(name string) (time.Duration, error) {
	prop := e.Props.Get(name)
	if prop == nil {
		return time.Duration(0), propertyNotFoundError(name)
	}
	dur, err := prop.Duration()
	if err != nil {
		return time.Duration(0), fmt.Errorf("%s: parse duration %q: %w", name, prop.Value, err)
	}
	return dur, nil
}
func (e Event) setDuration(name string, value time.Duration) {
	prop := ical.NewProp(name)
	prop.SetDuration(value)
	e.Props.Set(prop)
}

const (
	date_format         = "20060102"
	datetime_format     = "20060102T150405"
	datetime_utc_format = "20060102T150405Z"
)

func (e Event) getDatetime(name string) (Datetime, error) {
	prop := e.Props.Get(name)
	if prop == nil {
		return Datetime{}, propertyNotFoundError(name)
	}
	tz, err := getTzidParam(prop, e.Timezone)
	if err != nil {
		return Datetime{}, fmt.Errorf("%s: load timezone: %w", name, err)
	}
	d, err := parseDateText(prop.Value, tz)
	if err != nil {
		return Datetime{}, fmt.Errorf("%s: parse datetime %q: %w", name, prop.Value, err)
	}
	return d, nil
}
func (e Event) setDatetime(name string, datetime Datetime) {
	prop := ical.NewProp(name)
	prop.Value = serializeDateText(datetime)
	if !datetime.Floating {
		setTzidParam(prop, datetime.Stamp.Location())
	}
	e.Props.Set(prop)
}

func (e Event) getDatetimeList(name string) ([]Datetime, error) {
	prop := e.Props.Get(name)
	if prop == nil {
		return nil, propertyNotFoundError(name)
	}
	tz, err := getTzidParam(prop, e.Timezone)
	if err != nil {
		return nil, fmt.Errorf("%s: load timezone: %w", name, err)
	}

	segments := strings.Split(prop.Value, ",")
	dates := make([]Datetime, len(segments))

	for i, s := range segments {
		parsed, err := parseDateText(s, tz)
		if err != nil {
			return nil, fmt.Errorf("%s: parse datetime item %d %q: %w", name, i, s, err)
		}
		if parsed.Floating && tz != nil {
			parsed.Floating = false
		}
		dates[i] = parsed
	}
	return dates, nil
}
func (e Event) setDatetimeList(name string, datetimes []Datetime) {
	prop := ical.NewProp(name)

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
	e.Props.Set(prop)
}
