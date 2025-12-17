package main

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

// All declarations

type Code struct {
	Imports []string
	Aliases map[string]reflect.Type
	Router  *BridgeTypeRouter
}

func NewCode() *Code {
	return &Code{
		Aliases: make(map[string]reflect.Type),
		Router: &BridgeTypeRouter{
			Routes: []BridgeTypeRoute{
				// custom types
				rruleRoute,
				urlRoute,
				timestampRoute,
				durationRoute,

				// primitive types
				structRoute,
				mapRoute,
				sliceRoute,
				pointerRoute,
				primitiveRoute,
			},
			KnownTypes: make(map[string]CachedBridgeType),
		},
	}
}

// User-facing fucntions

func (c *Code) AddImport(imp string) {
	c.Imports = append(c.Imports, imp)
}

func (c *Code) Use(typename string, t reflect.Type) {
	c.Router.Lookup(t)
	c.Aliases[typename] = t
}

// Rendering / generating functions

func (*Code) fmtTypeDecl(out io.Writer, t reflect.Type, expr string) {
	fmt.Fprintf(out, "var %s = %s\n", TypeDeclSyntaxID(t), expr)
}

func (*Code) fmtFromDecl(out io.Writer, t reflect.Type, body string) {
	returnType := t.String()
	fmt.Fprintf(out, "func %s(v nu.Value) (out %s, err error) {\n", FromDeclSyntaxID(t), returnType)
	fmt.Fprintf(out, "defer func() { if err != nil { err = fmt.Errorf(\"%s: %%w\", err) } }()\n", returnType)
	out.Write([]byte(body))
	fmt.Fprintln(out, "}")
}

func (*Code) fmtToDecl(out io.Writer, t reflect.Type, body string) {
	paramType := t.String()
	fmt.Fprintf(out, "func %s(v %s) (out nu.Value, err error) {\n", ToDeclSyntaxID(t), paramType)
	fmt.Fprintf(out, "defer func() { if err != nil { err = fmt.Errorf(\"%s: %%w\", err) } }()\n", paramType)
	out.Write([]byte(body))
	fmt.Fprintln(out, "}")
}

func (c *Code) Render(out io.Writer, pkg string) {
	fmt.Fprintf(out, "package %s\n", pkg)
	for _, imp := range c.Imports {
		fmt.Fprintf(out, "import \"%s\"\n", imp)
	}

	for _, bt := range c.Router.KnownTypes {
		c.fmtTypeDecl(out, bt.GoType(), bt.TypeExpr())
		c.fmtFromDecl(out, bt.GoType(), bt.FromBody())
		c.fmtToDecl(out, bt.GoType(), bt.ToBody())
	}
	for alias, typ := range c.Aliases {
		fmt.Fprintf(out, "var %sType = %s\n", alias, TypeDeclSyntaxID(typ))
		fmt.Fprintf(out, "var %sFromNu = %s\n", alias, FromDeclSyntaxID(typ))
		fmt.Fprintf(out, "var %sToNu = %s\n", alias, ToDeclSyntaxID(typ))
	}
}

// convenience functions

func pascalToSnakeCase(pascalcase string) string {
	var result strings.Builder
	for i, c := range pascalcase {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				result.WriteString("_")
			}
			result.WriteRune(c + ('a' - 'A'))
			continue
		}
		result.WriteRune(c)
	}
	return result.String()
}
