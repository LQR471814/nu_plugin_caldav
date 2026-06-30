package main

import (
	"context"
	"errors"
	"fmt"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/LQR471814/nu_plugin_caldav/internal/dto"
	"github.com/LQR471814/nu_plugin_caldav/internal/events"
	"github.com/LQR471814/nu_plugin_caldav/internal/nuconv"
	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
	"github.com/google/uuid"
)

var falseNu = nu.ToValue(false)
var defaultParallelism = nu.ToValue(4)

var saveEventsCmd = &nu.Command{
	Signature: nu.PluginSignature{
		Name:        "caldav save events",
		Category:    "Network",
		Desc:        "Saves events to a calendar",
		SearchTerms: []string{"caldav", "save", "events"},
		Named: []nu.Flag{
			{
				Long:    "update",
				Short:   'u',
				Default: &falseNu,
				Desc:    "Update events if they already exist instead of erroring. Note: When using this option, ",
			},
			{
				Long:    "parallel",
				Short:   'p',
				Default: &defaultParallelism,
				Desc:    "Controls the amount of requests that can be made in parallel.",
			},
		},
		RequiredPositional: []nu.PositionalArg{
			{
				Name:  "calendar_path",
				Desc:  "The `path` attribute of the calendar record returned by `caldav query calendars`.",
				Shape: syntaxshape.String(),
			},
		},
		InputOutputTypes: []nu.InOutTypes{
			{
				// TODO: fix typing later
				// since some fields may be omitted and nushell does not yet
				// support optional typing
				In:  types.Any(),
				Out: types.Nothing(),
			},
		},
	},
	OnRun: saveEventsCmdExec,
}

func init() {
	commands = append(commands, saveEventsCmd)
}

type saveEventCtx struct {
	ctx          context.Context
	calendarPath string
	client       *caldav.Client
}

// returns full event object(s) with updates applied
func makeUpdatedObjects(
	ctx saveEventCtx,
	objectReplicas []dto.EventObject,
) (out []events.EventObject, err error) {
	if len(objectReplicas) == 0 {
		return
	}

	for _, replica := range objectReplicas {
		if replica.ObjectPath == nil || *replica.ObjectPath == "" {
			// this is an assert, so it bypasses the regular error path
			panic(fmt.Errorf("event object must have object_path defined for update: %v", replica))
		}
	}

	paths := make([]string, len(objectReplicas))
	for i, replica := range objectReplicas {
		paths[i] = *replica.ObjectPath
	}
	objects, err := ctx.client.MultiGetCalendar(ctx.ctx, ctx.calendarPath, &caldav.CalendarMultiGet{
		Paths: paths,
		CompRequest: caldav.CalendarCompRequest{
			Name:     ical.CompEvent,
			AllProps: true,
		},
	})
	if err != nil {
		return
	}

	out = make([]events.EventObject, len(objects))
	for i, o := range objects {
		out[i] = events.EventObject{
			ObjectPath: o.Path,
		}
		for _, child := range o.Data.Children {
			prop := child.Props.Get(ical.PropRecurrenceID)
			ev := events.Event{
				Event:    ical.Event{Component: child},
				Timezone: time.Local,
			}
			if prop != nil {
				out[i].Overrides = append(out[i].Overrides, ev)
				continue
			}
			out[i].Main = ev
		}
	}

	return
}

type putEventObjectJob struct {
	calpath string
	client  *caldav.Client
	obj     events.EventObject
}

func (j putEventObjectJob) Do(ctx context.Context) (err error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	objpath := j.obj.ObjectPath
	if objpath == "" {
		uid, err := j.obj.Main.GetUID()
		if err != nil {
			return fmt.Errorf("get UID for calendar object path: %w", err)
		}
		objpath = path.Join(j.calpath, uid)
	}
	_, err = j.client.PutCalendarObject(ctx, objpath, j.obj.ToCalendar())
	return
}

func escapeTextProperty(name string, get func() (string, error), set func(*string)) error {
	text, err := get()
	if errors.Is(err, events.ErrPropertyNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("get %s: %w", name, err)
	}
	text = strings.ReplaceAll(text, "\n", "\\n")
	set(&text)
	return nil
}

