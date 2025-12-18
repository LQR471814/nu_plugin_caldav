package main

import (
	"net/url"
	"reflect"
	"time"

	"github.com/LQR471814/nu_plugin_caldav/internal/dto"
)

var timeType = reflect.TypeOf(time.Time{}).String()
var durType = reflect.TypeOf(time.Duration(0)).String()
var urlType = reflect.TypeOf(&url.URL{}).String()
var rruleType = reflect.TypeOf(dto.RRule{}).String()

// time.Time support
type timestampBridge struct {
	t reflect.Type
}

func timestampRoute(router *BridgeTypeRouter, t reflect.Type) GoNuBridgeType {
	if t.String() != timeType {
		return nil
	}
	return timestampBridge{t: t}
}

func (t timestampBridge) GoType() reflect.Type {
	return t.t
}

func (t timestampBridge) TypeExpr() string {
	return "types.Date()"
}

func (t timestampBridge) FromBody() string {
	return `out, ok := v.Value.(time.Time)
if !ok { return out, fmt.Errorf("expected time.Time got %T", v.Value) }
return`
}

func (t timestampBridge) ToBody() string {
	return "return nu.ToValue(v), nil"
}

// time.Duration support
type durationBridge struct {
	t reflect.Type
}

func durationRoute(router *BridgeTypeRouter, t reflect.Type) GoNuBridgeType {
	if t.String() != durType {
		return nil
	}
	return durationBridge{t: t}
}

func (t durationBridge) GoType() reflect.Type {
	return t.t
}

func (t durationBridge) TypeExpr() string {
	return "types.Duration()"
}

func (t durationBridge) FromBody() string {
	return `out, ok := v.Value.(time.Duration)
if !ok { return out, fmt.Errorf("expected time.Duration got %T", v.Value) }
return`
}

func (t durationBridge) ToBody() string {
	return "return nu.ToValue(v), nil"
}

// *url.URL support
type urlBridge struct {
	t reflect.Type
}

func urlRoute(router *BridgeTypeRouter, t reflect.Type) GoNuBridgeType {
	if t.String() != urlType {
		return nil
	}
	return urlBridge{t: t}
}

func (t urlBridge) GoType() reflect.Type {
	return t.t
}

func (t urlBridge) TypeExpr() string {
	return "types.String()"
}

func (t urlBridge) FromBody() string {
	return `if v.Value == nil { return nil, nil }
parsed, err := url.Parse(v.Value.(string))
if err != nil { return nil, err }
return parsed, nil`
}

func (t urlBridge) ToBody() string {
	return `if v == nil { return nu.Value{Value: nil}, nil }
return nu.ToValue(v.String()), nil`
}

// *rrule.RRule support
type rruleBridge struct {
	t reflect.Type
}

func rruleRoute(router *BridgeTypeRouter, t reflect.Type) GoNuBridgeType {
	if t.String() != rruleType {
		return nil
	}
	return rruleBridge{t: t}
}

func (t rruleBridge) GoType() reflect.Type {
	return t.t
}

func (t rruleBridge) TypeExpr() string {
	return "types.String()"
}

func (t rruleBridge) FromBody() string {
	return `if v.Value == nil { return dto.RRule{}, nil }
parsed, err := rrule.StrToRRule(v.Value.(string))
if err != nil { return dto.RRule{}, err }
return dto.RRule{RRule: parsed}, nil`
}

func (t rruleBridge) ToBody() string {
	return `if v.RRule == nil { return nu.Value{Value: nil}, nil }
return nu.ToValue(v.String()), nil`
}
