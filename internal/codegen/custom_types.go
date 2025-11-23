package main

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/teambition/rrule-go"
)

var timeType = reflect.TypeOf(time.Time{}).String()
var durType = reflect.TypeOf(time.Duration(0)).String()
var urlType = reflect.TypeOf(&url.URL{}).String()
var rruleType = reflect.TypeOf(&rrule.RRule{}).String()

func init() {
	// time.Time support
	customTypes = append(customTypes, func(t reflect.Type) (typefn typeDeclFn, fromfn fromDeclFn, tofn toDeclFn) {
		if t.String() != timeType {
			return
		}
		typefn = func(cache map[uint64]TypeDecl, t reflect.Type) (out TypeDecl) {
			out.TypeId = typeId(t)
			if existing, ok := cache[out.TypeId]; ok {
				return existing
			}
			defer func() { cache[out.TypeId] = out }()

			out.Value = "types.Date()"
			return
		}
		fromfn = func(cache map[uint64]FromDecl, t reflect.Type) (out FromDecl) {
			out.TypeId = typeId(t)
			if existing, ok := cache[out.TypeId]; ok {
				return existing
			}
			defer func() { cache[out.TypeId] = out }()

			(&out).SetTypeStr("time.Time")

			var sb strings.Builder
			fmt.Fprintln(&sb, "out, ok := v.Value.(time.Time)")
			fmt.Fprintf(&sb, "if !ok { return out, fmt.Errorf(\"expected time.Time got %%T\", v.Value) }\n")
			fmt.Fprintln(&sb, "return")
			out.Body = sb.String()
			return
		}
		tofn = func(cache map[uint64]ToDecl, t reflect.Type) (out ToDecl) {
			out.TypeId = typeId(t)
			if existing, ok := cache[out.TypeId]; ok {
				return existing
			}
			defer func() { cache[out.TypeId] = out }()

			(&out).SetTypeStr("time.Time")
			out.Body = "return nu.ToValue(v), nil"
			return
		}
		return
	})

	// time.Duration support
	customTypes = append(customTypes, func(t reflect.Type) (typefn typeDeclFn, fromfn fromDeclFn, tofn toDeclFn) {
		if t.String() != durType {
			return
		}
		typefn = func(cache map[uint64]TypeDecl, t reflect.Type) (out TypeDecl) {
			out.TypeId = typeId(t)
			if existing, ok := cache[out.TypeId]; ok {
				return existing
			}
			defer func() { cache[out.TypeId] = out }()

			out.Value = "types.Duration()"
			return
		}
		fromfn = func(cache map[uint64]FromDecl, t reflect.Type) (out FromDecl) {
			out.TypeId = typeId(t)
			if existing, ok := cache[out.TypeId]; ok {
				return existing
			}
			defer func() { cache[out.TypeId] = out }()

			(&out).SetTypeStr("time.Duration")

			var sb strings.Builder
			fmt.Fprintln(&sb, "out, ok := v.Value.(time.Duration)")
			fmt.Fprintf(&sb, "if !ok { return out, fmt.Errorf(\"expected time.Duration got %%T\", v.Value) }\n")
			fmt.Fprintln(&sb, "return")
			out.Body = sb.String()
			return
		}
		tofn = func(cache map[uint64]ToDecl, t reflect.Type) (out ToDecl) {
			out.TypeId = typeId(t)
			if existing, ok := cache[out.TypeId]; ok {
				return existing
			}
			defer func() { cache[out.TypeId] = out }()

			(&out).SetTypeStr("time.Duration")
			out.Body = "return nu.ToValue(v), nil"
			return
		}
		return
	})

	// *url.URL support
	customTypes = append(customTypes, func(t reflect.Type) (typefn typeDeclFn, fromfn fromDeclFn, tofn toDeclFn) {
		if t.String() != urlType {
			return
		}
		typefn = func(cache map[uint64]TypeDecl, t reflect.Type) (out TypeDecl) {
			out.TypeId = typeId(t)
			if existing, ok := cache[out.TypeId]; ok {
				return existing
			}
			defer func() { cache[out.TypeId] = out }()

			out.Value = "types.String()"
			return
		}
		fromfn = func(cache map[uint64]FromDecl, t reflect.Type) (out FromDecl) {
			out.TypeId = typeId(t)
			if existing, ok := cache[out.TypeId]; ok {
				return existing
			}
			defer func() { cache[out.TypeId] = out }()

			(&out).SetTypeStr("*url.URL")
			var sb strings.Builder
			fmt.Fprintln(&sb, "if v.Value == nil { return nil, nil }")
			fmt.Fprintln(&sb, "parsed, err := url.Parse(v.Value.(string))")
			fmt.Fprintln(&sb, "if err != nil { return nil, err }")
			fmt.Fprintln(&sb, "return parsed, nil")
			out.Body = sb.String()
			return
		}
		tofn = func(cache map[uint64]ToDecl, t reflect.Type) (out ToDecl) {
			out.TypeId = typeId(t)
			if existing, ok := cache[out.TypeId]; ok {
				return existing
			}
			defer func() { cache[out.TypeId] = out }()

			(&out).SetTypeStr("*url.URL")
			var sb strings.Builder
			fmt.Fprintln(&sb, "if v == nil { return nu.Value{Value: nil}, nil }")
			fmt.Fprintln(&sb, "return nu.ToValue(v.String()), nil")
			out.Body = sb.String()
			return
		}
		return
	})

	// *rrule.RRule support
	customTypes = append(customTypes, func(t reflect.Type) (typefn typeDeclFn, fromfn fromDeclFn, tofn toDeclFn) {
		if t.String() != rruleType {
			return
		}
		typefn = func(cache map[uint64]TypeDecl, t reflect.Type) (out TypeDecl) {
			out.TypeId = typeId(t)
			if existing, ok := cache[out.TypeId]; ok {
				return existing
			}
			defer func() { cache[out.TypeId] = out }()

			out.Value = "types.String()"
			return
		}
		fromfn = func(cache map[uint64]FromDecl, t reflect.Type) (out FromDecl) {
			out.TypeId = typeId(t)
			if existing, ok := cache[out.TypeId]; ok {
				return existing
			}
			defer func() { cache[out.TypeId] = out }()

			(&out).SetTypeStr("*rrule.RRule")
			var sb strings.Builder
			fmt.Fprintln(&sb, "if v.Value == nil { return nil, nil }")
			fmt.Fprintln(&sb, "parsed, err := rrule.StrToRRule(v.Value.(string))")
			fmt.Fprintln(&sb, "if err != nil { return nil, err }")
			fmt.Fprintln(&sb, "return parsed, nil")
			out.Body = sb.String()
			return
		}
		tofn = func(cache map[uint64]ToDecl, t reflect.Type) (out ToDecl) {
			out.TypeId = typeId(t)
			if existing, ok := cache[out.TypeId]; ok {
				return existing
			}
			defer func() { cache[out.TypeId] = out }()

			(&out).SetTypeStr("*rrule.RRule")
			var sb strings.Builder
			fmt.Fprintln(&sb, "if v == nil { return nu.Value{Value: nil}, nil }")
			fmt.Fprintln(&sb, "return nu.ToValue(v.String()), nil")
			out.Body = sb.String()
			return
		}
		return
	})
}
