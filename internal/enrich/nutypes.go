package enrich

import (
	"github.com/LQR471814/nu_plugin_caldav/internal/nuconv"
	"github.com/ainvaltin/nu-plugin/types"
)

var (
	// EventRecordDef describes the Nushell record produced by Event.ToNu.
	EventRecordDef = types.RecordDef{
		"uid":                        types.String(),
		"summary":                    types.String(),
		"location":                   types.String(),
		"description":                types.String(),
		"categories":                 types.List(types.String()),
		"datetime_stamp":             types.Record(nuconv.DatetimeType),
		"created":                    types.Record(nuconv.DatetimeType),
		"last_modified":              types.Record(nuconv.DatetimeType),
		"class":                      nuconv.EventClassType,
		"geo":                        types.Record(nuconv.EventGeoType),
		"priority":                   types.Int(),
		"sequence":                   types.Int(),
		"status":                     nuconv.EventStatusType,
		"transparency":               nuconv.EventTransparencyType,
		"url":                        types.String(),
		"comment":                    types.String(),
		"attach":                     types.String(),
		"contact":                    types.String(),
		"organizer":                  types.String(),
		"start":                      types.Record(nuconv.DatetimeType),
		"end":                        types.Record(nuconv.DatetimeType),
		"duration":                   types.Duration(),
		"recurrence_rule":            nuconv.RRuleType,
		"recurrence_dates":           types.Table(nuconv.DatetimeType),
		"recurrence_exception_dates": types.Table(nuconv.DatetimeType),
		"recurrence_instance":        types.Record(nuconv.DatetimeType),
		"trigger":                    types.Record(nuconv.EventTriggerType),
	}

	// EventType is the Nushell type annotation for Event.
	EventType types.Type = types.Record(EventRecordDef)

	// EventTableType is the Nushell table annotation for a list of Event records.
	EventTableType types.Type = types.Table(EventRecordDef)

	// CalendarObjectRecordDef describes the Nushell record produced by CalendarObject.ToNu.
	CalendarObjectRecordDef = types.RecordDef{
		"object_path": types.String(),
		"main":        EventType,
		"overrides":   EventTableType,
	}

	// CalendarObjectType is the Nushell type annotation for CalendarObject.
	CalendarObjectType types.Type = types.Record(CalendarObjectRecordDef)

	// CalendarObjectTableType is the Nushell table annotation for a list of CalendarObject records.
	CalendarObjectTableType types.Type = types.Table(CalendarObjectRecordDef)
)
