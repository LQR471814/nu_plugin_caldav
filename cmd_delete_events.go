package main

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/types"
	"github.com/emersion/go-webdav/caldav"
)

var deleteEventsCmd = &nu.Command{
	Signature: nu.PluginSignature{
		Name:        "caldav delete events",
		Category:    "Network",
		Desc:        "Deletes event objects given their paths",
		SearchTerms: []string{"caldav", "delete", "events"},
		Named: []nu.Flag{
			{
				Long:    "parallel",
				Short:   'p',
				Default: &defaultParallelism,
				Desc:    "Controls the amount of requests that can be made in parallel.",
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

type deleteEventJob struct {
	client  *caldav.Client
	objpath string
}

func (j deleteEventJob) Do(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return j.client.RemoveAll(ctx, j.objpath)
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
	parallel := 1
	v, ok := call.FlagValue("parallel")
	if ok {
		parallel = v.Value.(int)
	}

	inputs, err := recvListInput(call, func(v nu.Value) (string, error) { return tryCast[string](v) })
	if err != nil {
		return
	}

	jobs := make([]job, len(inputs))
	for i, objpath := range inputs {
		jobs[i] = deleteEventJob{
			client:  client,
			objpath: objpath,
		}
	}

	err = parallelizeJobs(ctx, jobs, parallel)
	if err != nil {
		return
	}

	return
}
