package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/shibukawa/configdir"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

const db_version = 1

func setupDB(ctx context.Context, tx *sql.Tx, txqry *Queries) (err error) {
	_, err = tx.ExecContext(ctx, schema)
	if err != nil {
		return
	}
	err = txqry.PutMetadata(ctx, db_version)
	return
}

const state_file = "state.db"

func Purge() (err error) {
	dirs := configdir.New("LQR471814", "nu_plugin_caldav")
	cache := dirs.QueryCacheFolder().Path
	return os.Remove(filepath.Join(cache, state_file))
}

func Open(ctx context.Context) (driver *sql.DB, qry *Queries, err error) {
	dirs := configdir.New("LQR471814", "nu_plugin_caldav")
	cache := dirs.QueryCacheFolder().Path

	err = os.MkdirAll(cache, 0777)
	if err != nil {
		return
	}
	driver, err = sql.Open("sqlite", fmt.Sprintf(
		"file:%s?"+
			"_journal_mode=WAL&"+
			"_synchronous=NORMAL&"+
			"_busy_timeout=10000",
		filepath.Join(cache, state_file),
	))
	if err != nil {
		return
	}
	err = driver.PingContext(ctx)
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
		err = setupDB(ctx, tx, txqry)
		if err != nil {
			return
		}
		tx.Commit()
		return
	}

	// if some unexpected error
	return
}
