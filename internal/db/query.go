package db

import (
	"context"
	"database/sql"
	"fmt"
)

const readEvents = `-- name: ReadEvents :many
select path, dto from event_object
where calendar_path = ?
and (:start is NULL or max_end > :start)
and (:end is NULL or min_start < :end)
`

type ReadEventsRow struct {
	Path string
	Dto  []byte
}

// ReadEventsRange limits the range of events which will be read.
//
// Specifically it will read all events such that both:
//   - event min_start < range end
//   - event max_end > range start
//
// What this means is that it also includes events which start before the
// start of the range, but end inside of the range, and events which end
// after the range end but start inside the range. Effectively functioning
// as an "inclusive" range selector.
//
// This reasoning for this is that it is more flexible (because the user
// can filter out undesirable events downstream, rather than the
// opposite).
//
// If start is undefined, start = -infinity
// If end is undefined, end = infinity
type ReadEventsRange struct {
	Start sql.NullTime
	End   sql.NullTime
}

func (q *Queries) ReadEvents(
	ctx context.Context,
	calendarPath string,
	out chan ReadEventsRow,
	rng ReadEventsRange,
) error {
	rows, err := q.db.QueryContext(
		ctx,
		readEvents, calendarPath,
		sql.Named("start", rng.Start),
		sql.Named("end", rng.End),
	)
	if err != nil {
		err = fmt.Errorf("query sql: %w", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var i ReadEventsRow
		if err := rows.Scan(&i.Path, &i.Dto); err != nil {
			err = fmt.Errorf("scan row: %w", err)
			return err
		}
		out <- i
	}

	if err := rows.Close(); err != nil {
		err = fmt.Errorf("close rows: %w", err)
		return err
	}
	if err := rows.Err(); err != nil {
		err = fmt.Errorf("iter: %w", err)
		return err
	}

	return nil
}
