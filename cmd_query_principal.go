package main

import (
	"context"

	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/types"
)

var queryPrincipal = &nu.Command{
	Signature: nu.PluginSignature{
		Name:        "caldav query principal",
		Category:    "Network",
		Desc:        "Finds the principal associated with the current user.",
		SearchTerms: caldavKeywordsQuery("calendar", "principal"),
		InputOutputTypes: []nu.InOutTypes{
			{
				In:  types.Nothing(),
				Out: types.String(),
			},
		},
	},
	OnRun: queryPrincipalExec,
}

func init() {
	commands = append(commands, queryPrincipal)
}

func queryPrincipalExec(ctx context.Context, call *nu.ExecCommand) (err error) {
	client, err := getClient(ctx, call)
	if err != nil {
		return
	}
	principal, err := client.FindCurrentUserPrincipal(ctx)
	if err != nil {
		return
	}
	err = call.ReturnValue(ctx, nu.ToValue(principal))
	return
}
