package main

import (
	"context"

	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
)

var calendarsCmd = &nu.Command{
	Signature: nu.PluginSignature{
		Name:        "caldav query calendars",
		Category:    "Network",
		Desc:        "Reads calendars for a given homeset from CalDAV.",
		SearchTerms: []string{"caldav", "query", "calendars"},
		RequiredPositional: []nu.PositionalArg{
			{
				Name:  "homeset",
				Desc:  "The string returned by `caldav query homeset`.",
				Shape: syntaxshape.String(),
			},
		},
		InputOutputTypes: []nu.InOutTypes{
			{
				In:  types.Nothing(),
				Out: calendarType,
			},
		},
	},
	OnRun: calendarsCmdExec,
}

func init() {
	commands = append(commands, calendarsCmd)
}

func calendarsCmdExec(ctx context.Context, call *nu.ExecCommand) (err error) {
	client, err := getClientFromEnv(ctx, call)
	if err != nil {
		return
	}

	homeset, err := tryCast[string](call.Positional[0])
	if err != nil {
		return
	}

	calendars, err := client.FindCalendars(ctx, homeset)
	if err != nil {
		return
	}

	values := make([]nu.Value, len(calendars))
	for i, cal := range calendars {
		values[i] = nu.ToValue(nu.Record{
			"path":                    nu.ToValue(cal.Path),
			"name":                    nu.ToValue(cal.Name),
			"description":             nu.ToValue(cal.Description),
			"max_resource_size":       nu.ToValue(cal.MaxResourceSize),
			"supported_component_set": nu.ToValue(cal.SupportedComponentSet),
		})
	}

	err = call.ReturnValue(ctx, nu.ToValue(values))
	return
}
