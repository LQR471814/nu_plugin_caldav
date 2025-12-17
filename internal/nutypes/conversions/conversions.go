package conversions

import "net/url"
import "time"
import "fmt"
import "github.com/ainvaltin/nu-plugin"
import "github.com/ainvaltin/nu-plugin/types"
import "github.com/LQR471814/nu_plugin_caldav/events"
import "github.com/LQR471814/nu_plugin_caldav/internal/nutypes"
import "github.com/teambition/rrule-go"
import "github.com/emersion/go-webdav/caldav"

var type_17862013815172309399 = type_15613163272824911089

func type_17862013815172309399_FromNu(v nu.Value) (out *string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*string: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	res, err := type_15613163272824911089_FromNu(v)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func type_17862013815172309399_ToNu(v *string) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*string: %w", err)
		}
	}()
	if v == nil {
		return nu.Value{}, nil
	}
	return type_15613163272824911089_ToNu(*v)
}

var type_17860233973098560385 = types.Float()

func type_17860233973098560385_FromNu(v nu.Value) (out float64, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("float64: %w", err)
		}
	}()
	casted, ok := v.Value.(float64)
	converted := float64(casted)
	if !ok {
		return converted, fmt.Errorf("expected float64 got %v", v.Value)
	}
	return converted, nil
}
func type_17860233973098560385_ToNu(v float64) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("float64: %w", err)
		}
	}()
	return nu.ToValue(v), nil
}

var type_8971279483973357571 = type_7057708295081751301

func type_8971279483973357571_FromNu(v nu.Value) (out *events.EventTransparency, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*events.EventTransparency: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	res, err := type_7057708295081751301_FromNu(v)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func type_8971279483973357571_ToNu(v *events.EventTransparency) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*events.EventTransparency: %w", err)
		}
	}()
	if v == nil {
		return nu.Value{}, nil
	}
	return type_7057708295081751301_ToNu(*v)
}

var type_18334623996676649874 = types.String()

func type_18334623996676649874_FromNu(v nu.Value) (out *rrule.RRule, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*rrule.RRule: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	parsed, err := rrule.StrToRRule(v.Value.(string))
	if err != nil {
		return nil, err
	}
	return parsed, nil
}
func type_18334623996676649874_ToNu(v *rrule.RRule) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*rrule.RRule: %w", err)
		}
	}()
	if v == nil {
		return nu.Value{Value: nil}, nil
	}
	return nu.ToValue(v.String()), nil
}

var type_4202073097562803312 = types.RecordDef{
	"now":           type_8047992331715851194,
	"duration":      type_16589689216511618220,
	"active_events": type_6714134590091148549,
}

func type_4202073097562803312_FromNu(v nu.Value) (out nutypes.TimeSegment, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("nutypes.TimeSegment: %w", err)
		}
	}()
	record, ok := v.Value.(nu.Record)
	if !ok {
		return out, fmt.Errorf("expected nu.Record got %T", v.Value)
	}
	var val nu.Value
	val, _ = record["now"]
	out.Now, err = type_8047992331715851194_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["duration"]
	out.Duration, err = type_16589689216511618220_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["active_events"]
	out.ActiveEvents, err = type_6714134590091148549_FromNu(val)
	if err != nil {
		return out, err
	}
	return out, nil
}
func type_4202073097562803312_ToNu(v nutypes.TimeSegment) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("nutypes.TimeSegment: %w", err)
		}
	}()
	rec := nu.Record{}
	rec["now"], err = type_8047992331715851194_ToNu(v.Now)
	if err != nil {
		return nu.Value{}, err
	}
	rec["duration"], err = type_16589689216511618220_ToNu(v.Duration)
	if err != nil {
		return nu.Value{}, err
	}
	rec["active_events"], err = type_6714134590091148549_ToNu(v.ActiveEvents)
	if err != nil {
		return nu.Value{}, err
	}
	return nu.Value{Value: rec}, nil
}

var type_15139881813094606131 = types.Int()

func type_15139881813094606131_FromNu(v nu.Value) (out int64, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("int64: %w", err)
		}
	}()
	casted, ok := v.Value.(int64)
	converted := int64(casted)
	if !ok {
		return converted, fmt.Errorf("expected int64 got %v", v.Value)
	}
	return converted, nil
}
func type_15139881813094606131_ToNu(v int64) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("int64: %w", err)
		}
	}()
	return nu.ToValue(v), nil
}

var type_4572434499988200822 = types.Table(type_18369289839240265122)