// apply default property updates to new/modified events
func applyDefaultUpdates(event events.Event, objectPath string, dtstamp *ical.Prop, now events.Datetime) error {
	event.Props.Set(dtstamp)

	// escape text
	for _, prop := range []struct {
		name string
		get  func() (string, error)
		set  func(*string)
	}{
		{name: ical.PropLocation, get: event.GetLocation, set: event.SetLocation},
		{name: ical.PropComment, get: event.GetComment, set: event.SetComment},
		{name: ical.PropDescription, get: event.GetDescription, set: event.SetDescription},
		{name: ical.PropContact, get: event.GetContact, set: event.SetContact},
	} {
		if err := escapeTextProperty(prop.name, prop.get, prop.set); err != nil {
			return err
		}
	}

	uid, err := event.GetUID()
	if err != nil && !errors.Is(err, events.ErrPropertyNotFound) {
		return fmt.Errorf("get UID: %w", err)
	}
	if errors.Is(err, events.ErrPropertyNotFound) || uid == "" {
		uid, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		event.SetUID(uid.String())
	}

	if objectPath == "" {
		event.SetLastModified(&now)
		return nil
	}
	event.SetCreated(&now)
	return nil
}

func saveEventsCmdExec(ctx context.Context, call *nu.ExecCommand) (err error) {
	defer func() {
		res := recover()
		if res != nil {
			err = fmt.Errorf("Panic: %v\n%s", res, string(debug.Stack()))
		}
	}()

	currentTime := time.Now()

	// parse flags
	client, err := getClient(ctx, call)
	if err != nil {
		return
	}
	calendarPath, err := tryCast[string](call.Positional[0])
	if err != nil {
		return
	}
	update := false
	v, ok := call.FlagValue("update")
	if ok {
		update = v.Value.(bool)
	}
	parallel := 1
	v, ok = call.FlagValue("parallel")
	if ok {
		parallel = v.Value.(int)
	}

	subctx := saveEventCtx{
		ctx:          ctx,
		client:       client,
		calendarPath: calendarPath,
	}

	// process input events
	var putObjects []events.EventObject
	inputObjectReplicas, err := recvListInput(call, nuconv.EventObjectFromNu)
	if err != nil {
		return
	}
	if !update {
		for _, replica := range inputObjectReplicas {
			if replica.ObjectPath != nil && *replica.ObjectPath != "" {
				err = fmt.Errorf("Found event with object_path set. If you wish to create while updating existing events, please provide the -update flag.\n%v", replica)
				return
			}
		}
	} else {
		var updateObjectReplicas []dto.EventObject
		for _, replica := range inputObjectReplicas {
			if replica.ObjectPath != nil && *replica.ObjectPath == "" {
				continue
			}
			updateObjectReplicas = append(updateObjectReplicas, replica)
		}
		// add events to be updated to PUT queue
		putObjects, err = makeUpdatedObjects(subctx, updateObjectReplicas)
		if err != nil {
			return
		}
	}

	// make events that will be newly created
	for _, replica := range inputObjectReplicas {
		if replica.ObjectPath != nil && *replica.ObjectPath != "" {
			continue
		}
		obj := events.EventObject{}
		obj.Main = events.Event{Timezone: time.Local, Event: *ical.NewEvent()}
		err = replica.Main.Apply(obj.Main)
		if err != nil {
			return fmt.Errorf("apply main event: %w", err)
		}
		for _, override := range replica.Overrides {
			ev := events.Event{Timezone: time.Local, Event: *ical.NewEvent()}
			err = override.Apply(ev)
			if err != nil {
				return fmt.Errorf("apply override event: %w", err)
			}
			obj.Overrides = append(obj.Overrides, ev)
		}
		putObjects = append(putObjects, obj)
	}

	// apply default property updates to new/modified objects
	dtstamp := ical.NewProp(ical.PropDateTimeStamp)
	dtstamp.SetDateTime(currentTime)
	now := events.Datetime{Stamp: currentTime}
	for _, obj := range putObjects {
		err = applyDefaultUpdates(obj.Main, obj.ObjectPath, dtstamp, now)
		if err != nil {
			return
		}
		for _, override := range obj.Overrides {
			err = applyDefaultUpdates(override, obj.ObjectPath, dtstamp, now)
			if err != nil {
				return
			}
		}
	}

	jobs := make([]job, len(putObjects))
	for i, obj := range putObjects {
		jobs[i] = putEventObjectJob{
			calpath: calendarPath,
			client:  client,
			obj:     obj,
		}
	}
	err = parallelizeJobs(ctx, jobs, parallel)
	if err != nil {
		return
	}

	return
}
