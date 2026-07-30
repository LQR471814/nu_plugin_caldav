package main

import (
	"math"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/LQR471814/nu_plugin_caldav/internal/dto"
	"github.com/LQR471814/nu_plugin_caldav/internal/events"
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
	//   - start < overall_end

	// outputs:
	// - list[TimeSegment]:
	//   - start: int
	//   - dur: int
	//   - start+dur <= overall_end
	//   - active: list[Event]:
	//     - let s = the current time segment
	//     - forall e in list[Event] (e.start < s.start+s.dur & e.end >= s.start+s.dur)
	//   - no two consecutive time segments have the same active events

	const event_sample_size = 65565

	f.Add(uint64(0), uint64(0), uint64(0))
	f.Fuzz(func(t *testing.T, seed, overallStartSec, overallDuration uint64) {
		// since time.Time uses an int64, we prune values of uint64 that will
		// exceed the largest possible int64
		if overallStartSec > uint64(math.MaxInt64) || overallDuration > uint64(math.MaxInt64) {
			return
		}

		// we ensure that end doesn't overflow int64
		if uint64(math.MaxInt64)-overallStartSec < overallDuration {
			return
		}

		r := rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))

		eventList := make([]dto.Event, event_sample_size)
		for i := range event_sample_size {
			// events may be constructed completely disregarding the timeline
			// overall span because the timeline function should still be safe
			// to call regardless

			start := r.Int64()
			dur := time.Duration(r.Float64() * float64(math.MaxInt64-start))

			startStamp := time.Unix(start, 0)
			endStamp := startStamp.Add(dur)

			eventList[i] = dto.Event{
				Start: events.Datetime{
					Stamp: startStamp,
				},
				End: events.Datetime{
					Stamp: endStamp,
				},
			}
		}

		overallStart := time.Unix(abs(int64(overallStartSec)), 0)
		overallDur := time.Second * time.Duration(abs(int64(overallDuration)))
		overallEnd := overallStart.Add(overallDur)

		// all invariants are validated within the function itself which will
		// cause the function to panic if an invariant is violated
		convertToTimeline(eventList, overallStart, overallEnd)
	})
}
