package main

import (
	"context"

	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
)

var homesetCmd = &nu.Command{
	Signature: nu.PluginSignature{
		Name:        "caldav query homeset",
		Category:    "Network",
		Desc:        "Find a homeset ID from CalDAV (optionally given a principal username).",
		SearchTerms: []string{"caldav", "query", "homeset"},
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
	client, err := getClientFromEnv(ctx, call)
	if err != nil {
		return
	}

	var principal string
	if len(call.Positional) > 0 {
		principal, err = tryCast[string](call.Positional[0])
		if err != nil {
			return
		}
	}

	homeSet, err := client.FindCalendarHomeSet(ctx, principal)
	if err != nil {
		return
	}

	err = call.ReturnValue(ctx, nu.ToValue(homeSet))
	return
}
