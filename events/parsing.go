package events

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/emersion/go-ical"
)

func (e Event) getString(name string) (string, bool) {
	prop := e.Props.Get(name)
	if prop == nil {
		return "", false
	}
	return prop.Value, true
}
func (e Event) setString(name, value string) {
	prop := ical.NewProp(name)
	prop.Value = value
	e.Props.Set(prop)
}

func (e Event) getInt(name string) (int, bool) {
	prop := e.Props.Get(name)
	if prop == nil {
		return 0, false
	}
	v, err := prop.Int()
	if err != nil {
		panic(err)
	}
	return v, true
}
func (e Event) setInt(name string, value int) {
	prop := ical.NewProp(name)
	prop.Value = fmt.Sprint(value)
	e.Props.Set(prop)
}

func (e Event) getURL(name string) (*url.URL, bool) {
	prop := e.Props.Get(name)
	if prop == nil {
		return nil, false
	}
	u, err := url.Parse(prop.Value)
	if err != nil {
		panic(err)
	}
	return u, true
}
func (e Event) setURL(name string, value *url.URL) {
	prop := ical.NewProp(name)
	prop.Value = value.String()
	e.Props.Set(prop)
}

func (e Event) getStringList(name string) ([]string, bool) {
	prop := e.Props.Get(name)
	if prop == nil {
		return nil, false
	}
	list, err := prop.TextList()
	if err != nil {
		panic(err)
	}
	return list, true
}
func (e Event) setStringList(name string, value []string) {
	prop := &ical.Prop{Name: name}
	prop.SetTextList(value)
	e.Props.Set(prop)
}

func (e Event) getDuration(name string) (time.Duration, bool) {
	prop := e.Props.Get(name)
	if prop == nil {
		return time.Duration(0), false
	}
	dur, err := prop.Duration()
	if err != nil {
		panic(err)
	}
	return dur, true
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

func (e Event) getDatetime(name string) (Datetime, bool) {
	prop := e.Props.Get(name)
	if prop == nil {
		return Datetime{}, false
	}
	tz, err := getTzidParam(prop, e.Timezone)
	if err != nil {
		panic(err)
	}
	d, err := parseDateText(prop.Value, tz)
	if err != nil {
		panic(err)
	}
	return d, true
}
func (e Event) setDatetime(name string, datetime Datetime) {
	prop := ical.NewProp(name)
	prop.Value = serializeDateText(datetime)
	if !datetime.Floating {
		setTzidParam(prop, datetime.Stamp.Location())
	}
	e.Props.Set(prop)
}

func (e Event) getDatetimeList(name string) ([]Datetime, bool) {
	prop := e.Props.Get(name)
	if prop == nil {
		return nil, false
	}
	tz, err := getTzidParam(prop, e.Timezone)
	if err != nil {
		panic(err)
	}

	segments := strings.Split(prop.Value, ",")
	dates := make([]Datetime, len(segments))

	for i, s := range segments {
		parsed, err := parseDateText(s, tz)
		if err != nil {
			panic(err)
		}
		if parsed.Floating && tz != nil {
			parsed.Floating = false
		}
		dates[i] = parsed
	}
	return dates, true
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
