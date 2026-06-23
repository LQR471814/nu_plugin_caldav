package props

import "time"

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
