package main

import (
	"context"
	"fmt"

	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
)

var homesetCmd = &nu.Command{
	Signature: nu.PluginSignature{
		Name:        "caldav query homeset",
		Category:    "Network",
		Desc:        "Find a homeset ID from CalDAV (optionally given a principal username).",
		SearchTerms: []string{"caldav", "homeset"},
		OptionalPositional: []nu.PositionalArg{
			{
				Name:  "principal",
				Desc:  "Usually the username of the one that owns the calendar, can be left blank if CalDAV URL already includes the principal path.",
				Shape: syntaxshape.String(),
			},
		},
		InputOutputTypes: []nu.InOutTypes{
			{
				In:  types.Nothing(),
				Out: types.String(),
			},
		},
	},
	OnRun: homesetCmdExec,
}

func init() {
	commands = append(commands, homesetCmd)
}

func homesetCmdExec(ctx context.Context, call *nu.ExecCommand) (err error) {
	client, err := getClientFromEnv()
	if err != nil {
		return
	}

	var principal string
	switch in := call.Input.(type) {
	case nil:
	case nu.Value:
		switch vt := in.Value.(type) {
		case string:
			principal = vt
		default:
			err = fmt.Errorf("unsupported input type %T", call.Input)
			return
		}
	default:
		err = fmt.Errorf("unsupported input type %T", call.Input)
		return
	}

	homeSet, err := client.FindCalendarHomeSet(ctx, principal)
	if err != nil {
		return
	}

	err = call.ReturnValue(ctx, nu.ToValue(homeSet))
	return
}
