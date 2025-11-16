package main

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
	"github.com/emersion/go-webdav/caldav"
)

var deleteEventsCmd = &nu.Command{
	Signature: nu.PluginSignature{
		Name:        "caldav delete events",
		Category:    "Network",
		Desc:        "Deletes events from a calendar",
		SearchTerms: []string{"caldav", "delete", "events"},
		Named: []nu.Flag{
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
				// a list of object_paths
				In:  types.List(types.String()),
				Out: types.Nothing(),
			},
		},
	},
	OnRun: deleteEventsCmdExec,
}

func init() {
	commands = append(commands, deleteEventsCmd)
}

type deleteEventCtx struct {
	ctx          context.Context
	calendarPath string
	client       *caldav.Client
}

func deleteEventObjectWorker(
	ctx deleteEventCtx,
	objectPaths chan string,
	status chan error,
) {
	for objpath := range objectPaths {
		err := ctx.client.RemoveAll(ctx.ctx, objpath)
		status <- err
	}
}

func multiDeleteObjects(
	ctx deleteEventCtx,
	objects []string,
	parallel int,
) (err error) {
	objectChan := make(chan string)
	status := make(chan error)
	defer close(objectChan)
	defer close(status)

	for range parallel {
		go deleteEventObjectWorker(ctx, objectChan, status)
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

func deleteEventsCmdExec(ctx context.Context, call *nu.ExecCommand) (err error) {
	defer func() {
		res := recover()
		if res != nil {
			err = fmt.Errorf("Panic: %v\n%s", res, string(debug.Stack()))
		}
	}()

	// parse flags
	client, err := getClient(ctx, call)
	if err != nil {
		return
	}
	calendarPath, err := tryCast[string](call.Positional[0])
	if err != nil {
		return
	}
	parallel := 1
	v, ok := call.FlagValue("parallel")
	if ok {
		parallel = v.Value.(int)
	}

	inputs, err := recvListInput(call, func(v nu.Value) string { return v.Value.(string) })
	if err != nil {
		return
	}

	// start workers & send jobs
	err = multiDeleteObjects(deleteEventCtx{
		ctx:          ctx,
		client:       client,
		calendarPath: calendarPath,
	}, inputs, parallel)
	if err != nil {
		return
	}

	return
}
