package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/shibukawa/configdir"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

const db_version = 1

func setupDB(ctx context.Context, tx *sql.Tx, qry *Queries) (err error) {
	err = qry.PutMetadata(ctx, db_version)
	if err != nil {
		return
	}
	_, err = tx.ExecContext(ctx, schema)
	return
}

func Open(ctx context.Context) (driver *sql.DB, qry *Queries, err error) {
	dirs := configdir.New("LQR471814", "nu_plugin_caldav")
	cache := dirs.QueryCacheFolder().Path

	driver, err = sql.Open("sqlite", fmt.Sprintf(
		"file:%s?"+
			"cache=shared&"+
			"_journal_mode=WAL&"+
			"_synchronous=NORMAL&"+
			"_busy_timeout=10000&"+
			"_foreign_keys=1&"+
			"_pragma=temp_store(MEMORY)&"+
			"_pragma=cache_size(-20000)",
		filepath.Join(cache, "state.db"),
	))
	if err != nil {
		return
	}
	qry = New(driver)

	tx, err := driver.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer tx.Rollback()
	txqry := qry.WithTx(tx)

	version, err := txqry.ReadMetadata(ctx)

	// if db is already setup
	if err == nil && version == db_version {
		return
	}

	// if empty db
	if strings.Contains(err.Error(), "no such table") ||
		errors.Is(err, sql.ErrNoRows) ||
		version != db_version {
		err = setupDB(ctx, tx, qry)
		if err != nil {
			return
		}
		tx.Commit()
		return
	}

	// if some unexpected error
	return
}
