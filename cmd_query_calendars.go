package main

import (
	"context"

	"github.com/LQR471814/nu_plugin_caldav/internal/nutypes/conversions"
	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
)

var queryCalendarsCmd = &nu.Command{
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
				Out: conversions.CalendarListType,
			},
		},
	},
	OnRun: queryCalendarsCmdExec,
}

func init() {
	commands = append(commands, queryCalendarsCmd)
}

func queryCalendarsCmdExec(ctx context.Context, call *nu.ExecCommand) (err error) {
	client, err := getClient(ctx, call)
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
	out, err := conversions.CalendarListToNu(calendars)
	if err != nil {
		return
	}
	err = call.ReturnValue(ctx, out)
	return
}