func type_4572434499988200822_FromNu(v nu.Value) (out nutypes.CalendarList, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("nutypes.CalendarList: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	arr, ok := v.Value.([]nu.Value)
	if !ok {
		return nil, fmt.Errorf("expected []nu.Value got %T", v.Value)
	}
	out = make(nutypes.CalendarList, len(arr))
	for i, e := range arr {
		out[i], err = type_18369289839240265122_FromNu(e)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
func type_4572434499988200822_ToNu(v nutypes.CalendarList) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("nutypes.CalendarList: %w", err)
		}
	}()
	list := make([]nu.Value, len(v))
	for i, e := range v {
		list[i], err = type_18369289839240265122_ToNu(e)
		if err != nil {
			return nu.Value{}, err
		}
	}
	return nu.Value{Value: list}, nil
}

var type_15613163272824911089 = types.String()

func type_15613163272824911089_FromNu(v nu.Value) (out string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("string: %w", err)
		}
	}()
	casted, ok := v.Value.(string)
	converted := string(casted)
	if !ok {
		return converted, fmt.Errorf("expected string got %v", v.Value)
	}
	return converted, nil
}
func type_15613163272824911089_ToNu(v string) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("string: %w", err)
		}
	}()
	return nu.ToValue(v), nil
}

var type_8047992331715851194 = types.Date()

func type_8047992331715851194_FromNu(v nu.Value) (out time.Time, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("time.Time: %w", err)
		}
	}()
	out, ok := v.Value.(time.Time)
	if !ok {
		return out, fmt.Errorf("expected time.Time got %T", v.Value)
	}
	return
}
func type_8047992331715851194_ToNu(v time.Time) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("time.Time: %w", err)
		}
	}()
	return nu.ToValue(v), nil
}

var type_9664538759823739797 = type_2493169154543297135

func type_9664538759823739797_FromNu(v nu.Value) (out *events.EventClass, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*events.EventClass: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	res, err := type_2493169154543297135_FromNu(v)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func type_9664538759823739797_ToNu(v *events.EventClass) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*events.EventClass: %w", err)
		}
	}()
	if v == nil {
		return nu.Value{}, nil
	}
	return type_2493169154543297135_ToNu(*v)
}

var type_7163250051298988498 = types.Record(type_7161572108068222122)

func type_7163250051298988498_FromNu(v nu.Value) (out *events.EventGeo, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*events.EventGeo: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	res, err := type_7161572108068222122_FromNu(v)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func type_7163250051298988498_ToNu(v *events.EventGeo) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*events.EventGeo: %w", err)
		}
	}()
	if v == nil {
		return nu.Value{}, nil
	}
	return type_7161572108068222122_ToNu(*v)
}

var type_2584899110032584934 = type_10890016574791629639

func type_2584899110032584934_FromNu(v nu.Value) (out *int, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*int: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	res, err := type_10890016574791629639_FromNu(v)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func type_2584899110032584934_ToNu(v *int) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*int: %w", err)
		}
	}()
	if v == nil {
		return nu.Value{}, nil
	}
	return type_10890016574791629639_ToNu(*v)
}

var type_3931126380996215332 = types.Table(type_5454485661162817076)

func type_3931126380996215332_FromNu(v nu.Value) (out []events.Datetime, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("[]events.Datetime: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	arr, ok := v.Value.([]nu.Value)
	if !ok {
		return nil, fmt.Errorf("expected []nu.Value got %T", v.Value)
	}
	out = make([]events.Datetime, len(arr))
	for i, e := range arr {
		out[i], err = type_5454485661162817076_FromNu(e)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
func type_3931126380996215332_ToNu(v []events.Datetime) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("[]events.Datetime: %w", err)
		}
	}()
	list := make([]nu.Value, len(v))
	for i, e := range v {
		list[i], err = type_5454485661162817076_ToNu(e)
		if err != nil {
			return nu.Value{}, err
		}
	}
	return nu.Value{Value: list}, nil
}

var type_5863190983406162214 = type_16589689216511618220

func type_5863190983406162214_FromNu(v nu.Value) (out *time.Duration, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*time.Duration: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	res, err := type_16589689216511618220_FromNu(v)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func type_5863190983406162214_ToNu(v *time.Duration) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*time.Duration: %w", err)
		}
	}()
	if v == nil {
		return nu.Value{}, nil
	}
	return type_16589689216511618220_ToNu(*v)
}

var type_4505817543918974569 = types.RecordDef{
	"uid":                        type_17862013815172309399,
	"summary":                    type_17862013815172309399,
	"location":                   type_17862013815172309399,
	"description":                type_17862013815172309399,
	"categories":                 type_11669970230249425419,
	"datetime_stamp":             type_12480522309550428545,
	"created":                    type_12480522309550428545,
	"last_modified":              type_12480522309550428545,
	"class":                      type_9664538759823739797,
	"geo":                        type_7163250051298988498,
	"priority":                   type_2584899110032584934,
	"sequence":                   type_2584899110032584934,
	"status":                     type_784588192188755836,
	"transparency":               type_8971279483973357571,
	"url":                        type_5363327835607766502,
	"comment":                    type_17862013815172309399,
	"attach":                     type_5363327835607766502,
	"contact":                    type_17862013815172309399,
	"organizer":                  type_5363327835607766502,
	"start":                      types.Record(type_5454485661162817076),
	"end":                        types.Record(type_5454485661162817076),
	"recurrence_rule":            type_18334623996676649874,
	"recurrence_dates":           type_3931126380996215332,
	"recurrence_exception_dates": type_3931126380996215332,
	"recurrence_instance":        type_12480522309550428545,
	"trigger":                    type_9520111014888170891,
	"other":                      type_8207267565485682006,
}

