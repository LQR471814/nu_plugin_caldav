# nu_plugin_caldav

> A nushell plugin for interfacing with a CalDAV server.

## Commands

| Command                                | Description                                                         |
|----------------------------------------|---------------------------------------------------------------------|
| `caldav query homeset [principal]`     | Find a homeset from CalDAV (optionally given a principal username). |
| `caldav query calendars <homeset>`     | Reads calendars for a given homeset from CalDAV.                    |
| `caldav query events <calendar_path>`  | Reads events from a given CalDAV calendar.                          |
| `caldav upsert events <calendar_path>` | Updates or inserts events from the given input.                     |

## Configuration

Server configuration is done through environment variables:

- `NU_PLUGIN_CALDAV_URL`: Full URL (ex. https://hostname/...)
  to the CalDAV server.
- `NU_PLUGIN_CALDAV_USERNAME`: Username for authentication with
  the CalDAV server, if not necessary, one can leave it unset or
  blank.
- `NU_PLUGIN_CALDAV_PASSWORD`: Password for authentication with
  the CalDAV server, if not necessary, one can leave it unset or
  blank.
- `NU_PLUGIN_CALDAV_INSECURE`: Set to `1` if HTTPS security errors
  should be ignored, optional.

## Types

```go
var calendarType = types.Table(types.RecordDef{
	// path to the calendar
	"path": types.String(),
	// name of the calendar
	"name": types.String(),
	// description of the calendar
	"description": types.String(),
	// quota information
	"max_resource_size": types.Int(),
	// support information
	"supported_component_set": types.List(types.String()),
})

var dateType = types.Record(types.RecordDef{
	// datetime represents the underlying datetime information of the date &
	// time
	"datetime": types.Date(),
	// if all_day is true, just the date should be considered and the time
	// ignored
	"all_day": types.Bool(),
	// if floating is true, the timezone of the date is always at the local
	// time of whoever is using the calendar and not "official" timezone should
	// be given to it
	"floating": types.Bool(),
})

var eventsType = types.Table(types.RecordDef{
	// universal id for the event, unique across calendars
	"uid": types.String(),
	// name of the event
	"name": types.String(),
	// location the event occurs at
	"location": types.String(),
	// description of the event (possibly multiline)
	"description": types.String(),
	// list of "tags" for the event
	"categories": types.List(types.String()),
	// event start
	"start": dateType,
	// event end
	"end": dateType,
	// recurrence_id (optional) designates the event as an override of another
	// recurring event (if set, recurrence_set must not be set)
	"recurrence_id": dateType,
	// recurrence_set (optional) designates the event as the originator of a
	// recurring event (if set, recurrence_id must not be set)
	"recurrence_set": types.Record(types.RecordDef{
		// recurrence rule (required)
		"rule": types.String(),
		// defines which recurrences should not occur (optional)
		//
		// note: the timezone of these dates can only be one of:
		// 1. UTC
		// 2. Floating time (no timezone)
		// 3. A single other explicit time zone
		//
		// Ex.
		// 	OK:
		// 		(UTC, floating, UTC, America/Los_Angeles, America/Los_Angeles)
		// 	BAD:
		// 		(UTC, floating, UTC, Asia/Shanghai, America/Los_Angeles)
		// 	There can only be one other explicit time zone outside of UTC.
		"exceptions": types.List(dateType),
		// defines which additional dates recurrences should occur (optional)
        //
		// note: all recurrences' timezone will be:
		// 1. UTC
		// 2. Floating time (no timezone)
		// 3. Specific time zone
		"dates": types.List(dateType),
	}),
	// trigger (optional) defines a notification trigger for the event
	"trigger": types.Record(types.RecordDef{
		// if set, notification will be triggered this duration before the
		// event (if set absolute should not be set)
		"relative": types.Duration(),
		// if set, notification will be triggered at a given absolute time (if
		// set, relative should not be set)
		"absolute": dateType,
	}),
})
```

## Limitations

- For now, this plugin will fetch **all** events from the server
  everytime `caldav events` is called, while this is ok for
  functionality, it may not be performant.

