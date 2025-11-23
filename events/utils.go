package events

import (
	"time"

	"github.com/emersion/go-ical"
	"github.com/thlib/go-timezone-local/tzlocal"
)

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

func getTzidParam(prop *ical.Prop, defaultzone *time.Location) (tz *time.Location, err error) {
	if prop.Params == nil {
		return
	}
	tz = defaultzone
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

var tzname string

func init() {
	var err error
	tzname, err = tzlocal.RuntimeTZ()
	if err != nil {
		panic(err)
	}
}

func setTzidParam(prop *ical.Prop, tz *time.Location) {
	if prop.Params == nil {
		prop.Params = make(ical.Params)
	}
	if tz == nil {
		prop.Params.Del(ical.PropTimezoneID)
		return
	}
	name := tz.String()
	if name == "Local" {
		name = tzname
	}
	prop.Params.Set(ical.PropTimezoneID, name)
}
