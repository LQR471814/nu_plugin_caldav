package db

import (
	"context"
)

const readEvents = `-- name: ReadEvents :many
select path, dto from event_object where calendar_path = ?
`

type ReadEventsRow struct {
	Path string
	Dto  []byte
}

func (q *Queries) ReadEvents(ctx context.Context, calendarPath string, out chan ReadEventsRow) error {
	rows, err := q.db.QueryContext(ctx, readEvents, calendarPath)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var i ReadEventsRow
		if err := rows.Scan(&i.Path, &i.Dto); err != nil {
			return err
		}
		out <- i
	}
	if err := rows.Close(); err != nil {
		return err
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}
