package main

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

// Type declarations

type TypeDecl struct {
	TypeId uint64
	Value  string
}

func TypeDeclId(typeId uint64) string {
	return fmt.Sprintf("type_%d", typeId)
}

func (d TypeDecl) String() string {
	return fmt.Sprintf("var %s = %s", TypeDeclId(d.TypeId), d.Value)
}

// From nu value declarations

type FromDecl struct {
	TypeId  uint64
	TypeStr string
	Body    string
}

func (f *FromDecl) SetTypeStr(s string) {
	if s == "" {
		panic("expect non-empty string for TypeStr")
	}
	f.TypeStr = s
}

func FromDeclId(typeId uint64) string {
	return fmt.Sprintf("type_%d_FromNu", typeId)
}

func (d FromDecl) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "func %s(v nu.Value) %s {\n", FromDeclId(d.TypeId), d.TypeStr)
	fmt.Fprint(&sb, d.Body)
	fmt.Fprint(&sb, "}")
	return sb.String()
}

// To nu value declarations

type ToDecl struct {
	TypeId  uint64
	TypeStr string
	Body    string
}

func (f *ToDecl) SetTypeStr(s string) {
	if s == "" {
		panic("expect non-empty string for TypeStr")
	}
	f.TypeStr = s
}

func ToDeclId(typeId uint64) string {
	return fmt.Sprintf("type_%d_ToNu", typeId)
}

func (d ToDecl) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "func %s(v %s) nu.Value {\n", ToDeclId(d.TypeId), d.TypeStr)
	fmt.Fprint(&sb, d.Body)
	fmt.Fprint(&sb, "}")
	return sb.String()
}

// All declarations

type Code struct {
	Imports []string
	Aliases map[string]uint64
	Types   map[uint64]TypeDecl
	Froms   map[uint64]FromDecl
	Tos     map[uint64]ToDecl
}

func NewCode() *Code {
	return &Code{
		Aliases: make(map[string]uint64),
		Types:   make(map[uint64]TypeDecl),
		Froms:   make(map[uint64]FromDecl),
		Tos:     make(map[uint64]ToDecl),
	}
}

func (d *Code) AddImport(imp string) {
	d.Imports = append(d.Imports, imp)
}

func (d *Code) Use(typename string, t reflect.Type) {
	typeDecl(d.Types, t)
	fromDecl(d.Froms, t)
	toDecl(d.Tos, t)
	d.Aliases[typename] = typeId(t)
}

// Rendering / generating functions

func (d Code) Render(out io.Writer, pkg string) {
	fmt.Fprintf(out, "package %s\n", pkg)
	for _, imp := range d.Imports {
		fmt.Fprintf(out, "import \"%s\"\n", imp)
	}
	for _, t := range d.Types {
		fmt.Fprintln(out, t.String())
	}
	for _, f := range d.Froms {
		fmt.Fprintln(out, f.String())
	}
	for _, t := range d.Tos {
		fmt.Fprintln(out, t.String())
	}
	for alias, id := range d.Aliases {
		fmt.Fprintf(out, "var %sType = %s\n", alias, TypeDeclId(id))
		fmt.Fprintf(out, "var %sFromNu = %s\n", alias, FromDeclId(id))
		fmt.Fprintf(out, "var %sToNu = %s\n", alias, ToDeclId(id))
	}
}
