package main

import (
	"context"

	"github.com/LQR471814/nu_plugin_caldav/internal/db"
	"github.com/LQR471814/nu_plugin_caldav/internal/nuconv"
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
				Out: nuconv.CalendarListType,
			},
		},
	},
	OnRun: queryCalendarsCmdExec,
}

func init() {
	commands = append(commands, queryCalendarsCmd)
}

func queryCalendarsCmdExec(ctx context.Context, call *nu.ExecCommand) (err error) {
	driver, qry, err := db.Open(ctx)
	if err != nil {
		return
	}
	defer driver.Close()

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

	// this is done to clean up stale calendars (and their events) which are no
	// longer present in the given home set
	var paths []string
	for _, c := range calendars {
		paths = append(paths, c.Path)
	}
	err = qry.DeleteAllInPrincipalExcept(ctx, db.DeleteAllInPrincipalExceptParams{})
	if err != nil {
		return
	}

	out, err := nuconv.CalendarListToNu(calendars)
	if err != nil {
		return
	}
	err = call.ReturnValue(ctx, out)
	return
}
