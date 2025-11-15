package main

import (
	"context"
	"fmt"
	"path"
	"runtime/debug"
	"strings"
	"time"

	"github.com/LQR471814/nu_plugin_caldav/events"
	"github.com/LQR471814/nu_plugin_caldav/internal/nutypes"
	"github.com/LQR471814/nu_plugin_caldav/internal/nutypes/conversions"
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
		SearchTerms: []string{"caldav", "upsert", "events"},
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
	objectReplicas []nutypes.EventObjectReplica,
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

func putEventObjectWorker(
	ctx saveEventCtx,
	objects chan events.EventObject,
	status chan error,
) {
	for obj := range objects {
		objpath := obj.ObjectPath
		if objpath == "" {
			uid, _ := obj.Main.GetUID()
			objpath = path.Join(ctx.calendarPath, uid)
		}
		_, err := ctx.client.PutCalendarObject(ctx.ctx, objpath, obj.ToCalendar())
		status <- err
	}
}

func multiPutObjects(
	ctx saveEventCtx,
	objects []events.EventObject,
	parallel int,
) (err error) {
	objectChan := make(chan events.EventObject)
	status := make(chan error)
	defer close(objectChan)
	defer close(status)

	for range parallel {
		go putEventObjectWorker(ctx, objectChan, status)
	}
	finished := 0
	for _, obj := range objects {
		select {
		case <-ctx.ctx.Done():
			err = fmt.Errorf("context canceled")
			return
		case objectChan <- obj:
		case err = <-status:
			if err != nil {
				return
			}
			finished++
		}
	}
	for finished < len(objects) {
		select {
		case <-ctx.ctx.Done():
			err = fmt.Errorf("context canceled")
			return
		case err = <-status:
			if err != nil {
				return
			}
			finished++
		}
	}
	return
}

// apply default property updates to new/modified events
func applyDefaultUpdates(event events.Event, objectPath string, dtstamp *ical.Prop, now events.Datetime) {
	event.Props.Set(dtstamp)

	// escape text
	text, ok := event.GetLocation()
	if ok {
		text = strings.ReplaceAll(text, "\n", "\\n")
		event.SetLocation(&text)
	}
	text, ok = event.GetComment()
	if ok {
		text = strings.ReplaceAll(text, "\n", "\\n")
		event.SetComment(&text)
	}
	text, ok = event.GetDescription()
	if ok {
		text = strings.ReplaceAll(text, "\n", "\\n")
		event.SetDescription(&text)
	}
	text, ok = event.GetContact()
	if ok {
		text = strings.ReplaceAll(text, "\n", "\\n")
		event.SetContact(&text)
	}

	if uid, _ := event.GetUID(); uid == "" {
		uid, err := uuid.NewRandom()
		if err != nil {
			panic(err)
		}
		event.SetUID(uid.String())
	}

	if objectPath == "" {
		event.SetLastModified(&now)
		return
	}
	event.SetCreated(&now)
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
	inputObjectReplicas, err := recvListInput(call, conversions.EventObjectReplicaFromNu)
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
		var updateObjectReplicas []nutypes.EventObjectReplica
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
		replica.Main.Apply(obj.Main)
		for _, override := range replica.Overrides {
			ev := events.Event{Timezone: time.Local, Event: *ical.NewEvent()}
			override.Apply(ev)
			obj.Overrides = append(obj.Overrides, ev)
		}
		putObjects = append(putObjects, obj)
	}

	// apply default property updates to new/modified objects
	dtstamp := ical.NewProp(ical.PropDateTimeStamp)
	dtstamp.SetDateTime(currentTime)
	now := events.Datetime{Stamp: currentTime}
	for _, obj := range putObjects {
		applyDefaultUpdates(obj.Main, obj.ObjectPath, dtstamp, now)
		for _, override := range obj.Overrides {
			applyDefaultUpdates(override, obj.ObjectPath, dtstamp, now)
		}
	}

	// start workers & send jobs
	err = multiPutObjects(subctx, putObjects, parallel)
	if err != nil {
		return
	}

	return
}
