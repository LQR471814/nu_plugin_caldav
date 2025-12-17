package main

import (
	"fmt"
	"reflect"
	"strings"
)

type primitiveBridge struct {
	t reflect.Type
}

func primitiveRoute(router *BridgeTypeRouter, t reflect.Type) GoNuBridgeType {
	return primitiveBridge{t}
}

func (p primitiveBridge) GoType() reflect.Type {
	return p.t
}

var err_type = reflect.TypeOf((*error)(nil)).Elem()

func (p primitiveBridge) TypeExpr() (out string) {
	if p.t.Implements(err_type) {
		out = "types.Error()"
		return
	}

	switch p.t.String() {
	case timeType:
		out = "types.Date()"
		return
	case durType:
		out = "types.Duration()"
		return
	}

	switch p.t.Kind() {
	case reflect.Slice:
		if p.t.Elem().Kind() == reflect.Uint8 {
			out = "types.Binary()"
			return
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out = "types.Int()"
		return
	case reflect.Float32, reflect.Float64:
		out = "types.Float()"
		return
	case reflect.String:
		out = "types.String()"
		return
	case reflect.Bool:
		out = "types.Bool()"
		return
	}

	panic(fmt.Errorf("unsupported type: %v", p.t.String()))
}

func (p primitiveBridge) FromBody() (out string) {
	underlyingType := ""
	switch p.t.Kind() {
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
	fmt.Fprintf(&sb, "converted := %s(casted)\n", p.t.String())
	fmt.Fprintf(&sb, "if !ok { return converted, fmt.Errorf(\"expected %s got %%v\", v.Value) }\n", underlyingType)
	fmt.Fprintln(&sb, "return converted, nil")
	out = sb.String()
	return
}

func (p primitiveBridge) ToBody() (out string) {
	out = "return nu.ToValue(v), nil"
	return
}

// pointer types

type pointerBridge struct {
	t      reflect.Type
	router *BridgeTypeRouter
}

func pointerRoute(router *BridgeTypeRouter, t reflect.Type) GoNuBridgeType {
	if t.Kind() != reflect.Pointer {
		return nil
	}
	return pointerBridge{t: t, router: router}
}

func (p pointerBridge) GoType() reflect.Type {
	// TODO: oneof types don't work in plugin type definitions
	return p.t
}

func (p pointerBridge) TypeExpr() string {
	child := p.router.Lookup(p.t.Elem())
	if _, ok := child.(structBridge); ok {
		return fmt.Sprintf("types.Record(%s)", TypeDeclSyntaxID(child.GoType()))
	}
	return TypeDeclSyntaxID(child.GoType())
}

func (p pointerBridge) FromBody() (out string) {
	var sb strings.Builder
	fmt.Fprintln(&sb, "if v.Value == nil { return nil, nil }")
	child := p.router.Lookup(p.t.Elem()).GoType()
	fmt.Fprintf(&sb, "res, err := %s(v)\n", FromDeclSyntaxID(child))
	fmt.Fprintln(&sb, "if err != nil { return nil, err }")
	fmt.Fprintf(&sb, "return &res, nil")
	out = sb.String()
	return
}

func (p pointerBridge) ToBody() (out string) {
	var sb strings.Builder
	fmt.Fprintln(&sb, "if v == nil { return nu.Value{}, nil }")
	child := p.router.Lookup(p.t.Elem()).GoType()
	fmt.Fprintf(&sb, "return %s(*v)", ToDeclSyntaxID(child))
	out = sb.String()
	return
}

// slice types

type sliceBridge struct {
	t      reflect.Type
	router *BridgeTypeRouter
}

func sliceRoute(router *BridgeTypeRouter, t reflect.Type) GoNuBridgeType {
	if t.Kind() != reflect.Slice {
		return nil
	}
	return sliceBridge{t: t, router: router}
}

func (s sliceBridge) GoType() reflect.Type {
	return s.t
}

func (s sliceBridge) TypeExpr() string {
	child := s.router.Lookup(s.t.Elem()).GoType()
	if s.t.Elem().Kind() == reflect.Struct {
		return fmt.Sprintf("types.Table(%s)", TypeDeclSyntaxID(child))
	}
	return fmt.Sprintf("types.List(%s)", TypeDeclSyntaxID(child))
}

func (s sliceBridge) FromBody() (out string) {
	child := s.router.Lookup(s.t.Elem()).GoType()

	var sb strings.Builder
	fmt.Fprintln(&sb, "if v.Value == nil { return nil, nil }")
	fmt.Fprintln(&sb, "arr, ok := v.Value.([]nu.Value)")
	fmt.Fprintf(&sb, "if !ok { return nil, fmt.Errorf(\"expected []nu.Value got %%T\", v.Value) }\n")
	fmt.Fprintf(&sb, "out = make(%s, len(arr))\n", s.t.String())
	fmt.Fprintln(&sb, "for i, e := range arr {")
	fmt.Fprintf(&sb, "out[i], err = %s(e)\n", FromDeclSyntaxID(child))
	fmt.Fprintln(&sb, "if err != nil { return nil, err }")
	fmt.Fprintln(&sb, "}")
	fmt.Fprintln(&sb, "return out, nil")

	return sb.String()
}

func (s sliceBridge) ToBody() (out string) {
	child := s.router.Lookup(s.t.Elem()).GoType()

	var sb strings.Builder
	fmt.Fprintln(&sb, "list := make([]nu.Value, len(v))")
	fmt.Fprintln(&sb, "for i, e := range v {")
	fmt.Fprintf(&sb, "list[i], err = %s(e)\n", ToDeclSyntaxID(child))
	fmt.Fprintln(&sb, "if err != nil { return nu.Value{}, err }")
	fmt.Fprintln(&sb, "}")
	fmt.Fprintln(&sb, "return nu.Value{Value: list}, nil")

	out = sb.String()
	return
}

// map types

type mapBridge struct {
	t      reflect.Type
	router *BridgeTypeRouter
}

func mapRoute(router *BridgeTypeRouter, t reflect.Type) GoNuBridgeType {
	if t.Kind() != reflect.Map {
		return nil
	}
	return mapBridge{t: t, router: router}
}

func (m mapBridge) GoType() reflect.Type {
	return m.t
}

func (m mapBridge) TypeExpr() string {
	// TODO: looks like nushell's type system isn't expressive enough yet
	return "types.Any()"
}

func (m mapBridge) FromBody() (out string) {
	val := m.router.Lookup(m.t.Elem()).GoType()

	var sb strings.Builder
	fmt.Fprintln(&sb, "if v.Value == nil { return nil, nil }")
	fmt.Fprintln(&sb, "dict, ok := v.Value.(nu.Record)")
	fmt.Fprintf(&sb, "if !ok { return nil, fmt.Errorf(\"expected nu.Record got %%T\", v.Value) }\n")
	fmt.Fprintf(&sb, "out = make(%s, len(dict))\n", m.t.String())
	fmt.Fprintln(&sb, "for k, v := range dict {")
	fmt.Fprintf(&sb, "out[k], err = %s(v)\n", FromDeclSyntaxID(val))
	fmt.Fprintln(&sb, "if err != nil { return nil, err }")
	fmt.Fprintln(&sb, "}")
	fmt.Fprintln(&sb, "return out, nil")

	out = sb.String()
	return
}

func (m mapBridge) ToBody() (out string) {
	val := m.router.Lookup(m.t.Elem()).GoType()

	var sb strings.Builder
	fmt.Fprintln(&sb, "dict := make(nu.Record, len(v))")
	fmt.Fprintln(&sb, "for k, v := range v {")
	fmt.Fprintf(&sb, "dict[k], err = %s(v)\n", ToDeclSyntaxID(val))
	fmt.Fprintln(&sb, "if err != nil { return nu.Value{}, err }")
	fmt.Fprintln(&sb, "}")
	fmt.Fprintln(&sb, "return nu.Value{Value: dict}, nil")

	out = sb.String()
	return
}

// struct types

type structBridge struct {
	t      reflect.Type
	router *BridgeTypeRouter
}

func structRoute(router *BridgeTypeRouter, t reflect.Type) GoNuBridgeType {
	if t.Kind() != reflect.Struct {
		return nil
	}
	return structBridge{t: t, router: router}
}

func (s structBridge) GoType() reflect.Type {
	return s.t
}

func (s structBridge) TypeExpr() string {
	if s.t.Kind() != reflect.Struct {
		panic("cannot use non-struct as struct")
	}

	var sb strings.Builder
	fmt.Fprintln(&sb, "types.RecordDef{")
	for i := range s.t.NumField() {
		f := s.t.Field(i)
		recordField := nuFieldName(s.t, f)
		if recordField == "-" {
			continue
		}
		child := s.router.Lookup(f.Type)
		syntaxID := TypeDeclSyntaxID(child.GoType())
		fmt.Printf("%s %T\n", f.Name, child)
		if _, ok := child.(structBridge); ok {
			syntaxID = fmt.Sprintf("types.Record(%s)", syntaxID)
		}
		fmt.Fprintf(&sb, "\"%s\": %s,\n", recordField, syntaxID)
	}
	fmt.Fprint(&sb, "}")

	return sb.String()
}

func (s structBridge) FromBody() string {
	if s.t.Kind() != reflect.Struct {
		panic("cannot use non-struct as struct")
	}

	var sb strings.Builder

	fmt.Fprintln(&sb, "record, ok := v.Value.(nu.Record)")
	fmt.Fprintf(&sb, "if !ok { return out, fmt.Errorf(\"expected nu.Record got %%T\", v.Value) }\n")
	fmt.Fprintln(&sb, "var val nu.Value")

	for i := range s.t.NumField() {
		f := s.t.Field(i)

		recordField := nuFieldName(s.t, f)
		if recordField == "-" {
			continue
		}
		defaultVal := f.Tag.Get("default")

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

		child := s.router.Lookup(f.Type).GoType()
		fmt.Fprintf(
			&sb,
			// all nullables are represented with a pointer, therefore if the
			// key doesn't exist it will simply return nil, which is the
			// correct behavior
			"out.%s, err = %s(val)\n",
			f.Name,
			FromDeclSyntaxID(child),
		)
		fmt.Fprintln(&sb, "if err != nil { return out, err }")
		if defaultVal != "" {
			fmt.Fprintln(&sb, "}")
		}
	}
	fmt.Fprintln(&sb, "return out, nil")

	return sb.String()
}

func (s structBridge) ToBody() string {
	t := s.GoType()
	if t.Kind() != reflect.Struct {
		panic("cannot use non-struct as struct")
	}

	var sb strings.Builder
	fmt.Fprintln(&sb, "rec := nu.Record{}")
	for i := range t.NumField() {
		f := t.Field(i)

		recordField := nuFieldName(t, f)
		if recordField == "-" {
			continue
		}

		child := s.router.Lookup(f.Type).GoType()
		fmt.Fprintf(&sb, "rec[\"%s\"], err = %s(v.%s)\n", recordField, ToDeclSyntaxID(child), f.Name)
		fmt.Fprintln(&sb, "if err != nil { return nu.Value{}, err }")
	}
	fmt.Fprintln(&sb, "return nu.Value{Value: rec}, nil")

	return sb.String()
}

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
