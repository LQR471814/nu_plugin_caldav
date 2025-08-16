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
				In: types.Nothing(),
				Out: types.Table(types.RecordDef{
					"path":                    types.String(),
					"name":                    types.String(),
					"description":             types.String(),
					"max_resource_size":       types.Int(),
					"supported_component_set": types.List(types.String()),
				}),
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

	values := make([]map[string]any, len(calendars))
	for i, cal := range calendars {
		values[i] = map[string]any{
			"path":                    cal.Path,
			"name":                    cal.Name,
			"description":             cal.Description,
			"max_resource_size":       cal.MaxResourceSize,
			"supported_component_set": cal.SupportedComponentSet,
		}
	}

	err = call.ReturnValue(ctx, nu.ToValue(values))
	return
}
