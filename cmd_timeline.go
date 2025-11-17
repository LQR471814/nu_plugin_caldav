package main

import (
	"context"
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/LQR471814/nu_plugin_caldav/internal/nutypes"
	"github.com/LQR471814/nu_plugin_caldav/internal/nutypes/conversions"
	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
	"github.com/teambition/rrule-go"
)

var timelineCmd = &nu.Command{
	Signature: nu.PluginSignature{
		Name:        "caldav timeline",
		Desc:        "Combines multiple event objects into a single timeline.",
		SearchTerms: []string{"caldav", "timeline"},
		Category:    "Viewers",
		Named: []nu.Flag{
			{
				Long:  "start",
				Short: 's',
				Desc:  "Filter for all events after this start time.",
				Shape: syntaxshape.DateTime(),
			},
			{
				Long:  "end",
				Short: 'e',
				Desc:  "Filter for all events before this end time.",
				Shape: syntaxshape.DateTime(),
			},
		},
		InputOutputTypes: []nu.InOutTypes{
			{
				// TODO: fix typing later
				// since some fields may be omitted and nushell does not yet
				// support optional typing
				In:  types.Any(),
				Out: conversions.TimelineType,
			},
		},
	},
	OnRun: expandEventsCmdExec,
}

func init() {
	commands = append(commands, timelineCmd)
}

func expandEvents(out *[]nutypes.EventReplica, object nutypes.EventObjectReplica, start, end time.Time) {
	set := &rrule.Set{}

	if object.Main.RecurrenceRule != nil {
		set.RRule(object.Main.RecurrenceRule)
	} else {
		// single recurrence rule
		rule, err := rrule.NewRRule(rrule.ROption{
			Dtstart: object.Main.Start.Stamp,
			Count:   1,
		})
		if err != nil {
			panic(err)
		}
		set.RRule(rule)
	}

	rdates := make([]time.Time, len(object.Main.RecurrenceDates))
	for i, d := range object.Main.RecurrenceDates {
		rdates[i] = d.Stamp
	}
	set.SetRDates(rdates)

	exdates := make([]time.Time, len(object.Main.RecurrenceExceptionDates))
	for i, d := range object.Main.RecurrenceExceptionDates {
		exdates[i] = d.Stamp
	}
	// not a part of EXDATE, but added to prevent rrule from generating
	// recurrences on those specific date instances as overrides will
	// be added separately
	for _, override := range object.Overrides {
		exdates = append(exdates, override.RecurrenceInstance.Stamp)
		if override.Start.Stamp.Before(start) || override.Start.Stamp.After(end) {
			continue
		}
		*out = append(*out, override)
	}
	set.SetExDates(exdates)

	times := set.Between(start, end, true)
	for _, startTime := range times {
		if startTime.Before(start) {
			panic(fmt.Errorf("recurrence start time should not be before overall start time: %v", startTime))
		}
		if startTime.After(end) {
			panic(fmt.Errorf("recurrence start time should not be after overall end time: %v", startTime))
		}
		replica := object.Main
		dur := replica.End.Stamp.Sub(replica.Start.Stamp)
		replica.Start.Stamp = startTime
		replica.End.Stamp = startTime.Add(dur)
		*out = append(*out, replica)
	}
}

