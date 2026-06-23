package nuconv

import "time"
import "fmt"
import "github.com/ainvaltin/nu-plugin"
import "github.com/ainvaltin/nu-plugin/types"
import "github.com/LQR471814/nu_plugin_caldav/internal/enrich/props"
import "github.com/LQR471814/nu_plugin_caldav/internal/enrich/dto"
import "github.com/teambition/rrule-go"
import "github.com/emersion/go-webdav/caldav"

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

var type_11339315157932765317 = types.RecordDef{
	"relative":    type_5863190983406162214,
	"relative_to": type_13250047047225666367,
	"absolute":    type_15050730807189225719,
}

func type_11339315157932765317_FromNu(v nu.Value) (out props.EventTrigger, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("props.EventTrigger: %w", err)
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
	out.RelativeTo, err = type_13250047047225666367_FromNu(val)
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
func type_11339315157932765317_ToNu(v props.EventTrigger) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("props.EventTrigger: %w", err)
		}
	}()
	rec := nu.Record{}
	rec["relative"], err = type_5863190983406162214_ToNu(v.Relative)
	if err != nil {
		return nu.Value{}, err
	}
	rec["relative_to"], err = type_13250047047225666367_ToNu(v.RelativeTo)
	if err != nil {
		return nu.Value{}, err
	}
	rec["absolute"], err = type_15050730807189225719_ToNu(v.Absolute)
	if err != nil {
		return nu.Value{}, err
	}
	return nu.Value{Value: rec}, nil
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

var type_3391629412554670232 = types.String()

func type_3391629412554670232_FromNu(v nu.Value) (out props.EventClass, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("props.EventClass: %w", err)
		}
	}()
	casted, ok := v.Value.(string)
	converted := props.EventClass(casted)
	if !ok {
		return converted, fmt.Errorf("expected string got %v", v.Value)
	}
	return converted, nil
}
func type_3391629412554670232_ToNu(v props.EventClass) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("props.EventClass: %w", err)
		}
	}()
	return nu.ToValue(v), nil
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

var type_15563090433295620771 = types.Table(type_15848486708820041698)

func type_15563090433295620771_FromNu(v nu.Value) (out dto.PropValueList, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("dto.PropValueList: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	arr, ok := v.Value.([]nu.Value)
	if !ok {
		return nil, fmt.Errorf("expected []nu.Value got %T", v.Value)
	}
	out = make(dto.PropValueList, len(arr))
	for i, e := range arr {
		out[i], err = type_15848486708820041698_FromNu(e)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
func type_15563090433295620771_ToNu(v dto.PropValueList) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("dto.PropValueList: %w", err)
		}
	}()
	list := make([]nu.Value, len(v))
	for i, e := range v {
		list[i], err = type_15848486708820041698_ToNu(e)
		if err != nil {
			return nu.Value{}, err
		}
	}
	return nu.Value{Value: list}, nil
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

var type_11923325321682739420 = types.Table(type_1233005477764658533)

func type_11923325321682739420_FromNu(v nu.Value) (out dto.Timeline, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("dto.Timeline: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	arr, ok := v.Value.([]nu.Value)
	if !ok {
		return nil, fmt.Errorf("expected []nu.Value got %T", v.Value)
	}
	out = make(dto.Timeline, len(arr))
	for i, e := range arr {
		out[i], err = type_1233005477764658533_FromNu(e)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
func type_11923325321682739420_ToNu(v dto.Timeline) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("dto.Timeline: %w", err)
		}
	}()
	list := make([]nu.Value, len(v))
	for i, e := range v {
		list[i], err = type_1233005477764658533_ToNu(e)
		if err != nil {
			return nu.Value{}, err
		}
	}
	return nu.Value{Value: list}, nil
}

var type_16184234370936045305 = types.String()

func type_16184234370936045305_FromNu(v nu.Value) (out props.EventStatus, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("props.EventStatus: %w", err)
		}
	}()
	casted, ok := v.Value.(string)
	converted := props.EventStatus(casted)
	if !ok {
		return converted, fmt.Errorf("expected string got %v", v.Value)
	}
	return converted, nil
}
func type_16184234370936045305_ToNu(v props.EventStatus) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("props.EventStatus: %w", err)
		}
	}()
	return nu.ToValue(v), nil
}

