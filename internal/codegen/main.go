package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/LQR471814/nu_plugin_caldav/internal/enrich"
	"github.com/LQR471814/nu_plugin_caldav/internal/enrich/dto"
)

func underlyingType(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Array, reflect.Map:
		return underlyingType(t.Elem())
	}
	return t
}

func newCode() (c *Code, err error) {
	c = NewCode()

	// c.AddImport("net/url")
	c.AddImport("time")
	c.AddImport("fmt")
	c.AddImport("github.com/ainvaltin/nu-plugin")
	c.AddImport("github.com/ainvaltin/nu-plugin/types")
	c.AddImport("github.com/LQR471814/nu_plugin_caldav/internal/enrich/props")
	c.AddImport("github.com/LQR471814/nu_plugin_caldav/internal/enrich/dto")
	c.AddImport("github.com/teambition/rrule-go")
	c.AddImport("github.com/emersion/go-webdav/caldav")

	for _, reg := range enrich.AllFields() {
		typ := underlyingType(reflect.TypeOf(reg.Zero))
		if !strings.HasSuffix(typ.PkgPath(), "enrich/props") {
			continue
		}
		c.Use(typ.Name(), typ)
	}

	c.Use("PropValue", reflect.TypeFor[dto.PropValueList]())
	c.Use("PropValueList", reflect.TypeFor[dto.PropValueList]())
	c.Use("CalendarList", reflect.TypeFor[dto.CalendarList]())
	c.Use("TimeSegment", reflect.TypeFor[dto.TimeSegment]())
	c.Use("Timeline", reflect.TypeFor[dto.Timeline]())
	c.Use("RRule", reflect.TypeFor[dto.RRule]())

	return
}

func render(pkg string) (out []byte, err error) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: ", err, string(debug.Stack()))
		}
	}()

	buf := bytes.NewBuffer(nil)

	code, err := newCode()
	if err != nil {
		return
	}
	code.Render(buf, pkg)

	src, err := format.Source(buf.Bytes())
	if err != nil {
		err = fmt.Errorf("format failed: %v\n%v", err, buf.String())
		return
	}
	out = src
	return
}

func main() {
	out := flag.String("out", "", "The file to output generated go code to.")
	pkg := flag.String("pkg", "main", "The package name to use in the generated go code.")
	flag.Parse()

	if *out == "" {
		fmt.Fprintln(os.Stderr, "Error: must specify an output directory via -out")
		os.Exit(1)
	}

	f, err := os.Create(*out)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}
	defer f.Close()

	contents, err := render(*pkg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(1)
	}
	f.Write(contents)
}
