package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/zeebo/xxh3"
)

// built-in primitives

var errType = reflect.TypeOf((*error)(nil)).Elem()

func builtinTypeDecl(cache map[uint64]TypeDecl, t reflect.Type) (out TypeDecl) {
	out.TypeId = typeId(t)
	if existing, ok := cache[out.TypeId]; ok {
		return existing
	}
	defer func() { cache[out.TypeId] = out }()

	if t.Implements(errType) {
		out.Value = "types.Error()"
		return
	}

	switch t.String() {
	case timeType:
		out.Value = "types.Date()"
		return
	case durType:
		out.Value = "types.Duration()"
		return
	}

	switch t.Kind() {
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			out.Value = "types.Binary()"
			return
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out.Value = "types.Int()"
		return
	case reflect.Float32, reflect.Float64:
		out.Value = "types.Float()"
		return
	case reflect.String:
		out.Value = "types.String()"
		return
	case reflect.Bool:
		out.Value = "types.Bool()"
		return
	}

	panic(fmt.Errorf("unsupported type: %v", t.String()))
}

func builtinFromDecl(cache map[uint64]FromDecl, t reflect.Type) (out FromDecl) {
	out.TypeId = typeId(t)
	if existing, ok := cache[out.TypeId]; ok {
		return existing
	}
	defer func() { cache[out.TypeId] = out }()

	(&out).SetTypeStr(t.String())

	underlyingType := ""

	switch t.Kind() {
	case reflect.String:
		underlyingType = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		underlyingType = "int64"
	case reflect.Float32, reflect.Float64:
		underlyingType = "float64"
	case reflect.Bool:
		underlyingType = "bool"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "casted, ok := v.Value.(%s)\n", underlyingType)
	fmt.Fprintf(&sb, "converted := %s(casted)\n", t.String())
	fmt.Fprintf(&sb, "if !ok { return converted, fmt.Errorf(\"expected %s got %%v\", v.Value) }\n", underlyingType)
	fmt.Fprintln(&sb, "return converted, nil")
	out.Body = sb.String()

	return
}

func builtinToDecl(cache map[uint64]ToDecl, t reflect.Type) (out ToDecl) {
	out.TypeId = typeId(t)
	if existing, ok := cache[out.TypeId]; ok {
		return existing
	}
	defer func() { cache[out.TypeId] = out }()

	(&out).SetTypeStr(t.String())

	out.Body = "return nu.ToValue(v), nil"
	return
}

// pointer types

func pointerTypeDecl(cache map[uint64]TypeDecl, t reflect.Type) (out TypeDecl) {
	return typeDecl(cache, t.Elem())
}

func pointerFromDecl(cache map[uint64]FromDecl, t reflect.Type) (out FromDecl) {
	out.TypeId = typeId(t)
	if existing, ok := cache[out.TypeId]; ok {
		return existing
	}
	defer func() { cache[out.TypeId] = out }()
	(&out).SetTypeStr(t.String())

	var sb strings.Builder
	fmt.Fprintln(&sb, "if v.Value == nil { return nil, nil }")
	decl := fromDecl(cache, t.Elem())
	fmt.Fprintf(&sb, "res, err := %s(v)\n", FromDeclId(decl.TypeId))
	fmt.Fprintln(&sb, "if err != nil { return nil, err }")
	fmt.Fprintf(&sb, "return &res, nil")
	out.Body = sb.String()

	return
}

func pointerToDecl(cache map[uint64]ToDecl, t reflect.Type) (out ToDecl) {
	out.TypeId = typeId(t)
	if existing, ok := cache[out.TypeId]; ok {
		return existing
	}
	defer func() { cache[out.TypeId] = out }()

	(&out).SetTypeStr(t.String())

	var sb strings.Builder
	fmt.Fprintln(&sb, "if v == nil { return nu.Value{}, nil }")
	decl := toDecl(cache, t.Elem())
	fmt.Fprintf(&sb, "return %s(*v)", ToDeclId(decl.TypeId))
	out.Body = sb.String()

	return
}

// slice types

