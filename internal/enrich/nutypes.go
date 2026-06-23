package enrich

import (
	"github.com/LQR471814/nu_plugin_caldav/internal/nuconv"
	"github.com/ainvaltin/nu-plugin/types"
)

func optionalType(t types.Type) types.Type {
	return types.OneOf(t, types.Nothing())
}

var (
	// EventRecordDef describes the Nushell record produced by Event.ToNu.
	EventRecordDef = types.RecordDef{
		"uid":                        types.String(),
		"summary":                    optionalType(types.String()),
		"location":                   optionalType(types.String()),
		"description":                optionalType(types.String()),
		"categories":                 optionalType(types.List(types.String())),
		"datetime_stamp":             optionalType(types.Record(nuconv.DatetimeType)),
		"created":                    optionalType(types.Record(nuconv.DatetimeType)),
		"last_modified":              optionalType(types.Record(nuconv.DatetimeType)),
		"class":                      optionalType(nuconv.EventClassType),
		"geo":                        optionalType(types.Record(nuconv.EventGeoType)),
		"priority":                   optionalType(types.Int()),
		"sequence":                   optionalType(types.Int()),
		"status":                     optionalType(nuconv.EventStatusType),
		"transparency":               optionalType(nuconv.EventTransparencyType),
		"url":                        optionalType(types.String()),
		"comment":                    optionalType(types.String()),
		"attach":                     optionalType(types.String()),
		"contact":                    optionalType(types.String()),
		"organizer":                  optionalType(types.String()),
		"start":                      types.Record(nuconv.DatetimeType),
		"end":                        types.Record(nuconv.DatetimeType),
		"duration":                   optionalType(types.Duration()),
		"recurrence_rule":            optionalType(nuconv.RRuleType),
		"recurrence_dates":           optionalType(types.Table(nuconv.DatetimeType)),
		"recurrence_exception_dates": optionalType(types.Table(nuconv.DatetimeType)),
		"recurrence_instance":        optionalType(types.Record(nuconv.DatetimeType)),
		"trigger":                    optionalType(types.Record(nuconv.EventTriggerType)),
	}

	// EventType is the Nushell type annotation for Event.
	EventType types.Type = types.Record(EventRecordDef)

	// EventTableType is the Nushell table annotation for a list of Event records.
	EventTableType types.Type = types.Table(EventRecordDef)

	// CalendarObjectRecordDef describes the Nushell record produced by CalendarObject.ToNu.
	CalendarObjectRecordDef = types.RecordDef{
		"object_path": optionalType(types.String()),
		"main":        EventType,
		"overrides":   optionalType(EventTableType),
	}

	// CalendarObjectType is the Nushell type annotation for CalendarObject.
	CalendarObjectType types.Type = types.Record(CalendarObjectRecordDef)

	// CalendarObjectTableType is the Nushell table annotation for a list of CalendarObject records.
	CalendarObjectTableType types.Type = types.Table(CalendarObjectRecordDef)
)
