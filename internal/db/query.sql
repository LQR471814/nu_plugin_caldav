-- name: ReadMetadata :one
select version from metadata
where id = 1;

-- name: PutMetadata :exec
insert into metadata (id, version)
values (1, ?)
on conflict (id) do update set
	version = excluded.version;

-- name: PutCalendar :exec
insert into calendar (path, dto, sync_token)
values (?, ?, ?)
on conflict (path) do update set
	dto = excluded.dto,
	sync_token = excluded.sync_token;

-- name: UpdateCalSyncToken :exec
insert into calendar (path, sync_token)
values (?, ?)
on conflict (path) do update set
	sync_token = excluded.sync_token;

-- name: PutEvent :exec
insert into event_object (path, calendar_path, dto)
values (?, ?, ?)
on conflict (path) do update set
	calendar_path = excluded.calendar_path,
	dto = excluded.dto;

-- name: ReadCalendars :many
select * from calendar;

-- name: ReadCalendarSyncToken :one
select sync_token from calendar where path = ?;

-- name: ReadCalendarDto :one
select dto from calendar where path = ?;

-- name: ReadEvents :many
select path, dto from event_object where calendar_path = ?;

-- name: DeleteEvents :exec
delete from event_object
where path in (sqlc.slice('paths'));