var type_1233005477764658533 = types.RecordDef{
	"now":           type_8047992331715851194,
	"duration":      type_16589689216511618220,
	"active_events": type_11669970230249425419,
}

func type_1233005477764658533_FromNu(v nu.Value) (out dto.TimeSegment, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("dto.TimeSegment: %w", err)
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
	out.ActiveEvents, err = type_11669970230249425419_FromNu(val)
	if err != nil {
		return out, err
	}
	return out, nil
}
func type_1233005477764658533_ToNu(v dto.TimeSegment) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("dto.TimeSegment: %w", err)
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
	rec["active_events"], err = type_11669970230249425419_ToNu(v.ActiveEvents)
	if err != nil {
		return nu.Value{}, err
	}
	return nu.Value{Value: rec}, nil
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

var type_18311196003507852142 = types.RecordDef{
	"latitude":  type_17860233973098560385,
	"longitude": type_17860233973098560385,
}

func type_18311196003507852142_FromNu(v nu.Value) (out props.EventGeo, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("props.EventGeo: %w", err)
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
func type_18311196003507852142_ToNu(v props.EventGeo) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("props.EventGeo: %w", err)
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

var type_6387110243123677632 = types.String()

func type_6387110243123677632_FromNu(v nu.Value) (out props.EventTransparency, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("props.EventTransparency: %w", err)
		}
	}()
	casted, ok := v.Value.(string)
	converted := props.EventTransparency(casted)
	if !ok {
		return converted, fmt.Errorf("expected string got %v", v.Value)
	}
	return converted, nil
}
func type_6387110243123677632_ToNu(v props.EventTransparency) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("props.EventTransparency: %w", err)
		}
	}()
	return nu.ToValue(v), nil
}

var type_15848486708820041698 = types.RecordDef{
	"value":  type_15613163272824911089,
	"params": type_14293658896741725053,
}

func type_15848486708820041698_FromNu(v nu.Value) (out dto.PropValue, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("dto.PropValue: %w", err)
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
func type_15848486708820041698_ToNu(v dto.PropValue) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("dto.PropValue: %w", err)
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

var type_1838685811995560013 = types.Table(type_18369289839240265122)

func type_1838685811995560013_FromNu(v nu.Value) (out dto.CalendarList, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("dto.CalendarList: %w", err)
		}
	}()
	if v.Value == nil {
		return nil, nil
	}
	arr, ok := v.Value.([]nu.Value)
	if !ok {
		return nil, fmt.Errorf("expected []nu.Value got %T", v.Value)
	}
	out = make(dto.CalendarList, len(arr))
	for i, e := range arr {
		out[i], err = type_18369289839240265122_FromNu(e)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
func type_1838685811995560013_ToNu(v dto.CalendarList) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("dto.CalendarList: %w", err)
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

var type_7406295723486674371 = types.String()

func type_7406295723486674371_FromNu(v nu.Value) (out dto.RRule, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("dto.RRule: %w", err)
		}
	}()
	if v.Value == nil {
		return dto.RRule{}, nil
	}
	parsed, err := rrule.StrToRRule(v.Value.(string))
	if err != nil {
		return dto.RRule{}, err
	}
	return dto.RRule{RRule: parsed}, nil
}
func type_7406295723486674371_ToNu(v dto.RRule) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("dto.RRule: %w", err)
		}
	}()
	if v.RRule == nil {
		return nu.Value{Value: nil}, nil
	}
	return nu.ToValue(v.String()), nil
}