func sliceTypeDecl(cache map[uint64]TypeDecl, t reflect.Type) (out TypeDecl) {
	out.TypeId = typeId(t)
	if existing, ok := cache[out.TypeId]; ok {
		return existing
	}
	defer func() { cache[out.TypeId] = out }()

	decl := typeDecl(cache, t.Elem())

	if t.Elem().Kind() == reflect.Struct {
		recordDefId := xxh3.Hash([]byte(t.Elem().String()))
		return TypeDecl{
			TypeId: typeId(t),
			Value:  fmt.Sprintf("types.Table(%s)", TypeDeclId(recordDefId)),
		}
	}

	return TypeDecl{
		TypeId: typeId(t),
		Value:  fmt.Sprintf("types.List(%s)", TypeDeclId(decl.TypeId)),
	}
}

func sliceFromDecl(cache map[uint64]FromDecl, t reflect.Type) (out FromDecl) {
	out.TypeId = typeId(t)
	if existing, ok := cache[out.TypeId]; ok {
		return existing
	}
	defer func() { cache[out.TypeId] = out }()

	(&out).SetTypeStr(t.String())

	decl := fromDecl(cache, t.Elem())

	var sb strings.Builder
	fmt.Fprintln(&sb, "if v.Value == nil { return nil, nil }")
	fmt.Fprintln(&sb, "arr, ok := v.Value.([]nu.Value)")
	fmt.Fprintf(&sb, "if !ok { return nil, fmt.Errorf(\"expected []nu.Value got %%T\", v.Value) }\n")
	fmt.Fprintf(&sb, "out = make(%s, len(arr))\n", out.TypeStr)
	fmt.Fprintln(&sb, "for i, e := range arr {")
	fmt.Fprintf(&sb, "out[i], err = %s(e)\n", FromDeclId(decl.TypeId))
	fmt.Fprintln(&sb, "if err != nil { return nil, err }")
	fmt.Fprintln(&sb, "}")
	fmt.Fprintln(&sb, "return out, nil")

	out.Body = sb.String()
	return
}

func sliceToDecl(cache map[uint64]ToDecl, t reflect.Type) (out ToDecl) {
	out.TypeId = typeId(t)
	if existing, ok := cache[out.TypeId]; ok {
		return existing
	}
	defer func() { cache[out.TypeId] = out }()

	(&out).SetTypeStr(t.String())

	decl := toDecl(cache, t.Elem())

	var sb strings.Builder
	fmt.Fprintln(&sb, "list := make([]nu.Value, len(v))")
	fmt.Fprintln(&sb, "for i, e := range v {")
	fmt.Fprintf(&sb, "list[i], err = %s(e)\n", ToDeclId(decl.TypeId))
	fmt.Fprintln(&sb, "if err != nil { return nu.Value{}, err }")
	fmt.Fprintln(&sb, "}")
	fmt.Fprintln(&sb, "return nu.Value{Value: list}, nil")

	out.Body = sb.String()
	return
}

// map types

func mapTypeDecl(cache map[uint64]TypeDecl, t reflect.Type) (out TypeDecl) {
	out.TypeId = typeId(t)
	if existing, ok := cache[out.TypeId]; ok {
		return existing
	}
	defer func() { cache[out.TypeId] = out }()

	return TypeDecl{
		TypeId: typeId(t),
		// TODO: looks like nushell's type system isn't expressive enough yet
		Value: "types.Any()",
	}
}

func mapFromDecl(cache map[uint64]FromDecl, t reflect.Type) (out FromDecl) {
	out.TypeId = typeId(t)
	if existing, ok := cache[out.TypeId]; ok {
		return existing
	}
	defer func() { cache[out.TypeId] = out }()

	(&out).SetTypeStr(t.String())

	valDecl := fromDecl(cache, t.Elem())

	var sb strings.Builder
	fmt.Fprintln(&sb, "dict, ok := v.Value.(nu.Record)")
	fmt.Fprintf(&sb, "if !ok { return nil, fmt.Errorf(\"expected nu.Record got %%T\", v.Value) }\n")
	fmt.Fprintf(&sb, "out = make(%s, len(dict))\n", out.TypeStr)
	fmt.Fprintln(&sb, "for k, v := range dict {")
	fmt.Fprintf(&sb, "out[k], err = %s(v)\n", FromDeclId(valDecl.TypeId))
	fmt.Fprintln(&sb, "if err != nil { return nil, err }")
	fmt.Fprintln(&sb, "}")
	fmt.Fprintln(&sb, "return out, nil")

	out.Body = sb.String()
	return
}

