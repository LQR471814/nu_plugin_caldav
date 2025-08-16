# nu_plugin_caldav

> A nushell plugin for interfacing with a CalDAV server.

## Commands

| Command                                | Description                                                         |
|----------------------------------------|---------------------------------------------------------------------|
| `caldav query homeset [principal]`     | Find a homeset from CalDAV (optionally given a principal username). |
| `caldav query calendars <homeset>`     | Reads calendars for a given homeset from CalDAV.                    |
| `caldav query events <calendar_path>`  | Reads events from a given CalDAV calendar.                          |
| `caldav upsert events <calendar_path>` | Updates or inserts events from the given input.                     |
| `<event record> \| caldav recurrences` | Render recurrences of a given event passed in.                      |

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

```typescript
// note: passing nothing to a field that does not have the `?`
// annotation will cause the plugin to return an error.

type Calendar = {
    path: string
    name: string
    description: string
    max_resource_size: int
    supported_component_set: string[]
}

type Event = {
	uid: string
	name: string
	location: string
	description: string
    categories: string[]
    start: datetime
    end: datetime
    recurrence_id?: datetime
    recurrence_rule?: string
    recurrence_exceptions?: datetime[]
    trigger?: {
        relative?: duration
        absolute?: datetime
    }
}

type UpdatedEvent = {
	uid: string
	name?: string
	location?: string
	description?: string
    categories?: string[]
    start?: datetime
    end?: datetime
    recurrence_id?: datetime
    recurrence_rule?: string
    recurrence_exceptions?: datetime[]
    trigger?: {
        relative?: duration
        absolute?: datetime
    }
}
```

## Limitations

- For now, this plugin will fetch **all** events from the server
  everytime `caldav events` is called, while this is ok for
  functionality, it may not be performant.