func type_4505817543918974569_FromNu(v nu.Value) (out nutypes.EventReplica, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("nutypes.EventReplica: %w", err)
		}
	}()
	record, ok := v.Value.(nu.Record)
	if !ok {
		return out, fmt.Errorf("expected nu.Record got %T", v.Value)
	}
	var val nu.Value
	val, _ = record["uid"]
	out.Uid, err = type_17862013815172309399_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["summary"]
	out.Summary, err = type_17862013815172309399_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["location"]
	out.Location, err = type_17862013815172309399_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["description"]
	out.Description, err = type_17862013815172309399_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["categories"]
	out.Categories, err = type_11669970230249425419_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["datetime_stamp"]
	out.DatetimeStamp, err = type_12480522309550428545_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["created"]
	out.Created, err = type_12480522309550428545_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["last_modified"]
	out.LastModified, err = type_12480522309550428545_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["class"]
	out.Class, err = type_9664538759823739797_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["geo"]
	out.Geo, err = type_7163250051298988498_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["priority"]
	out.Priority, err = type_2584899110032584934_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["sequence"]
	out.Sequence, err = type_2584899110032584934_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["status"]
	out.Status, err = type_784588192188755836_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["transparency"]
	out.Transparency, err = type_8971279483973357571_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["url"]
	out.URL, err = type_5363327835607766502_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["comment"]
	out.Comment, err = type_17862013815172309399_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["attach"]
	out.Attach, err = type_5363327835607766502_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["contact"]
	out.Contact, err = type_17862013815172309399_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["organizer"]
	out.Organizer, err = type_5363327835607766502_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["start"]
	out.Start, err = type_5454485661162817076_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["end"]
	out.End, err = type_5454485661162817076_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["recurrence_rule"]
	out.RecurrenceRule, err = type_18334623996676649874_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["recurrence_dates"]
	out.RecurrenceDates, err = type_3931126380996215332_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["recurrence_exception_dates"]
	out.RecurrenceExceptionDates, err = type_3931126380996215332_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["recurrence_instance"]
	out.RecurrenceInstance, err = type_12480522309550428545_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["trigger"]
	out.Trigger, err = type_9520111014888170891_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["other"]
	out.Other, err = type_8207267565485682006_FromNu(val)
	if err != nil {
		return out, err
	}
	return out, nil
}
func type_4505817543918974569_ToNu(v nutypes.EventReplica) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("nutypes.EventReplica: %w", err)
		}
	}()
	rec := nu.Record{}
	rec["uid"], err = type_17862013815172309399_ToNu(v.Uid)
	if err != nil {
		return nu.Value{}, err
	}
	rec["summary"], err = type_17862013815172309399_ToNu(v.Summary)
	if err != nil {
		return nu.Value{}, err
	}
	rec["location"], err = type_17862013815172309399_ToNu(v.Location)
	if err != nil {
		return nu.Value{}, err
	}
	rec["description"], err = type_17862013815172309399_ToNu(v.Description)
	if err != nil {
		return nu.Value{}, err
	}
	rec["categories"], err = type_11669970230249425419_ToNu(v.Categories)
	if err != nil {
		return nu.Value{}, err
	}
	rec["datetime_stamp"], err = type_12480522309550428545_ToNu(v.DatetimeStamp)
	if err != nil {
		return nu.Value{}, err
	}
	rec["created"], err = type_12480522309550428545_ToNu(v.Created)
	if err != nil {
		return nu.Value{}, err
	}
	rec["last_modified"], err = type_12480522309550428545_ToNu(v.LastModified)
	if err != nil {
		return nu.Value{}, err
	}
	rec["class"], err = type_9664538759823739797_ToNu(v.Class)
	if err != nil {
		return nu.Value{}, err
	}
	rec["geo"], err = type_7163250051298988498_ToNu(v.Geo)
	if err != nil {
		return nu.Value{}, err
	}
	rec["priority"], err = type_2584899110032584934_ToNu(v.Priority)
	if err != nil {
		return nu.Value{}, err
	}
	rec["sequence"], err = type_2584899110032584934_ToNu(v.Sequence)
	if err != nil {
		return nu.Value{}, err
	}
	rec["status"], err = type_784588192188755836_ToNu(v.Status)
	if err != nil {
		return nu.Value{}, err
	}
	rec["transparency"], err = type_8971279483973357571_ToNu(v.Transparency)
	if err != nil {
		return nu.Value{}, err
	}
	rec["url"], err = type_5363327835607766502_ToNu(v.URL)
	if err != nil {
		return nu.Value{}, err
	}
	rec["comment"], err = type_17862013815172309399_ToNu(v.Comment)
	if err != nil {
		return nu.Value{}, err
	}
	rec["attach"], err = type_5363327835607766502_ToNu(v.Attach)
	if err != nil {
		return nu.Value{}, err
	}
	rec["contact"], err = type_17862013815172309399_ToNu(v.Contact)
	if err != nil {
		return nu.Value{}, err
	}
	rec["organizer"], err = type_5363327835607766502_ToNu(v.Organizer)
	if err != nil {
		return nu.Value{}, err
	}
	rec["start"], err = type_5454485661162817076_ToNu(v.Start)
	if err != nil {
		return nu.Value{}, err
	}
	rec["end"], err = type_5454485661162817076_ToNu(v.End)
	if err != nil {
		return nu.Value{}, err
	}
	rec["recurrence_rule"], err = type_18334623996676649874_ToNu(v.RecurrenceRule)
	if err != nil {
		return nu.Value{}, err
	}
	rec["recurrence_dates"], err = type_3931126380996215332_ToNu(v.RecurrenceDates)
	if err != nil {
		return nu.Value{}, err
	}
	rec["recurrence_exception_dates"], err = type_3931126380996215332_ToNu(v.RecurrenceExceptionDates)
	if err != nil {
		return nu.Value{}, err
	}
	rec["recurrence_instance"], err = type_12480522309550428545_ToNu(v.RecurrenceInstance)
	if err != nil {
		return nu.Value{}, err
	}
	rec["trigger"], err = type_9520111014888170891_ToNu(v.Trigger)
	if err != nil {
		return nu.Value{}, err
	}
	rec["other"], err = type_8207267565485682006_ToNu(v.Other)
	if err != nil {
		return nu.Value{}, err
	}
	return nu.Value{Value: rec}, nil
}