func mapToDecl(cache map[uint64]ToDecl, t reflect.Type) (out ToDecl) {
	out.TypeId = typeId(t)
	if existing, ok := cache[out.TypeId]; ok {
		return existing
	}
	defer func() { cache[out.TypeId] = out }()

	(&out).SetTypeStr(t.String())

	valDecl := toDecl(cache, t.Elem())

	var sb strings.Builder
	fmt.Fprintln(&sb, "dict := make(nu.Record, len(v))")
	fmt.Fprintln(&sb, "for k, v := range v {")
	fmt.Fprintf(&sb, "dict[k], err = %s(v)\n", ToDeclId(valDecl.TypeId))
	fmt.Fprintln(&sb, "if err != nil { return nu.Value{}, err }")
	fmt.Fprintln(&sb, "}")
	fmt.Fprintln(&sb, "return nu.Value{Value: dict}, nil")

	out.Body = sb.String()
	return
}

// struct types

func nuFieldName(t reflect.Type, f reflect.StructField) string {
	fieldName, ok := f.Tag.Lookup("name")
	if !ok {
		fieldName = pascalToSnakeCase(f.Name)
	}
	if fieldName == "" {
		panic(fmt.Sprintf(
			"cannot have empty fieldname for field '%s.%s'",
			t.Name(),
			f.Name,
		))
	}
	return fieldName
}

func structRecordDef(cache map[uint64]TypeDecl, t reflect.Type) (out TypeDecl) {
	if t.Kind() != reflect.Struct {
		panic("cannot use non-struct as struct")
	}

	out.TypeId = typeId(t)
	if existing, ok := cache[out.TypeId]; ok {
		return existing
	}
	defer func() { cache[out.TypeId] = out }()

	var sb strings.Builder
	fmt.Fprintln(&sb, "types.RecordDef{")
	for i := range t.NumField() {
		f := t.Field(i)

		recordField := nuFieldName(t, f)
		if recordField == "-" {
			continue
		}

		decl := typeDecl(cache, f.Type)
		fmt.Fprintf(&sb, "\"%s\": %s,\n", recordField, TypeDeclId(decl.TypeId))
	}
	fmt.Fprint(&sb, "}")

	out.Value = sb.String()
	return
}

func structTypeDecl(cache map[uint64]TypeDecl, t reflect.Type) (out TypeDecl) {
	if t.Kind() != reflect.Struct {
		panic("cannot use non-struct as struct")
	}

	id := xxh3.Hash([]byte(t.String() + " wrapped"))
	out.TypeId = id
	if existing, ok := cache[out.TypeId]; ok {
		return existing
	}
	defer func() { cache[out.TypeId] = out }()

	defDecl := structRecordDef(cache, t)

	var sb strings.Builder
	fmt.Fprintf(&sb, "types.Record(%s)", TypeDeclId(defDecl.TypeId))
	out.Value = sb.String()

	return
}

func structFromDecl(cache map[uint64]FromDecl, t reflect.Type) (out FromDecl) {
	if t.Kind() != reflect.Struct {
		panic("cannot use non-struct as struct")
	}

	out.TypeId = typeId(t)
	if existing, ok := cache[out.TypeId]; ok {
		return existing
	}
	defer func() { cache[out.TypeId] = out }()

	(&out).SetTypeStr(t.String())

	var sb strings.Builder

	fmt.Fprintln(&sb, "record, ok := v.Value.(nu.Record)")
	fmt.Fprintf(&sb, "if !ok { return out, fmt.Errorf(\"expected nu.Record got %%T\", v.Value) }\n")
	fmt.Fprintln(&sb, "var val nu.Value")

	for i := range t.NumField() {
		f := t.Field(i)

		recordField := nuFieldName(t, f)
		if recordField == "-" {
			continue
		}
		defaultVal := f.Tag.Get("default")

		decl := fromDecl(cache, f.Type)

		if defaultVal != "" {
			fmt.Fprint(&sb, "val, ok")
		} else {
			fmt.Fprint(&sb, "val, _")
		}
		fmt.Fprintf(&sb, " = record[\"%s\"]\n", recordField)

		if defaultVal != "" {
			fmt.Fprintf(&sb, "if !ok { out.%s = %s } ", f.Name, defaultVal)
			fmt.Fprintln(&sb, "else {")
		}
		fmt.Fprintf(
			&sb,
			// all nullables are represented with a pointer, therefore if the
			// key doesn't exist it will simply return nil, which is the
			// correct behavior
			"out.%s, err = %s(val)\n",
			f.Name,
			FromDeclId(decl.TypeId),
		)
		fmt.Fprintln(&sb, "if err != nil { return out, err }")
		if defaultVal != "" {
			fmt.Fprintln(&sb, "}")
		}
	}
	fmt.Fprintln(&sb, "return out, nil")
	out.Body = sb.String()

	return
}

