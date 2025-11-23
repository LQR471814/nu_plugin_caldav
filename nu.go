package main

import (
	"fmt"
	"io"
	"reflect"

	"github.com/ainvaltin/nu-plugin"
)

func tryCast[T any](val nu.Value) (T, error) {
	cast, ok := val.Value.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("expected type %T, got %T", zero, val.Value)
	}
	return cast, nil
}

func recvListInput[E any](call *nu.ExecCommand, mapping func(v nu.Value) (E, error)) (list []E, err error) {
	switch typed := call.Input.(type) {
	case nil:
		err = fmt.Errorf("cannot receive null as input")
		return
	case io.ReadCloser:
		err = fmt.Errorf("cannot receive raw stream as input")
		return
	case nu.Value:
		switch v := typed.Value.(type) {
		case []nu.Value:
			list = make([]E, len(v))
			for i, e := range v {
				list[i], err = mapping(e)
				if err != nil {
					return
				}
			}
			return
		case nu.Record:
			var mapped E
			mapped, err = mapping(typed)
			if err != nil {
				return
			}
			list = []E{mapped}
			return
		default:
			err = fmt.Errorf("unknown input type: %v", v)
			return
		}
	case <-chan nu.Value:
		for v := range typed {
			var mapped E
			mapped, err = mapping(v)
			if err != nil {
				return
			}
			list = append(list, mapped)
		}
		return
	}
	panic(fmt.Errorf(
		"unexpected type in input: %v",
		reflect.TypeOf(call.Input),
	))
}