var type_2493169154543297135 = types.String()

func type_2493169154543297135_FromNu(v nu.Value) (out events.EventClass, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.EventClass: %w", err)
		}
	}()
	casted, ok := v.Value.(string)
	converted := events.EventClass(casted)
	if !ok {
		return converted, fmt.Errorf("expected string got %v", v.Value)
	}
	return converted, nil
}
func type_2493169154543297135_ToNu(v events.EventClass) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.EventClass: %w", err)
		}
	}()
	return nu.ToValue(v), nil
}

var type_784588192188755836 = type_15385297846572725340

func type_784588192188755836_FromNu(v nu.Value) (out *events.EventStatus, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*events.EventStatus: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	res, err := type_15385297846572725340_FromNu(v)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func type_784588192188755836_ToNu(v *events.EventStatus) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*events.EventStatus: %w", err)
		}
	}()
	if v == nil {
		return nu.Value{}, nil
	}
	return type_15385297846572725340_ToNu(*v)
}

var type_15560982419391353847 = types.Int()

func type_15560982419391353847_FromNu(v nu.Value) (out events.EventTriggerRelative, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.EventTriggerRelative: %w", err)
		}
	}()
	casted, ok := v.Value.(int64)
	converted := events.EventTriggerRelative(casted)
	if !ok {
		return converted, fmt.Errorf("expected int64 got %v", v.Value)
	}
	return converted, nil
}
func type_15560982419391353847_ToNu(v events.EventTriggerRelative) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.EventTriggerRelative: %w", err)
		}
	}()
	return nu.ToValue(v), nil
}

var type_9520111014888170891 = types.Record(type_13545470577293064413)

func type_9520111014888170891_FromNu(v nu.Value) (out *events.EventTrigger, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*events.EventTrigger: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	res, err := type_13545470577293064413_FromNu(v)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func type_9520111014888170891_ToNu(v *events.EventTrigger) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*events.EventTrigger: %w", err)
		}
	}()
	if v == nil {
		return nu.Value{}, nil
	}
	return type_13545470577293064413_ToNu(*v)
}

var type_581713733206709016 = types.RecordDef{
	"object_path": type_17862013815172309399,
	"main":        types.Record(type_4505817543918974569),
	"overrides":   type_6714134590091148549,
}

