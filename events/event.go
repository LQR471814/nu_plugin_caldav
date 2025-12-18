package events

import (
	"time"

	"github.com/emersion/go-ical"
)

// Datetime defines a DATE-TIME property.
type Datetime struct {
	Stamp    time.Time
	AllDay   bool `default:"false"`
	Floating bool `default:"false"`
}

type EventTriggerRelative int

const (
	EVENT_TRIGGER_REL_START EventTriggerRelative = iota
	EVENT_TRIGGER_REL_END
)

// EventTrigger defines the time a notification for the event should go off.
type EventTrigger struct {
	Relative   *time.Duration
	RelativeTo EventTriggerRelative
	Absolute   *time.Time
}

type EventGeo struct {
	Latitude  float64
	Longitude float64
}

type EventClass string

const (
	EVENT_CLASS_PUBLIC       EventClass = "PUBLIC"
	EVENT_CLASS_PRIVATE      EventClass = "PRIVATE"
	EVENT_CLASS_CONFIDENTIAL EventClass = "CONFIDENTIAL"
)

type EventStatus string

const (
	EVENT_STATUS_TENTATIVE EventStatus = "TENTATIVE"
	EVENT_STATUS_CONFIRMED EventStatus = "CONFIRMED"
	EVENT_STATUS_CANCELLED EventStatus = "CANCELLED"
)

type EventTransparency string

const (
	EVENT_TRANSPARENCY_OPAQUE      EventTransparency = "OPAQUE"
	EVENT_TRANSPARENCY_TRANSPARENT EventTransparency = "TRANSPARENT"
)

type Event struct {
	Timezone *time.Location
	ical.Event
}

// EventObject is like EventContainer, but for caldav-facing code.
type EventObject struct {
	ObjectPath string `default:"\"\""`
	Main       Event
	Overrides  []Event
}

func (obj EventObject) ToCalendar() *ical.Calendar {
	cal := ical.NewCalendar()

	version := ical.NewProp(ical.PropVersion)
	version.Value = "2.0"
	cal.Props.Set(version)

	productId := ical.NewProp(ical.PropProductID)
	productId.Value = "-//LQR471814//Nushell CalDav Plugin 0.1//EN"
	cal.Props.Set(productId)

	cal.Children = append(cal.Children, obj.Main.Component)
	for _, ov := range obj.Overrides {
		cal.Children = append(cal.Children, ov.Event.Component)
	}
	return cal
}
