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
	dto blob,
	-- the min_start and max_end fields exist for filtering purposes.
	--
	-- - min_start should indicate the minimum starting time across all
	-- instances of this event
	-- - max_end should indicate the maximum ending time across all
	-- instances of this event
	--
	-- all day events should be an event that starts at the start of the
	-- day and ends at the start of the next
	min_start datetime not null,
	max_end datetime not null
);

-- calendar stores a calendar resource
create table calendar (
	path text primary key,
	sync_token text
);

