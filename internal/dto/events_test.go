package dto

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/LQR471814/nu_plugin_caldav/internal/events"
	"github.com/emersion/go-ical"
)

func newTestEvent() events.Event {
	event := events.Event{
		Event:    *ical.NewEvent(),
		Timezone: time.Local,
	}
	event.SetUID("test")
	start := events.Datetime{Stamp: time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)}
	end := events.Datetime{Stamp: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)}
	event.SetStart(start)
	event.SetEnd(end)
	return event
}

func TestNewEventReturnsErrorForMissingRequiredProperty(t *testing.T) {
	event := newTestEvent()
	event.Props.Del(ical.PropUID)

	_, err := NewEvent(event)
	if !errors.Is(err, events.ErrPropertyNotFound) {
		t.Fatalf("expected ErrPropertyNotFound, got %v", err)
	}
}

func TestNewEventReturnsErrorForInvalidOptionalProperty(t *testing.T) {
	event := newTestEvent()
	prop := ical.NewProp(ical.PropGeo)
	prop.Value = "not-geo"
	event.Props.Set(prop)

	_, err := NewEvent(event)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), ical.PropGeo) {
		t.Fatalf("expected GEO in error, got %v", err)
	}
}

func TestEventApplyReturnsErrorForInvalidTrigger(t *testing.T) {
	dtoEvent := Event{
		Start:   events.Datetime{Stamp: time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC)},
		End:     events.Datetime{Stamp: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)},
		Trigger: &events.EventTrigger{},
	}

	err := dtoEvent.Apply(newTestEvent())
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "event trigger") {
		t.Fatalf("expected trigger error, got %v", err)
	}
}
