package enrich

import "fmt"

// FieldNotFoundErr is returned when a property of name Field cannot be
// found on the event while reading a Field
type FieldNotFoundErr struct {
	Field string
}

func (e FieldNotFoundErr) Error() string {
	return fmt.Sprintf("field (%s) not found", e.Field)
}

// FieldNotFoundErr is returned when a property of name Field failed to be
// parsed while reading a Field
type FieldParseErr struct {
	Field string
	Inner error
}

func (e FieldParseErr) Error() string {
	return fmt.Sprintf("parse field (%s): %s", e.Field, e.Inner.Error())
}

func (e FieldParseErr) Unwrap() error {
	return e.Inner
}

// FieldSerializeErr is returned when a property of name Field failed to
// be serialized while writing a Field
type FieldSerializeErr struct {
	Field string
	Inner error
}

func (e FieldSerializeErr) Error() string {
	return fmt.Sprintf("serialize field (%s): %s", e.Field, e.Inner.Error())
}

func (e FieldSerializeErr) Unwrap() error {
	return e.Inner
}
