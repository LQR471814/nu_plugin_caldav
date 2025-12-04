# nu_plugin_caldav

> A nushell plugin for interfacing with a CalDAV server.

## Commands

| Command                                                                                      | Description                                                                               |
|----------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------|
| `caldav query principal`                                                                     | Get the principal user path for the current configured user.                              |
| `caldav query homeset [principal]`                                                           | Find a homeset (collection of calendars) from CalDAV (optionally given a principal path). |
| `caldav query calendars <homeset>`                                                           | Reads the list calendars of calendars under a homeset from the CalDAV server.             |
| `caldav query events <calendar_path> [--start] [--end] [--text-match] [--text-match-negate]` | Reads events from a given calendar.                                                       |
| `<calendar_events> \| caldav save events <calendar_path> [--update]`                         | Creates (optionally updates if already existing) events from the given input.             |
| `<calendar_events> \| caldav timeline [--start] [--end]`                                     | Orders events chronologically.                                                            |

## Configuration

Server configuration is done through environment variables:

- `NU_PLUGIN_CALDAV_URL`: Full URL (ex. https://hostname/...) to
  the CalDAV server.
- `NU_PLUGIN_CALDAV_USERNAME`: Username for authentication with
  the CalDAV server. (optional)
- `NU_PLUGIN_CALDAV_PASSWORD`: Password for authentication with
  the CalDAV server. (optional)
- `NU_PLUGIN_CALDAV_INSECURE`: Set to `1` if HTTPS security errors
  should be ignored. (optional)

## Example Usage

https://github.com/LQR471814/nu_plugin_caldav/blob/3fb5759ae5033a7cc5db553f61e0c8614b148e4d/example.nu#L1-L46

## Limitations

- Server-side event filtering is currently limited to just start
  and end.
- Static validation of event type is currently not possible due to
  nushell's lack of optional types.
- Various fields in event are currently not implemented:
    - Binary attachment
    - Event scheduling fields / management of RSVP
- Calendar creation/deletion/rename is currently not implemented.
- Support for VTODO, VJOURNAL, etc... is currently not
  implemented.