func type_581713733206709016_FromNu(v nu.Value) (out nutypes.EventObjectReplica, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("nutypes.EventObjectReplica: %w", err)
		}
	}()
	record, ok := v.Value.(nu.Record)
	if !ok {
		return out, fmt.Errorf("expected nu.Record got %T", v.Value)
	}
	var val nu.Value
	val, _ = record["object_path"]
	out.ObjectPath, err = type_17862013815172309399_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["main"]
	out.Main, err = type_4505817543918974569_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["overrides"]
	out.Overrides, err = type_6714134590091148549_FromNu(val)
	if err != nil {
		return out, err
	}
	return out, nil
}
func type_581713733206709016_ToNu(v nutypes.EventObjectReplica) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("nutypes.EventObjectReplica: %w", err)
		}
	}()
	rec := nu.Record{}
	rec["object_path"], err = type_17862013815172309399_ToNu(v.ObjectPath)
	if err != nil {
		return nu.Value{}, err
	}
	rec["main"], err = type_4505817543918974569_ToNu(v.Main)
	if err != nil {
		return nu.Value{}, err
	}
	rec["overrides"], err = type_6714134590091148549_ToNu(v.Overrides)
	if err != nil {
		return nu.Value{}, err
	}
	return nu.Value{Value: rec}, nil
}

var type_18369289839240265122 = types.RecordDef{
	"path":                    type_15613163272824911089,
	"name":                    type_15613163272824911089,
	"description":             type_15613163272824911089,
	"max_resource_size":       type_15139881813094606131,
	"supported_component_set": type_11669970230249425419,
}

func type_18369289839240265122_FromNu(v nu.Value) (out caldav.Calendar, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("caldav.Calendar: %w", err)
		}
	}()
	record, ok := v.Value.(nu.Record)
	if !ok {
		return out, fmt.Errorf("expected nu.Record got %T", v.Value)
	}
	var val nu.Value
	val, _ = record["path"]
	out.Path, err = type_15613163272824911089_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["name"]
	out.Name, err = type_15613163272824911089_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["description"]
	out.Description, err = type_15613163272824911089_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["max_resource_size"]
	out.MaxResourceSize, err = type_15139881813094606131_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["supported_component_set"]
	out.SupportedComponentSet, err = type_11669970230249425419_FromNu(val)
	if err != nil {
		return out, err
	}
	return out, nil
}
func type_18369289839240265122_ToNu(v caldav.Calendar) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("caldav.Calendar: %w", err)
		}
	}()
	rec := nu.Record{}
	rec["path"], err = type_15613163272824911089_ToNu(v.Path)
	if err != nil {
		return nu.Value{}, err
	}
	rec["name"], err = type_15613163272824911089_ToNu(v.Name)
	if err != nil {
		return nu.Value{}, err
	}
	rec["description"], err = type_15613163272824911089_ToNu(v.Description)
	if err != nil {
		return nu.Value{}, err
	}
	rec["max_resource_size"], err = type_15139881813094606131_ToNu(v.MaxResourceSize)
	if err != nil {
		return nu.Value{}, err
	}
	rec["supported_component_set"], err = type_11669970230249425419_ToNu(v.SupportedComponentSet)
	if err != nil {
		return nu.Value{}, err
	}
	return nu.Value{Value: rec}, nil
}

var type_11669970230249425419 = types.List(type_15613163272824911089)

func type_11669970230249425419_FromNu(v nu.Value) (out []string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("[]string: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	arr, ok := v.Value.([]nu.Value)
	if !ok {
		return nil, fmt.Errorf("expected []nu.Value got %T", v.Value)
	}
	out = make([]string, len(arr))
	for i, e := range arr {
		out[i], err = type_15613163272824911089_FromNu(e)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
func type_11669970230249425419_ToNu(v []string) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("[]string: %w", err)
		}
	}()
	list := make([]nu.Value, len(v))
	for i, e := range v {
		list[i], err = type_15613163272824911089_ToNu(e)
		if err != nil {
			return nu.Value{}, err
		}
	}
	return nu.Value{Value: list}, nil
}

var type_15385297846572725340 = types.String()

func type_15385297846572725340_FromNu(v nu.Value) (out events.EventStatus, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.EventStatus: %w", err)
		}
	}()
	casted, ok := v.Value.(string)
	converted := events.EventStatus(casted)
	if !ok {
		return converted, fmt.Errorf("expected string got %v", v.Value)
	}
	return converted, nil
}
func type_15385297846572725340_ToNu(v events.EventStatus) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.EventStatus: %w", err)
		}
	}()
	return nu.ToValue(v), nil
}

var type_8207267565485682006 = types.Any()

func type_8207267565485682006_FromNu(v nu.Value) (out map[string][]nutypes.PropValueReplica, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("map[string][]nutypes.PropValueReplica: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	dict, ok := v.Value.(nu.Record)
	if !ok {
		return nil, fmt.Errorf("expected nu.Record got %T", v.Value)
	}
	out = make(map[string][]nutypes.PropValueReplica, len(dict))
	for k, v := range dict {
		out[k], err = type_15678305217965743770_FromNu(v)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
func type_8207267565485682006_ToNu(v map[string][]nutypes.PropValueReplica) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("map[string][]nutypes.PropValueReplica: %w", err)
		}
	}()
	dict := make(nu.Record, len(v))
	for k, v := range v {
		dict[k], err = type_15678305217965743770_ToNu(v)
		if err != nil {
			return nu.Value{}, err
		}
	}
	return nu.Value{Value: dict}, nil
}

