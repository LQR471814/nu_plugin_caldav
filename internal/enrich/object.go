package enrich

import (
	"fmt"

	"github.com/ainvaltin/nu-plugin"
	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
)

// CalendarObject contains a VEVENT and fields related to it.
type CalendarObject struct {
	// ObjectPath is the event's calendar object path.
	ObjectPath *string
	// Main contains the main event for which the Overrides override.
	Main Event
	// Overrides contains all the recurrence overrides of the recurring event,
	// if the event is not recurring or there are no overrides, this list will
	// be empty/nil.
	Overrides []Event
}

func (cal CalendarObject) ToNu() (out nu.Value, err error) {
	record := make(nu.Record)

	record["object_path"] = nu.ToValue(cal.ObjectPath)
	record["main"], err = cal.Main.ToNu()

	overrides := make([]nu.Value, len(cal.Overrides))
	for i, ov := range cal.Overrides {
		overrides[i], err = ov.ToNu()
		if err != nil {
			return
		}
	}
	record["overrides"] = nu.ToValue(overrides)

	return
}

func CalendarObjectFromNu(val nu.Value) (out CalendarObject, err error) {
	record, ok := val.Value.(nu.Record)
	if !ok {
		err = fmt.Errorf("expected nu.Record got %T", val.Value)
		return
	}

	path := record["object_path"]
	switch path := path.Value.(type) {
	case string:
		out.ObjectPath = &path
	}

	out.Main, err = EventFromNu(record["main"])
	if err != nil {
		return
	}

	overrides := record["overrides"]
	switch overrides := overrides.Value.(type) {
	case []nu.Value:
		out.Overrides = make([]Event, len(overrides))
		for i, ov := range overrides {
			out.Overrides[i], err = EventFromNu(ov)
			if err != nil {
				return
			}
		}
	}

	return
}

func (cal CalendarObject) ToCaldav() (out *ical.Calendar) {
	out = ical.NewCalendar()

	version := ical.NewProp(ical.PropVersion)
	version.Value = "2.0"
	out.Props.Set(version)

	productId := ical.NewProp(ical.PropProductID)
	productId.Value = "-//LQR471814//Nushell CalDav Plugin 0.1//EN"
	out.Props.Set(productId)

	out.Children = append(out.Children, cal.Main.Underlying.Component)
	for _, ov := range cal.Overrides {
		out.Children = append(out.Children, ov.Underlying.Component)
	}
	return out
}

func CalendarFromCaldav(obj caldav.CalendarObject) (out CalendarObject) {
	out = CalendarObject{ObjectPath: &obj.Path}
	for _, component := range obj.Data.Children {
		if component.Name != ical.CompEvent {
			continue
		}
		event := EventFromCaldav(ical.Event{Component: component})
		prop := component.Props.Get(ical.PropRecurrenceID)
		if prop != nil {
			out.Overrides = append(out.Overrides, event)
			continue
		}
		out.Main = event
	}
	return out
}
