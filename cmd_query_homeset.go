package main

import (
	"context"

	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
)

var queryHomesetCmd = &nu.Command{
	Signature: nu.PluginSignature{
		Name:        "caldav query homeset",
		Category:    "Network",
		Desc:        "Finds a homeset ID from CalDAV (optionally given a principal username).",
		SearchTerms: caldavKeywordsQuery("calendar", "homeset"),
		OptionalPositional: []nu.PositionalArg{
			{
				Name:  "principal",
				Desc:  "Usually the username of the one that owns the calendar, can be left blank if CalDAV URL already includes the principal path. The principal of the current user can also be found by the `caldav query principal` command.",
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
	OnRun: queryHomesetCmdExec,
}

func init() {
	commands = append(commands, queryHomesetCmd)
}

func queryHomesetCmdExec(ctx context.Context, call *nu.ExecCommand) (err error) {
	client, err := getClient(ctx, call)
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