var type_14900610606431710770 = types.Table(type_4202073097562803312)

func type_14900610606431710770_FromNu(v nu.Value) (out nutypes.Timeline, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("nutypes.Timeline: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	arr, ok := v.Value.([]nu.Value)
	if !ok {
		return nil, fmt.Errorf("expected []nu.Value got %T", v.Value)
	}
	out = make(nutypes.Timeline, len(arr))
	for i, e := range arr {
		out[i], err = type_4202073097562803312_FromNu(e)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
func type_14900610606431710770_ToNu(v nutypes.Timeline) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("nutypes.Timeline: %w", err)
		}
	}()
	list := make([]nu.Value, len(v))
	for i, e := range v {
		list[i], err = type_4202073097562803312_ToNu(e)
		if err != nil {
			return nu.Value{}, err
		}
	}
	return nu.Value{Value: list}, nil
}

var type_5454485661162817076 = types.RecordDef{
	"stamp":    type_8047992331715851194,
	"all_day":  type_729807561129781588,
	"floating": type_729807561129781588,
}

func type_5454485661162817076_FromNu(v nu.Value) (out events.Datetime, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.Datetime: %w", err)
		}
	}()
	record, ok := v.Value.(nu.Record)
	if !ok {
		return out, fmt.Errorf("expected nu.Record got %T", v.Value)
	}
	var val nu.Value
	val, _ = record["stamp"]
	out.Stamp, err = type_8047992331715851194_FromNu(val)
	if err != nil {
		return out, err
	}
	val, ok = record["all_day"]
	if !ok {
		out.AllDay = false
	} else {
		out.AllDay, err = type_729807561129781588_FromNu(val)
		if err != nil {
			return out, err
		}
	}
	val, ok = record["floating"]
	if !ok {
		out.Floating = false
	} else {
		out.Floating, err = type_729807561129781588_FromNu(val)
		if err != nil {
			return out, err
		}
	}
	return out, nil
}
func type_5454485661162817076_ToNu(v events.Datetime) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.Datetime: %w", err)
		}
	}()
	rec := nu.Record{}
	rec["stamp"], err = type_8047992331715851194_ToNu(v.Stamp)
	if err != nil {
		return nu.Value{}, err
	}
	rec["all_day"], err = type_729807561129781588_ToNu(v.AllDay)
	if err != nil {
		return nu.Value{}, err
	}
	rec["floating"], err = type_729807561129781588_ToNu(v.Floating)
	if err != nil {
		return nu.Value{}, err
	}
	return nu.Value{Value: rec}, nil
}

var type_15050730807189225719 = type_8047992331715851194

func type_15050730807189225719_FromNu(v nu.Value) (out *time.Time, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*time.Time: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	res, err := type_8047992331715851194_FromNu(v)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func type_15050730807189225719_ToNu(v *time.Time) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*time.Time: %w", err)
		}
	}()
	if v == nil {
		return nu.Value{}, nil
	}
	return type_8047992331715851194_ToNu(*v)
}

var type_13545470577293064413 = types.RecordDef{
	"relative":    type_5863190983406162214,
	"relative_to": type_15560982419391353847,
	"absolute":    type_15050730807189225719,
}

func type_13545470577293064413_FromNu(v nu.Value) (out events.EventTrigger, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.EventTrigger: %w", err)
		}
	}()
	record, ok := v.Value.(nu.Record)
	if !ok {
		return out, fmt.Errorf("expected nu.Record got %T", v.Value)
	}
	var val nu.Value
	val, _ = record["relative"]
	out.Relative, err = type_5863190983406162214_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["relative_to"]
	out.RelativeTo, err = type_15560982419391353847_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["absolute"]
	out.Absolute, err = type_15050730807189225719_FromNu(val)
	if err != nil {
		return out, err
	}
	return out, nil
}
func type_13545470577293064413_ToNu(v events.EventTrigger) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.EventTrigger: %w", err)
		}
	}()
	rec := nu.Record{}
	rec["relative"], err = type_5863190983406162214_ToNu(v.Relative)
	if err != nil {
		return nu.Value{}, err
	}
	rec["relative_to"], err = type_15560982419391353847_ToNu(v.RelativeTo)
	if err != nil {
		return nu.Value{}, err
	}
	rec["absolute"], err = type_15050730807189225719_ToNu(v.Absolute)
	if err != nil {
		return nu.Value{}, err
	}
	return nu.Value{Value: rec}, nil
}

var type_2040221911794545828 = types.RecordDef{
	"value":  type_15613163272824911089,
	"params": type_14293658896741725053,
}

