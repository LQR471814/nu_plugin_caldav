package conversions

import "github.com/ainvaltin/nu-plugin"
import "github.com/ainvaltin/nu-plugin/types"
import "github.com/LQR471814/nu_plugin_caldav/events"
import "github.com/LQR471814/nu_plugin_caldav/internal/nutypes"
import "github.com/teambition/rrule-go"
import "github.com/emersion/go-webdav/caldav"
import "net/url"
import "time"

var type_2493169154543297135 = types.String()
var type_5363327835607766502 = types.String()
var type_581713733206709016 = types.RecordDef{
	"object_path": type_15613163272824911089,
	"main":        type_1808007570533396940,
	"overrides":   type_6714134590091148549,
}
var type_11669970230249425419 = types.List(type_15613163272824911089)
var type_8047992331715851194 = types.Date()
var type_7161572108068222122 = types.RecordDef{
	"latitude":  type_17860233973098560385,
	"longitude": type_17860233973098560385,
}
var type_4579713456482216358 = types.Record(type_13545470577293064413)
var type_3582939560623804181 = types.Record(type_581713733206709016)
var type_2392839876798024645 = types.Table(type_581713733206709016)
var type_10893214126964739625 = types.Record(type_18369289839240265122)
var type_729807561129781588 = types.Bool()
var type_5454485661162817076 = types.RecordDef{
	"stamp":    type_8047992331715851194,
	"all_day":  type_729807561129781588,
	"floating": type_729807561129781588,
}
var type_3931126380996215332 = types.Table(type_5454485661162817076)
var type_15139881813094606131 = types.Int()
var type_16589689216511618220 = types.Duration()
var type_1808007570533396940 = types.Record(type_4505817543918974569)
var type_6714134590091148549 = types.Table(type_4505817543918974569)
var type_18369289839240265122 = types.RecordDef{
	"path":                    type_15613163272824911089,
	"name":                    type_15613163272824911089,
	"description":             type_15613163272824911089,
	"max_resource_size":       type_15139881813094606131,
	"supported_component_set": type_11669970230249425419,
}
var type_9513336738104922479 = types.Record(type_5454485661162817076)
var type_15385297846572725340 = types.String()
var type_13545470577293064413 = types.RecordDef{
	"relative":    type_16589689216511618220,
	"relative_to": type_15560982419391353847,
	"absolute":    type_8047992331715851194,
}
var type_4505817543918974569 = types.RecordDef{
	"uid":                        type_15613163272824911089,
	"summary":                    type_15613163272824911089,
	"location":                   type_15613163272824911089,
	"description":                type_15613163272824911089,
	"categories":                 type_11669970230249425419,
	"created":                    type_9513336738104922479,
	"last_modified":              type_9513336738104922479,
	"class":                      type_2493169154543297135,
	"geo":                        type_18222208942348113876,
	"priority":                   type_10890016574791629639,
	"sequence":                   type_10890016574791629639,
	"status":                     type_15385297846572725340,
	"transparency":               type_7057708295081751301,
	"url":                        type_5363327835607766502,
	"comment":                    type_15613163272824911089,
	"attach":                     type_5363327835607766502,
	"contact":                    type_15613163272824911089,
	"organizer":                  type_5363327835607766502,
	"start":                      type_9513336738104922479,
	"end":                        type_9513336738104922479,
	"recurrence_rule":            type_18334623996676649874,
	"recurrence_dates":           type_3931126380996215332,
	"recurrence_exception_dates": type_3931126380996215332,
	"recurrence_instance":        type_9513336738104922479,
	"trigger":                    type_4579713456482216358,
}
var type_10890016574791629639 = types.Int()
var type_7057708295081751301 = types.String()
var type_15560982419391353847 = types.Int()
var type_15613163272824911089 = types.String()
var type_17860233973098560385 = types.Float()
var type_18222208942348113876 = types.Record(type_7161572108068222122)
var type_18334623996676649874 = types.String()
var type_4202073097562803312 = types.RecordDef{
	"now":           type_8047992331715851194,
	"duration":      type_16589689216511618220,
	"active_events": type_6714134590091148549,
}
var type_14033482930202471129 = types.Record(type_4202073097562803312)
var type_4572434499988200822 = types.Table(type_18369289839240265122)
var type_14900610606431710770 = types.Table(type_4202073097562803312)

