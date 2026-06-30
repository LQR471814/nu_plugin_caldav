package events

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/emersion/go-ical"
)

func newTestEvent() Event {
	return Event{
		Event:    *ical.NewEvent(),
		Timezone: time.Local,
	}
}

func TestGetMissingPropertyReturnsSentinelError(t *testing.T) {
	event := newTestEvent()

	_, err := event.GetUID()
	if !errors.Is(err, ErrPropertyNotFound) {
		t.Fatalf("expected ErrPropertyNotFound, got %v", err)
	}
}

func TestGetInvalidPriorityReturnsParseError(t *testing.T) {
	event := newTestEvent()
	prop := ical.NewProp(ical.PropPriority)
	prop.Value = "high"
	event.Props.Set(prop)

	_, err := event.GetPriority()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), ical.PropPriority) {
		t.Fatalf("expected property name in error, got %v", err)
	}
}

func TestGetInvalidGeoReturnsParseError(t *testing.T) {
	event := newTestEvent()
	prop := ical.NewProp(ical.PropGeo)
	prop.Value = "47.1"
	event.Props.Set(prop)

	_, err := event.GetGeo()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), ical.PropGeo) {
		t.Fatalf("expected property name in error, got %v", err)
	}
}