func type_2040221911794545828_FromNu(v nu.Value) (out nutypes.PropValueReplica, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("nutypes.PropValueReplica: %w", err)
		}
	}()
	record, ok := v.Value.(nu.Record)
	if !ok {
		return out, fmt.Errorf("expected nu.Record got %T", v.Value)
	}
	var val nu.Value
	val, _ = record["value"]
	out.Value, err = type_15613163272824911089_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["params"]
	out.Params, err = type_14293658896741725053_FromNu(val)
	if err != nil {
		return out, err
	}
	return out, nil
}
func type_2040221911794545828_ToNu(v nutypes.PropValueReplica) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("nutypes.PropValueReplica: %w", err)
		}
	}()
	rec := nu.Record{}
	rec["value"], err = type_15613163272824911089_ToNu(v.Value)
	if err != nil {
		return nu.Value{}, err
	}
	rec["params"], err = type_14293658896741725053_ToNu(v.Params)
	if err != nil {
		return nu.Value{}, err
	}
	return nu.Value{Value: rec}, nil
}

var type_6714134590091148549 = types.Table(type_4505817543918974569)

func type_6714134590091148549_FromNu(v nu.Value) (out []nutypes.EventReplica, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("[]nutypes.EventReplica: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	arr, ok := v.Value.([]nu.Value)
	if !ok {
		return nil, fmt.Errorf("expected []nu.Value got %T", v.Value)
	}
	out = make([]nutypes.EventReplica, len(arr))
	for i, e := range arr {
		out[i], err = type_4505817543918974569_FromNu(e)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
func type_6714134590091148549_ToNu(v []nutypes.EventReplica) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("[]nutypes.EventReplica: %w", err)
		}
	}()
	list := make([]nu.Value, len(v))
	for i, e := range v {
		list[i], err = type_4505817543918974569_ToNu(e)
		if err != nil {
			return nu.Value{}, err
		}
	}
	return nu.Value{Value: list}, nil
}

var type_729807561129781588 = types.Bool()

func type_729807561129781588_FromNu(v nu.Value) (out bool, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("bool: %w", err)
		}
	}()
	casted, ok := v.Value.(bool)
	converted := bool(casted)
	if !ok {
		return converted, fmt.Errorf("expected bool got %v", v.Value)
	}
	return converted, nil
}
func type_729807561129781588_ToNu(v bool) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("bool: %w", err)
		}
	}()
	return nu.ToValue(v), nil
}

var type_5363327835607766502 = types.String()

func type_5363327835607766502_FromNu(v nu.Value) (out *url.URL, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*url.URL: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	parsed, err := url.Parse(v.Value.(string))
	if err != nil {
		return nil, err
	}
	return parsed, nil
}
func type_5363327835607766502_ToNu(v *url.URL) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*url.URL: %w", err)
		}
	}()
	if v == nil {
		return nu.Value{Value: nil}, nil
	}
	return nu.ToValue(v.String()), nil
}

var type_16589689216511618220 = types.Duration()

func type_16589689216511618220_FromNu(v nu.Value) (out time.Duration, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("time.Duration: %w", err)
		}
	}()
	out, ok := v.Value.(time.Duration)
	if !ok {
		return out, fmt.Errorf("expected time.Duration got %T", v.Value)
	}
	return
}
func type_16589689216511618220_ToNu(v time.Duration) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("time.Duration: %w", err)
		}
	}()
	return nu.ToValue(v), nil
}

var type_15678305217965743770 = types.Table(type_2040221911794545828)

func type_15678305217965743770_FromNu(v nu.Value) (out []nutypes.PropValueReplica, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("[]nutypes.PropValueReplica: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	arr, ok := v.Value.([]nu.Value)
	if !ok {
		return nil, fmt.Errorf("expected []nu.Value got %T", v.Value)
	}
	out = make([]nutypes.PropValueReplica, len(arr))
	for i, e := range arr {
		out[i], err = type_2040221911794545828_FromNu(e)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
func type_15678305217965743770_ToNu(v []nutypes.PropValueReplica) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("[]nutypes.PropValueReplica: %w", err)
		}
	}()
	list := make([]nu.Value, len(v))
	for i, e := range v {
		list[i], err = type_2040221911794545828_ToNu(e)
		if err != nil {
			return nu.Value{}, err
		}
	}
	return nu.Value{Value: list}, nil
}

var type_7161572108068222122 = types.RecordDef{
	"latitude":  type_17860233973098560385,
	"longitude": type_17860233973098560385,
}