func structToDecl(cache map[uint64]ToDecl, t reflect.Type) (out ToDecl) {
	if t.Kind() != reflect.Struct {
		panic("cannot use non-struct as struct")
	}

	out.TypeId = typeId(t)
	if existing, ok := cache[out.TypeId]; ok {
		return existing
	}
	defer func() { cache[out.TypeId] = out }()

	(&out).SetTypeStr(t.String())

	var sb strings.Builder

	fmt.Fprintln(&sb, "rec := nu.Record{}")
	for i := range t.NumField() {
		f := t.Field(i)

		recordField := nuFieldName(t, f)
		if recordField == "-" {
			continue
		}

		decl := toDecl(cache, f.Type)
		fmt.Fprintf(&sb, "rec[\"%s\"], err = %s(v.%s)\n", recordField, ToDeclId(decl.TypeId), f.Name)
		fmt.Fprintln(&sb, "if err != nil { return nu.Value{}, err }")
	}
	fmt.Fprintln(&sb, "return nu.Value{Value: rec}, nil")
	out.Body = sb.String()

	return
}

// router functions

type typeDeclFn = func(cache map[uint64]TypeDecl, t reflect.Type) (out TypeDecl)
type fromDeclFn = func(cache map[uint64]FromDecl, t reflect.Type) (out FromDecl)
type toDeclFn = func(cache map[uint64]ToDecl, t reflect.Type) (out ToDecl)

var customTypes []func(t reflect.Type) (typefn typeDeclFn, fromfn fromDeclFn, tofn toDeclFn)

func typeDecl(cache map[uint64]TypeDecl, t reflect.Type) TypeDecl {
	for _, custom := range customTypes {
		fn, _, _ := custom(t)
		if fn != nil {
			return fn(cache, t)
		}
	}

	switch t.Kind() {
	case reflect.Struct:
		return structTypeDecl(cache, t)
	case reflect.Map:
		return mapTypeDecl(cache, t)
	case reflect.Slice:
		return sliceTypeDecl(cache, t)
	case reflect.Pointer:
		return pointerTypeDecl(cache, t)
	default:
		return builtinTypeDecl(cache, t)
	}
}

func fromDecl(cache map[uint64]FromDecl, t reflect.Type) FromDecl {
	for _, custom := range customTypes {
		_, fn, _ := custom(t)
		if fn != nil {
			return fn(cache, t)
		}
	}

	switch t.Kind() {
	case reflect.Struct:
		return structFromDecl(cache, t)
	case reflect.Map:
		return mapFromDecl(cache, t)
	case reflect.Slice:
		return sliceFromDecl(cache, t)
	case reflect.Pointer:
		return pointerFromDecl(cache, t)
	default:
		return builtinFromDecl(cache, t)
	}
}

func toDecl(cache map[uint64]ToDecl, t reflect.Type) ToDecl {
	for _, custom := range customTypes {
		_, _, fn := custom(t)
		if fn != nil {
			return fn(cache, t)
		}
	}

	switch t.Kind() {
	case reflect.Struct:
		return structToDecl(cache, t)
	case reflect.Map:
		return mapToDecl(cache, t)
	case reflect.Slice:
		return sliceToDecl(cache, t)
	case reflect.Pointer:
		return pointerToDecl(cache, t)
	default:
		return builtinToDecl(cache, t)
	}
}
