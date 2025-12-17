package main

import (
	"math/rand"
	"testing"
	"time"

	"github.com/LQR471814/nu_plugin_caldav/events"
	"github.com/LQR471814/nu_plugin_caldav/internal/dto"
)

func abs[T int64 | int32 | int16 | int8](v T) T {
	if v < 0 {
		return -v
	}
	return v
}

func FuzzConvertToTimeline(f *testing.F) {
	// inputs:
	// - overall_start: int
	// - overall_end: int
	// - overall_end >= overall_start
	// - list[Event]: (length: int)
	//   - start: int
	//   - end: int
	//   - end >= start
	//   - start >= overall_start

	// outputs:
	// - list[TimeSegment]:
	//   - start: int
	//   - dur: int
	//   - start+dur <= overall_end
	//   - active: list[Event]:
	//     - let s = the current time segment
	//     - forall e in list[Event] (e.start < s.start+s.dur & e.end >= s.start+s.dur)
	//   - no two consecutive time segments have the same active events

	f.Add(int64(0), int64(0), int32(0), uint16(0))
	f.Fuzz(func(t *testing.T, seed, overallStartSec int64, overallDuration int32, eventLength uint16) {
		r := rand.New(rand.NewSource(seed))

		overallStart := time.Unix(abs(overallStartSec), 0)
		overallEnd := overallStart.Add(time.Second * time.Duration(abs(overallDuration)))

		eventList := make([]dto.Event, eventLength)
		for i := range eventLength {
			offset := time.Duration(abs(r.Int31())) * time.Second
			dur := time.Duration(abs(r.Int31())) * time.Second
			start := overallStart.Add(offset)
			eventList[i] = dto.Event{
				Start: events.Datetime{
					Stamp: start,
				},
				End: events.Datetime{
					Stamp: start.Add(dur),
				},
			}
		}

		// all invariants are validated within the function itself which will
		// cause the function to panic if an invariant is violated
		convertToTimeline(eventList, overallStart, overallEnd)
	})
}
