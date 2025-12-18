-- metadata should contain exactly 1 row that contains metadata information for
-- this state
create table metadata (
	id int primary key,
	version int not null
);

-- event_object stores an event resource
create table event_object (
	path text primary key,
	calendar_path text not null references calendar(path)
		on update cascade
		on delete cascade,
	dto blob
);

-- calendar stores a calendar resource
create table calendar (
	path text primary key,
	sync_token text,
	dto blob
);