func convertToTimeline(eventList []nutypes.EventReplica, start, end time.Time) (out []nutypes.TimeSegment) {
	// pre-conditions
	if end.Before(start) {
		panic("timeline END cannot be before START")
	}
	for _, e := range eventList {
		if e.End.Stamp.Before(e.Start.Stamp) {
			panic(fmt.Errorf("event END cannot be before the event's START: %v", e))
		}
		if e.Start.Stamp.Before(start) {
			panic(fmt.Errorf("event START cannot be before timeline START: %v", e))
		}
	}

	if len(eventList) == 0 {
		return nil
	}
	slices.SortFunc(eventList, func(a, b nutypes.EventReplica) int {
		if a.Start.Stamp.Before(b.Start.Stamp) {
			return -1
		}
		if a.Start.Stamp.After(b.Start.Stamp) {
			return 1
		}
		return 0
	})

	// t is the current time
	t := start
	// nextEventCursor points to the next event whose start_time > t (may be undefined)
	nextEventCursor := 0
	// active keeps track of the events whose start_time < t && end_time > t (may be empty)
	var active []nutypes.EventReplica

	for !t.Equal(end) && !t.After(end) {
		// --- collect active events

		// remove inactive events
		var newActive []nutypes.EventReplica
		for _, activeEvent := range active {
			if activeEvent.End.Stamp.After(t) {
				newActive = append(newActive, activeEvent)
			}
		}

		// collect all active events and advasnce event cursor until reaching
		// an event whose start_time > t
		for nextEventCursor < len(eventList) {
			e := eventList[nextEventCursor]
			if e.Start.Stamp.After(t) {
				break
			}
			newActive = append(newActive, e)
			nextEventCursor++
		}

		active = newActive

		// --- advance to nearest next time

		// true if there is a next event's start or current active event's end
		// to advance to
		hasNext := false

		// find minimum next duration
		minNextDur := time.Duration(math.MaxInt64)
		for _, activeEvent := range active {
			dur := activeEvent.End.Stamp.Sub(t)
			if dur < minNextDur {
				minNextDur = dur
				hasNext = true
			}
		}
		if nextEventCursor < len(eventList) {
			dur := eventList[nextEventCursor].Start.Stamp.Sub(t)
			if dur < minNextDur {
				minNextDur = dur
				hasNext = true
			}
		}
		// ensure that the sum of the time segments' durations = end - start
		if !hasNext || t.Add(minNextDur).After(end) {
			minNextDur = end.Sub(t)
		}

		out = append(out, nutypes.TimeSegment{
			Now:          t,
			Duration:     minNextDur,
			ActiveEvents: active,
		})

		t = t.Add(minNextDur)
	}

	// post conditions
	total := time.Duration(0)
	for i, segment := range out {
		total += segment.Duration
		segmentEnd := segment.Now.Add(segment.Duration)
		if segmentEnd.After(end) {
			panic(fmt.Errorf("timeline segment cannot cross segment end: (idx: %d) of %+v", i, out))
		}
		for _, active := range segment.ActiveEvents {
			if active.Start.Stamp.After(segmentEnd) {
				panic(fmt.Errorf("active event cannot start after the time segment's end: (idx: %d) of %+v", i, out))
			}
			if active.End.Stamp.Before(segmentEnd) {
				panic(fmt.Errorf("active event should not end before the time segment's end: (idx: %d) of %+v", i, out))
			}
		}
		prev := i - 1
		if prev >= 0 && out[prev].ActiveEvents == nil && segment.ActiveEvents == nil {
			panic(fmt.Errorf(
				"no two consecutive time segments should both have nil active events: %v",
				out,
			))
		}
	}
	if total != end.Sub(start) {
		panic(fmt.Errorf(
			"sum of timeline segment total != end - start: got %+v expected %+v",
			total,
			end.Sub(start),
		))
	}

	return
}

func expandEventsCmdExec(ctx context.Context, call *nu.ExecCommand) (err error) {
	start, ok := call.Named["start"].Value.(time.Time)
	if !ok {
		err = fmt.Errorf("must specify -start")
		return
	}
	end, ok := call.Named["end"].Value.(time.Time)
	if !ok {
		err = fmt.Errorf("must specify -start")
		return
	}

	objects, err := recvListInput(call, conversions.EventObjectReplicaFromNu)
	if err != nil {
		return
	}

	var eventList []nutypes.EventReplica
	for _, obj := range objects {
		expandEvents(&eventList, obj, start, end)
	}
	timeline := convertToTimeline(eventList, start, end)

	err = call.ReturnValue(ctx, conversions.TimelineToNu(timeline))
	return
}
