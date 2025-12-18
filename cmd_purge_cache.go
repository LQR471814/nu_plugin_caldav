package main

import (
	"context"

	"github.com/LQR471814/nu_plugin_caldav/internal/db"
	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/types"
)

var purgeCacheCmd = &nu.Command{
	Signature: nu.PluginSignature{
		Name:        "caldav purge cache",
		Category:    "Misc",
		Desc:        "Completely clear cached events, calendars, and plugin state.",
		SearchTerms: []string{"caldav", "cache", "clear", "purge"},
		InputOutputTypes: []nu.InOutTypes{
			{
				In:  types.Nothing(),
				Out: types.Nothing(),
			},
		},
	},
	OnRun: purgeCacheCmdExec,
}

func init() {
	commands = append(commands, purgeCacheCmd)
}

func purgeCacheCmdExec(ctx context.Context, call *nu.ExecCommand) (err error) {
	return db.Purge()
}