var type_18418790197363493434 = types.RecordDef{
	"stamp":    type_8047992331715851194,
	"all_day":  type_729807561129781588,
	"floating": type_729807561129781588,
}

func type_18418790197363493434_FromNu(v nu.Value) (out props.Datetime, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("props.Datetime: %w", err)
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
func type_18418790197363493434_ToNu(v props.Datetime) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("props.Datetime: %w", err)
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

var type_13250047047225666367 = types.Int()

func type_13250047047225666367_FromNu(v nu.Value) (out props.EventTriggerRelative, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("props.EventTriggerRelative: %w", err)
		}
	}()
	casted, ok := v.Value.(int64)
	converted := props.EventTriggerRelative(casted)
	if !ok {
		return converted, fmt.Errorf("expected int64 got %v", v.Value)
	}
	return converted, nil
}
func type_13250047047225666367_ToNu(v props.EventTriggerRelative) (out nu.Value, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("props.EventTriggerRelative: %w", err)
		}
	}()
	return nu.ToValue(v), nil
}

var EventStatusType = type_16184234370936045305
var EventStatusFromNu = type_16184234370936045305_FromNu
var EventStatusToNu = type_16184234370936045305_ToNu
var EventClassType = type_3391629412554670232
var EventClassFromNu = type_3391629412554670232_FromNu
var EventClassToNu = type_3391629412554670232_ToNu
var PropValueType = type_15563090433295620771
var PropValueFromNu = type_15563090433295620771_FromNu
var PropValueToNu = type_15563090433295620771_ToNu
var TimeSegmentType = type_1233005477764658533
var TimeSegmentFromNu = type_1233005477764658533_FromNu
var TimeSegmentToNu = type_1233005477764658533_ToNu
var TimelineType = type_11923325321682739420
var TimelineFromNu = type_11923325321682739420_FromNu
var TimelineToNu = type_11923325321682739420_ToNu
var DatetimeType = type_18418790197363493434
var DatetimeFromNu = type_18418790197363493434_FromNu
var DatetimeToNu = type_18418790197363493434_ToNu
var EventGeoType = type_18311196003507852142
var EventGeoFromNu = type_18311196003507852142_FromNu
var EventGeoToNu = type_18311196003507852142_ToNu
var EventTransparencyType = type_6387110243123677632
var EventTransparencyFromNu = type_6387110243123677632_FromNu
var EventTransparencyToNu = type_6387110243123677632_ToNu
var PropValueListType = type_15563090433295620771
var PropValueListFromNu = type_15563090433295620771_FromNu
var PropValueListToNu = type_15563090433295620771_ToNu
var CalendarListType = type_1838685811995560013
var CalendarListFromNu = type_1838685811995560013_FromNu
var CalendarListToNu = type_1838685811995560013_ToNu
var RRuleType = type_7406295723486674371
var RRuleFromNu = type_7406295723486674371_FromNu
var RRuleToNu = type_7406295723486674371_ToNu
var EventTriggerType = type_11339315157932765317
var EventTriggerFromNu = type_11339315157932765317_FromNu
var EventTriggerToNu = type_11339315157932765317_ToNu