func type_8971279483973357571_FromNu(v nu.Value) *events.EventTransparency {
	if v.Value == nil {
		return nil
	}
	res := type_7057708295081751301_FromNu(v)
	return &res
}
func type_13545470577293064413_FromNu(v nu.Value) events.EventTrigger {
	out := events.EventTrigger{}
	record := v.Value.(nu.Record)
	var val nu.Value
	val, _ = record["relative"]
	out.Relative = type_5863190983406162214_FromNu(val)
	val, _ = record["relative_to"]
	out.RelativeTo = type_15560982419391353847_FromNu(val)
	val, _ = record["absolute"]
	out.Absolute = type_15050730807189225719_FromNu(val)
	return out
}
func type_14900610606431710770_FromNu(v nu.Value) nutypes.Timeline {
	if v.Value == nil {
		return nil
	}
	arr := v.Value.([]nu.Value)
	out := make(nutypes.Timeline, len(arr))
	for i, e := range arr {
		out[i] = type_4202073097562803312_FromNu(e)
	}
	return out
}
func type_12480522309550428545_FromNu(v nu.Value) *events.Datetime {
	if v.Value == nil {
		return nil
	}
	res := type_5454485661162817076_FromNu(v)
	return &res
}
func type_15385297846572725340_FromNu(v nu.Value) events.EventStatus {
	return events.EventStatus(v.Value.(string))
}
func type_16589689216511618220_FromNu(v nu.Value) time.Duration {
	return v.Value.(time.Duration)
}
func type_9520111014888170891_FromNu(v nu.Value) *events.EventTrigger {
	if v.Value == nil {
		return nil
	}
	res := type_13545470577293064413_FromNu(v)
	return &res
}
func type_10890016574791629639_FromNu(v nu.Value) int {
	return int(v.Value.(int64))
}
func type_7057708295081751301_FromNu(v nu.Value) events.EventTransparency {
	return events.EventTransparency(v.Value.(string))
}
func type_18369289839240265122_FromNu(v nu.Value) caldav.Calendar {
	out := caldav.Calendar{}
	record := v.Value.(nu.Record)
	var val nu.Value
	val, _ = record["path"]
	out.Path = type_15613163272824911089_FromNu(val)
	val, _ = record["name"]
	out.Name = type_15613163272824911089_FromNu(val)
	val, _ = record["description"]
	out.Description = type_15613163272824911089_FromNu(val)
	val, _ = record["max_resource_size"]
	out.MaxResourceSize = type_15139881813094606131_FromNu(val)
	val, _ = record["supported_component_set"]
	out.SupportedComponentSet = type_11669970230249425419_FromNu(val)
	return out
}
func type_4572434499988200822_FromNu(v nu.Value) nutypes.CalendarList {
	if v.Value == nil {
		return nil
	}
	arr := v.Value.([]nu.Value)
	out := make(nutypes.CalendarList, len(arr))
	for i, e := range arr {
		out[i] = type_18369289839240265122_FromNu(e)
	}
	return out
}
func type_11669970230249425419_FromNu(v nu.Value) []string {
	if v.Value == nil {
		return nil
	}
	arr := v.Value.([]nu.Value)
	out := make([]string, len(arr))
	for i, e := range arr {
		out[i] = type_15613163272824911089_FromNu(e)
	}
	return out
}
func type_9664538759823739797_FromNu(v nu.Value) *events.EventClass {
	if v.Value == nil {
		return nil
	}
	res := type_2493169154543297135_FromNu(v)
	return &res
}
func type_2584899110032584934_FromNu(v nu.Value) *int {
	if v.Value == nil {
		return nil
	}
	res := type_10890016574791629639_FromNu(v)
	return &res
}
func type_4505817543918974569_FromNu(v nu.Value) nutypes.EventReplica {
	out := nutypes.EventReplica{}
	record := v.Value.(nu.Record)
	var val nu.Value
	val, _ = record["uid"]
	out.Uid = type_17862013815172309399_FromNu(val)
	val, _ = record["summary"]
	out.Summary = type_17862013815172309399_FromNu(val)
	val, _ = record["location"]
	out.Location = type_17862013815172309399_FromNu(val)
	val, _ = record["description"]
	out.Description = type_17862013815172309399_FromNu(val)
	val, _ = record["categories"]
	out.Categories = type_11669970230249425419_FromNu(val)
	val, _ = record["created"]
	out.Created = type_12480522309550428545_FromNu(val)
	val, _ = record["last_modified"]
	out.LastModified = type_12480522309550428545_FromNu(val)
	val, _ = record["class"]
	out.Class = type_9664538759823739797_FromNu(val)
	val, _ = record["geo"]
	out.Geo = type_7163250051298988498_FromNu(val)
	val, _ = record["priority"]
	out.Priority = type_2584899110032584934_FromNu(val)
	val, _ = record["sequence"]
	out.Sequence = type_2584899110032584934_FromNu(val)
	val, _ = record["status"]
	out.Status = type_784588192188755836_FromNu(val)
	val, _ = record["transparency"]
	out.Transparency = type_8971279483973357571_FromNu(val)
	val, _ = record["url"]
	out.URL = type_5363327835607766502_FromNu(val)
	val, _ = record["comment"]
	out.Comment = type_17862013815172309399_FromNu(val)
	val, _ = record["attach"]
	out.Attach = type_5363327835607766502_FromNu(val)
	val, _ = record["contact"]
	out.Contact = type_17862013815172309399_FromNu(val)
	val, _ = record["organizer"]
	out.Organizer = type_5363327835607766502_FromNu(val)
	val, _ = record["start"]
	out.Start = type_5454485661162817076_FromNu(val)
	val, _ = record["end"]
	out.End = type_5454485661162817076_FromNu(val)
	val, _ = record["recurrence_rule"]
	out.RecurrenceRule = type_18334623996676649874_FromNu(val)
	val, _ = record["recurrence_dates"]
	out.RecurrenceDates = type_3931126380996215332_FromNu(val)
	val, _ = record["recurrence_exception_dates"]
	out.RecurrenceExceptionDates = type_3931126380996215332_FromNu(val)
	val, _ = record["recurrence_instance"]
	out.RecurrenceInstance = type_12480522309550428545_FromNu(val)
	val, _ = record["trigger"]
	out.Trigger = type_9520111014888170891_FromNu(val)
	return out
}
func type_5454485661162817076_FromNu(v nu.Value) events.Datetime {
	out := events.Datetime{}
	record := v.Value.(nu.Record)
	var val nu.Value
	var ok bool
	val, _ = record["stamp"]
	out.Stamp = type_8047992331715851194_FromNu(val)
	val, ok = record["all_day"]
	if !ok {
		out.AllDay = false
	} else {
		out.AllDay = type_729807561129781588_FromNu(val)
	}
	val, ok = record["floating"]
	if !ok {
		out.Floating = false
	} else {
		out.Floating = type_729807561129781588_FromNu(val)
	}
	return out
}
func type_784588192188755836_FromNu(v nu.Value) *events.EventStatus {
	if v.Value == nil {
		return nil
	}
	res := type_15385297846572725340_FromNu(v)
	return &res
}
func type_5863190983406162214_FromNu(v nu.Value) *time.Duration {
	if v.Value == nil {
		return nil
	}
	res := type_16589689216511618220_FromNu(v)
	return &res
}
func type_4202073097562803312_FromNu(v nu.Value) nutypes.TimeSegment {
	out := nutypes.TimeSegment{}
	record := v.Value.(nu.Record)
	var val nu.Value
	val, _ = record["now"]
	out.Now = type_8047992331715851194_FromNu(val)
	val, _ = record["duration"]
	out.Duration = type_16589689216511618220_FromNu(val)
	val, _ = record["active_events"]
	out.ActiveEvents = type_6714134590091148549_FromNu(val)
	return out
}
func type_17862013815172309399_FromNu(v nu.Value) *string {
	if v.Value == nil {
		return nil
	}
	res := type_15613163272824911089_FromNu(v)
	return &res
}
func type_17860233973098560385_FromNu(v nu.Value) float64 {
	return float64(v.Value.(float64))
}
func type_15560982419391353847_FromNu(v nu.Value) events.EventTriggerRelative {
	return events.EventTriggerRelative(v.Value.(int64))
}
func type_15050730807189225719_FromNu(v nu.Value) *time.Time {
	if v.Value == nil {
		return nil
	}
	res := type_8047992331715851194_FromNu(v)
	return &res
}
func type_729807561129781588_FromNu(v nu.Value) bool {
	return bool(v.Value.(bool))
}
func type_7161572108068222122_FromNu(v nu.Value) events.EventGeo {
	out := events.EventGeo{}
	record := v.Value.(nu.Record)
	var val nu.Value
	val, _ = record["latitude"]
	out.Latitude = type_17860233973098560385_FromNu(val)
	val, _ = record["longitude"]
	out.Longitude = type_17860233973098560385_FromNu(val)
	return out
}
func type_5363327835607766502_FromNu(v nu.Value) *url.URL {
	if v.Value == nil {
		return nil
	}
	parsed, err := url.Parse(v.Value.(string))
	if err != nil {
		panic(err)
	}
	return parsed
}
func type_581713733206709016_FromNu(v nu.Value) nutypes.EventObjectReplica {
	out := nutypes.EventObjectReplica{}
	record := v.Value.(nu.Record)
	var val nu.Value
	val, _ = record["object_path"]
	out.ObjectPath = type_17862013815172309399_FromNu(val)
	val, _ = record["main"]
	out.Main = type_4505817543918974569_FromNu(val)
	val, _ = record["overrides"]
	out.Overrides = type_6714134590091148549_FromNu(val)
	return out
}
func type_2392839876798024645_FromNu(v nu.Value) nutypes.EventObjectReplicaList {
	if v.Value == nil {
		return nil
	}
	arr := v.Value.([]nu.Value)
	out := make(nutypes.EventObjectReplicaList, len(arr))
	for i, e := range arr {
		out[i] = type_581713733206709016_FromNu(e)
	}
	return out
}
func type_15613163272824911089_FromNu(v nu.Value) string {
	return string(v.Value.(string))
}
func type_7163250051298988498_FromNu(v nu.Value) *events.EventGeo {
	if v.Value == nil {
		return nil
	}
	res := type_7161572108068222122_FromNu(v)
	return &res
}
func type_18334623996676649874_FromNu(v nu.Value) *rrule.RRule {
	if v.Value == nil {
		return nil
	}
	parsed, err := rrule.StrToRRule(v.Value.(string))
	if err != nil {
		panic(err)
	}
	return parsed
}
func type_3931126380996215332_FromNu(v nu.Value) []events.Datetime {
	if v.Value == nil {
		return nil
	}
	arr := v.Value.([]nu.Value)
	out := make([]events.Datetime, len(arr))
	for i, e := range arr {
		out[i] = type_5454485661162817076_FromNu(e)
	}
	return out
}
func type_6714134590091148549_FromNu(v nu.Value) []nutypes.EventReplica {
	if v.Value == nil {
		return nil
	}
	arr := v.Value.([]nu.Value)
	out := make([]nutypes.EventReplica, len(arr))
	for i, e := range arr {
		out[i] = type_4505817543918974569_FromNu(e)
	}
	return out
}
func type_15139881813094606131_FromNu(v nu.Value) int64 {
	return int64(v.Value.(int64))
}
func type_8047992331715851194_FromNu(v nu.Value) time.Time {
	return v.Value.(time.Time)
}
func type_2493169154543297135_FromNu(v nu.Value) events.EventClass {
	return events.EventClass(v.Value.(string))
}
func type_5363327835607766502_ToNu(v *url.URL) nu.Value {
	if v == nil {
		return nu.Value{Value: nil}
	}
	return nu.ToValue(v.String())
}
func type_4202073097562803312_ToNu(v nutypes.TimeSegment) nu.Value {
	return nu.Value{Value: nu.Record{
		"now":           type_8047992331715851194_ToNu(v.Now),
		"duration":      type_16589689216511618220_ToNu(v.Duration),
		"active_events": type_6714134590091148549_ToNu(v.ActiveEvents),
	}}
}
func type_8047992331715851194_ToNu(v time.Time) nu.Value {
	return nu.ToValue(v)
}
func type_2493169154543297135_ToNu(v events.EventClass) nu.Value {
	return nu.ToValue(v)
}
func type_10890016574791629639_ToNu(v int) nu.Value {
	return nu.ToValue(v)
}
func type_784588192188755836_ToNu(v *events.EventStatus) nu.Value {
	if v == nil {
		return nu.Value{Value: nil}
	}
	return type_15385297846572725340_ToNu(*v)
}
func type_18334623996676649874_ToNu(v *rrule.RRule) nu.Value {
	if v == nil {
		return nu.Value{Value: nil}
	}
	return nu.ToValue(v.String())
}
func type_15560982419391353847_ToNu(v events.EventTriggerRelative) nu.Value {
	return nu.ToValue(v)
}
func type_14900610606431710770_ToNu(v nutypes.Timeline) nu.Value {
	list := make([]nu.Value, len(v))
	for i, e := range v {
		list[i] = type_4202073097562803312_ToNu(e)
	}
	return nu.Value{Value: list}
}
func type_18369289839240265122_ToNu(v caldav.Calendar) nu.Value {
	return nu.Value{Value: nu.Record{
		"path":                    type_15613163272824911089_ToNu(v.Path),
		"name":                    type_15613163272824911089_ToNu(v.Name),
		"description":             type_15613163272824911089_ToNu(v.Description),
		"max_resource_size":       type_15139881813094606131_ToNu(v.MaxResourceSize),
		"supported_component_set": type_11669970230249425419_ToNu(v.SupportedComponentSet),
	}}
}
func type_15613163272824911089_ToNu(v string) nu.Value {
	return nu.ToValue(v)
}
func type_7057708295081751301_ToNu(v events.EventTransparency) nu.Value {
	return nu.ToValue(v)
}
func type_17862013815172309399_ToNu(v *string) nu.Value {
	if v == nil {
		return nu.Value{Value: nil}
	}
	return type_15613163272824911089_ToNu(*v)
}
func type_729807561129781588_ToNu(v bool) nu.Value {
	return nu.ToValue(v)
}
func type_17860233973098560385_ToNu(v float64) nu.Value {
	return nu.ToValue(v)
}
func type_7163250051298988498_ToNu(v *events.EventGeo) nu.Value {
	if v == nil {
		return nu.Value{Value: nil}
	}
	return type_7161572108068222122_ToNu(*v)
}
func type_9664538759823739797_ToNu(v *events.EventClass) nu.Value {
	if v == nil {
		return nu.Value{Value: nil}
	}
	return type_2493169154543297135_ToNu(*v)
}
func type_7161572108068222122_ToNu(v events.EventGeo) nu.Value {
	return nu.Value{Value: nu.Record{
		"latitude":  type_17860233973098560385_ToNu(v.Latitude),
		"longitude": type_17860233973098560385_ToNu(v.Longitude),
	}}
}
func type_16589689216511618220_ToNu(v time.Duration) nu.Value {
	return nu.ToValue(v)
}
func type_15050730807189225719_ToNu(v *time.Time) nu.Value {
	if v == nil {
		return nu.Value{Value: nil}
	}
	return type_8047992331715851194_ToNu(*v)
}
func type_9520111014888170891_ToNu(v *events.EventTrigger) nu.Value {
	if v == nil {
		return nu.Value{Value: nil}
	}
	return type_13545470577293064413_ToNu(*v)
}
func type_6714134590091148549_ToNu(v []nutypes.EventReplica) nu.Value {
	list := make([]nu.Value, len(v))
	for i, e := range v {
		list[i] = type_4505817543918974569_ToNu(e)
	}
	return nu.Value{Value: list}
}
func type_15139881813094606131_ToNu(v int64) nu.Value {
	return nu.ToValue(v)
}
func type_12480522309550428545_ToNu(v *events.Datetime) nu.Value {
	if v == nil {
		return nu.Value{Value: nil}
	}
	return type_5454485661162817076_ToNu(*v)
}
func type_3931126380996215332_ToNu(v []events.Datetime) nu.Value {
	list := make([]nu.Value, len(v))
	for i, e := range v {
		list[i] = type_5454485661162817076_ToNu(e)
	}
	return nu.Value{Value: list}
}
func type_5863190983406162214_ToNu(v *time.Duration) nu.Value {
	if v == nil {
		return nu.Value{Value: nil}
	}
	return type_16589689216511618220_ToNu(*v)
}
func type_5454485661162817076_ToNu(v events.Datetime) nu.Value {
	return nu.Value{Value: nu.Record{
		"stamp":    type_8047992331715851194_ToNu(v.Stamp),
		"all_day":  type_729807561129781588_ToNu(v.AllDay),
		"floating": type_729807561129781588_ToNu(v.Floating),
	}}
}
func type_15385297846572725340_ToNu(v events.EventStatus) nu.Value {
	return nu.ToValue(v)
}
func type_8971279483973357571_ToNu(v *events.EventTransparency) nu.Value {
	if v == nil {
		return nu.Value{Value: nil}
	}
	return type_7057708295081751301_ToNu(*v)
}
func type_4505817543918974569_ToNu(v nutypes.EventReplica) nu.Value {
	return nu.Value{Value: nu.Record{
		"uid":                        type_17862013815172309399_ToNu(v.Uid),
		"summary":                    type_17862013815172309399_ToNu(v.Summary),
		"location":                   type_17862013815172309399_ToNu(v.Location),
		"description":                type_17862013815172309399_ToNu(v.Description),
		"categories":                 type_11669970230249425419_ToNu(v.Categories),
		"created":                    type_12480522309550428545_ToNu(v.Created),
		"last_modified":              type_12480522309550428545_ToNu(v.LastModified),
		"class":                      type_9664538759823739797_ToNu(v.Class),
		"geo":                        type_7163250051298988498_ToNu(v.Geo),
		"priority":                   type_2584899110032584934_ToNu(v.Priority),
		"sequence":                   type_2584899110032584934_ToNu(v.Sequence),
		"status":                     type_784588192188755836_ToNu(v.Status),
		"transparency":               type_8971279483973357571_ToNu(v.Transparency),
		"url":                        type_5363327835607766502_ToNu(v.URL),
		"comment":                    type_17862013815172309399_ToNu(v.Comment),
		"attach":                     type_5363327835607766502_ToNu(v.Attach),
		"contact":                    type_17862013815172309399_ToNu(v.Contact),
		"organizer":                  type_5363327835607766502_ToNu(v.Organizer),
		"start":                      type_5454485661162817076_ToNu(v.Start),
		"end":                        type_5454485661162817076_ToNu(v.End),
		"recurrence_rule":            type_18334623996676649874_ToNu(v.RecurrenceRule),
		"recurrence_dates":           type_3931126380996215332_ToNu(v.RecurrenceDates),
		"recurrence_exception_dates": type_3931126380996215332_ToNu(v.RecurrenceExceptionDates),
		"recurrence_instance":        type_12480522309550428545_ToNu(v.RecurrenceInstance),
		"trigger":                    type_9520111014888170891_ToNu(v.Trigger),
	}}
}
func type_581713733206709016_ToNu(v nutypes.EventObjectReplica) nu.Value {
	return nu.Value{Value: nu.Record{
		"object_path": type_17862013815172309399_ToNu(v.ObjectPath),
		"main":        type_4505817543918974569_ToNu(v.Main),
		"overrides":   type_6714134590091148549_ToNu(v.Overrides),
	}}
}
func type_2392839876798024645_ToNu(v nutypes.EventObjectReplicaList) nu.Value {
	list := make([]nu.Value, len(v))
	for i, e := range v {
		list[i] = type_581713733206709016_ToNu(e)
	}
	return nu.Value{Value: list}
}
func type_4572434499988200822_ToNu(v nutypes.CalendarList) nu.Value {
	list := make([]nu.Value, len(v))
	for i, e := range v {
		list[i] = type_18369289839240265122_ToNu(e)
	}
	return nu.Value{Value: list}
}
func type_11669970230249425419_ToNu(v []string) nu.Value {
	list := make([]nu.Value, len(v))
	for i, e := range v {
		list[i] = type_15613163272824911089_ToNu(e)
	}
	return nu.Value{Value: list}
}
func type_2584899110032584934_ToNu(v *int) nu.Value {
	if v == nil {
		return nu.Value{Value: nil}
	}
	return type_10890016574791629639_ToNu(*v)
}
func type_13545470577293064413_ToNu(v events.EventTrigger) nu.Value {
	return nu.Value{Value: nu.Record{
		"relative":    type_5863190983406162214_ToNu(v.Relative),
		"relative_to": type_15560982419391353847_ToNu(v.RelativeTo),
		"absolute":    type_15050730807189225719_ToNu(v.Absolute),
	}}
}

var EventObjectReplicaListType = type_2392839876798024645
var EventObjectReplicaListFromNu = type_2392839876798024645_FromNu
var EventObjectReplicaListToNu = type_2392839876798024645_ToNu
var EventObjectReplicaType = type_581713733206709016
var EventObjectReplicaFromNu = type_581713733206709016_FromNu
var EventObjectReplicaToNu = type_581713733206709016_ToNu
var EventReplicaType = type_4505817543918974569
var EventReplicaFromNu = type_4505817543918974569_FromNu
var EventReplicaToNu = type_4505817543918974569_ToNu
var TimelineType = type_14900610606431710770
var TimelineFromNu = type_14900610606431710770_FromNu
var TimelineToNu = type_14900610606431710770_ToNu
var CalendarListType = type_4572434499988200822
var CalendarListFromNu = type_4572434499988200822_FromNu
var CalendarListToNu = type_4572434499988200822_ToNu