func type_7161572108068222122_FromNu(v nu.Value) (out events.EventGeo, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.EventGeo: %w", err)
		}
	}()
	record, ok := v.Value.(nu.Record)
	if !ok {
		return out, fmt.Errorf("expected nu.Record got %T", v.Value)
	}
	var val nu.Value
	val, _ = record["latitude"]
	out.Latitude, err = type_17860233973098560385_FromNu(val)
	if err != nil {
		return out, err
	}
	val, _ = record["longitude"]
	out.Longitude, err = type_17860233973098560385_FromNu(val)
	if err != nil {
		return out, err
	}
	return out, nil
}
func type_7161572108068222122_ToNu(v events.EventGeo) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.EventGeo: %w", err)
		}
	}()
	rec := nu.Record{}
	rec["latitude"], err = type_17860233973098560385_ToNu(v.Latitude)
	if err != nil {
		return nu.Value{}, err
	}
	rec["longitude"], err = type_17860233973098560385_ToNu(v.Longitude)
	if err != nil {
		return nu.Value{}, err
	}
	return nu.Value{Value: rec}, nil
}

var type_10890016574791629639 = types.Int()

func type_10890016574791629639_FromNu(v nu.Value) (out int, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("int: %w", err)
		}
	}()
	casted, ok := v.Value.(int64)
	converted := int(casted)
	if !ok {
		return converted, fmt.Errorf("expected int64 got %v", v.Value)
	}
	return converted, nil
}
func type_10890016574791629639_ToNu(v int) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("int: %w", err)
		}
	}()
	return nu.ToValue(v), nil
}

var type_7057708295081751301 = types.String()

func type_7057708295081751301_FromNu(v nu.Value) (out events.EventTransparency, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.EventTransparency: %w", err)
		}
	}()
	casted, ok := v.Value.(string)
	converted := events.EventTransparency(casted)
	if !ok {
		return converted, fmt.Errorf("expected string got %v", v.Value)
	}
	return converted, nil
}
func type_7057708295081751301_ToNu(v events.EventTransparency) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("events.EventTransparency: %w", err)
		}
	}()
	return nu.ToValue(v), nil
}

var type_14293658896741725053 = types.Any()

func type_14293658896741725053_FromNu(v nu.Value) (out map[string][]string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("map[string][]string: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	dict, ok := v.Value.(nu.Record)
	if !ok {
		return nil, fmt.Errorf("expected nu.Record got %T", v.Value)
	}
	out = make(map[string][]string, len(dict))
	for k, v := range dict {
		out[k], err = type_11669970230249425419_FromNu(v)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
func type_14293658896741725053_ToNu(v map[string][]string) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("map[string][]string: %w", err)
		}
	}()
	dict := make(nu.Record, len(v))
	for k, v := range v {
		dict[k], err = type_11669970230249425419_ToNu(v)
		if err != nil {
			return nu.Value{}, err
		}
	}
	return nu.Value{Value: dict}, nil
}

var type_12480522309550428545 = types.Record(type_5454485661162817076)

func type_12480522309550428545_FromNu(v nu.Value) (out *events.Datetime, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*events.Datetime: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	res, err := type_5454485661162817076_FromNu(v)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
func type_12480522309550428545_ToNu(v *events.Datetime) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("*events.Datetime: %w", err)
		}
	}()
	if v == nil {
		return nu.Value{}, nil
	}
	return type_5454485661162817076_ToNu(*v)
}

var type_2392839876798024645 = types.Table(type_581713733206709016)

func type_2392839876798024645_FromNu(v nu.Value) (out nutypes.EventObjectReplicaList, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("nutypes.EventObjectReplicaList: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	arr, ok := v.Value.([]nu.Value)
	if !ok {
		return nil, fmt.Errorf("expected []nu.Value got %T", v.Value)
	}
	out = make(nutypes.EventObjectReplicaList, len(arr))
	for i, e := range arr {
		out[i], err = type_581713733206709016_FromNu(e)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
func type_2392839876798024645_ToNu(v nutypes.EventObjectReplicaList) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("nutypes.EventObjectReplicaList: %w", err)
		}
	}()
	list := make([]nu.Value, len(v))
	for i, e := range v {
		list[i], err = type_581713733206709016_ToNu(e)
		if err != nil {
			return nu.Value{}, err
		}
	}
	return nu.Value{Value: list}, nil
}

var EventReplicaType = type_4505817543918974569
var EventReplicaFromNu = type_4505817543918974569_FromNu
var EventReplicaToNu = type_4505817543918974569_ToNu
var TimelineType = type_14900610606431710770
var TimelineFromNu = type_14900610606431710770_FromNu
var TimelineToNu = type_14900610606431710770_ToNu
var CalendarListType = type_4572434499988200822
var CalendarListFromNu = type_4572434499988200822_FromNu
var CalendarListToNu = type_4572434499988200822_ToNu
var EventObjectReplicaListType = type_2392839876798024645
var EventObjectReplicaListFromNu = type_2392839876798024645_FromNu
var EventObjectReplicaListToNu = type_2392839876798024645_ToNu
var EventObjectReplicaType = type_581713733206709016
var EventObjectReplicaFromNu = type_581713733206709016_FromNu
var EventObjectReplicaToNu = type_581713733206709016_ToNu