func Marshal(v any) (out nu.Value, err error) {
	switch v := v.(type) {
	case map[string][]string:
		out, err = type_14293658896741725053_ToNu(v)
		return
	case bool:
		out, err = type_729807561129781588_ToNu(v)
		return
	case props.EventClass:
		out, err = type_3391629412554670232_ToNu(v)
		return
	case []string:
		out, err = type_11669970230249425419_ToNu(v)
		return
	case dto.PropValueList:
		out, err = type_15563090433295620771_ToNu(v)
		return
	case int64:
		out, err = type_15139881813094606131_ToNu(v)
		return
	case caldav.Calendar:
		out, err = type_18369289839240265122_ToNu(v)
		return
	case dto.Timeline:
		out, err = type_11923325321682739420_ToNu(v)
		return
	case props.EventStatus:
		out, err = type_16184234370936045305_ToNu(v)
		return
	case dto.TimeSegment:
		out, err = type_1233005477764658533_ToNu(v)
		return
	case time.Time:
		out, err = type_8047992331715851194_ToNu(v)
		return
	case time.Duration:
		out, err = type_16589689216511618220_ToNu(v)
		return
	case props.EventGeo:
		out, err = type_18311196003507852142_ToNu(v)
		return
	case props.EventTransparency:
		out, err = type_6387110243123677632_ToNu(v)
		return
	case dto.PropValue:
		out, err = type_15848486708820041698_ToNu(v)
		return
	case dto.CalendarList:
		out, err = type_1838685811995560013_ToNu(v)
		return
	case dto.RRule:
		out, err = type_7406295723486674371_ToNu(v)
		return
	case props.Datetime:
		out, err = type_18418790197363493434_ToNu(v)
		return
	case *time.Duration:
		out, err = type_5863190983406162214_ToNu(v)
		return
	case props.EventTriggerRelative:
		out, err = type_13250047047225666367_ToNu(v)
		return
	case *time.Time:
		out, err = type_15050730807189225719_ToNu(v)
		return
	case props.EventTrigger:
		out, err = type_11339315157932765317_ToNu(v)
		return
	case float64:
		out, err = type_17860233973098560385_ToNu(v)
		return
	case string:
		out, err = type_15613163272824911089_ToNu(v)
		return
	}
	err = fmt.Errorf(`unsupported type: %T`, v)
	return
}
func Unmarshal(v nu.Value, out any) (err error) {
	switch out := out.(type) {
	case *map[string][]string:
		*out, err = type_14293658896741725053_FromNu(v)
		return
	case *bool:
		*out, err = type_729807561129781588_FromNu(v)
		return
	case *props.EventClass:
		*out, err = type_3391629412554670232_FromNu(v)
		return
	case *[]string:
		*out, err = type_11669970230249425419_FromNu(v)
		return
	case *dto.PropValueList:
		*out, err = type_15563090433295620771_FromNu(v)
		return
	case *int64:
		*out, err = type_15139881813094606131_FromNu(v)
		return
	case *caldav.Calendar:
		*out, err = type_18369289839240265122_FromNu(v)
		return
	case *dto.Timeline:
		*out, err = type_11923325321682739420_FromNu(v)
		return
	case *props.EventStatus:
		*out, err = type_16184234370936045305_FromNu(v)
		return
	case *dto.TimeSegment:
		*out, err = type_1233005477764658533_FromNu(v)
		return
	case *time.Time:
		*out, err = type_8047992331715851194_FromNu(v)
		return
	case *time.Duration:
		*out, err = type_16589689216511618220_FromNu(v)
		return
	case *props.EventGeo:
		*out, err = type_18311196003507852142_FromNu(v)
		return
	case *props.EventTransparency:
		*out, err = type_6387110243123677632_FromNu(v)
		return
	case *dto.PropValue:
		*out, err = type_15848486708820041698_FromNu(v)
		return
	case *dto.CalendarList:
		*out, err = type_1838685811995560013_FromNu(v)
		return
	case *dto.RRule:
		*out, err = type_7406295723486674371_FromNu(v)
		return
	case *props.Datetime:
		*out, err = type_18418790197363493434_FromNu(v)
		return
	case **time.Duration:
		*out, err = type_5863190983406162214_FromNu(v)
		return
	case *props.EventTriggerRelative:
		*out, err = type_13250047047225666367_FromNu(v)
		return
	case **time.Time:
		*out, err = type_15050730807189225719_FromNu(v)
		return
	case *props.EventTrigger:
		*out, err = type_11339315157932765317_FromNu(v)
		return
	case *float64:
		*out, err = type_17860233973098560385_FromNu(v)
		return
	case *string:
		*out, err = type_15613163272824911089_FromNu(v)
		return
	}
	err = fmt.Errorf(`unsupported type: %T`, out)
	return
}
