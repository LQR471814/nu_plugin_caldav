package main

import (
	"context"
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/LQR471814/nu_plugin_caldav/internal/dto"
	"github.com/LQR471814/nu_plugin_caldav/internal/nuconv"
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
				Out: nuconv.TimelineType,
			},
		},
	},
	OnRun: timelineCmdExec,
}

func init() {
	commands = append(commands, timelineCmd)
}

func expandEvents(out *[]dto.Event, object dto.EventObject, start, end time.Time) (err error) {
	set := &rrule.Set{}

	if object.Main.RecurrenceRule.RRule != nil {
		set.RRule(object.Main.RecurrenceRule.RRule)
	} else {
		// single recurrence rule
		var rule *rrule.RRule
		rule, err = rrule.NewRRule(rrule.ROption{
			Dtstart: object.Main.Start.Stamp,
			Count:   1,
		})
		if err != nil {
			return
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
			err = fmt.Errorf("recurrence start time cannot be < overall start time: %v", startTime)
			return
		}

		if startTime.After(end) {
			err = fmt.Errorf("recurrence start time cannot be > overall end time: %v", startTime)
			return
		}

		replica := object.Main
		dur := replica.End.Stamp.Sub(replica.Start.Stamp)
		replica.Start.Stamp = startTime
		replica.End.Stamp = startTime.Add(dur)
		*out = append(*out, replica)
	}

	return
}

func convertToTimeline(eventList []dto.Event, start, end time.Time) (out []dto.TimeSegment, err error) {
	// pre-conditions
	if end.Before(start) {
		err = fmt.Errorf("timeline END cannot be before START")
		return
	}

	for _, e := range eventList {
		if e.End.Stamp.Before(e.Start.Stamp) {
			err = fmt.Errorf("event END cannot be before the event's START: %v", e)
			return
		}
		if e.Start.Stamp.Before(start) {
			err = fmt.Errorf("event START cannot be before timeline START: %v", e)
			return
		}
	}

	if len(eventList) == 0 {
		return
	}

	slices.SortFunc(eventList, func(a, b dto.Event) int {
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
	var active []dto.Event

	for !t.Equal(end) && !t.After(end) {
		// --- collect active events

		// remove inactive events
		var newActive []dto.Event
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

		out = append(out, dto.TimeSegment{
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

func timelineCmdExec(ctx context.Context, call *nu.ExecCommand) (err error) {
	start, ok := call.Named["start"].Value.(time.Time)
	if !ok {
		err = fmt.Errorf("must specify -start")
		return
	}

	end, ok := call.Named["end"].Value.(time.Time)
	if !ok {
		err = fmt.Errorf("must specify -end")
		return
	}

	objects, err := recvListInput(call, nuconv.EventObjectFromNu)
	if err != nil {
		err = fmt.Errorf("recv list input: %w", err)
		return
	}

	var eventList []dto.Event
	for _, obj := range objects {
		expandEvents(&eventList, obj, start, end)
	}
	timeline, err := convertToTimeline(eventList, start, end)
	if err != nil {
		err = fmt.Errorf("convert to timeline: %w", err)
		return
	}

	out, err := nuconv.TimelineToNu(timeline)
	if err != nil {
		err = fmt.Errorf("timeline to nu: %w", err)
		return
	}

	err = call.ReturnValue(ctx, out)
	if err != nil {
		err = fmt.Errorf("return value: %w", err)
		return
	}
	return
}
