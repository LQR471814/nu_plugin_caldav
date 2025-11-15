package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"os"
	"reflect"

	"github.com/LQR471814/nu_plugin_caldav/internal/nutypes"
)

func code() *Code {
	c := NewCode()
	c.AddImport("github.com/ainvaltin/nu-plugin")
	c.AddImport("github.com/ainvaltin/nu-plugin/types")
	c.AddImport("github.com/LQR471814/nu_plugin_caldav/events")
	c.AddImport("github.com/LQR471814/nu_plugin_caldav/internal/nutypes")
	c.AddImport("github.com/teambition/rrule-go")
	c.AddImport("github.com/emersion/go-webdav/caldav")
	c.AddImport("net/url")
	c.AddImport("time")
	c.Use("EventObjectReplicaList", reflect.TypeOf(nutypes.EventObjectReplicaList{}))
	c.Use("EventObjectReplica", reflect.TypeOf(nutypes.EventObjectReplica{}))
	c.Use("EventReplica", reflect.TypeOf(nutypes.EventReplica{}))
	c.Use("Timeline", reflect.TypeOf(nutypes.Timeline{}))
	c.Use("CalendarList", reflect.TypeOf(nutypes.CalendarList{}))
	return c
}

func render(pkg string) []byte {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: ", err)
		}
	}()

	buf := bytes.NewBuffer(nil)

	code().Render(buf, pkg)

	src, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	return src
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
	f.Write(render(*pkg))
	f.Close()
}
