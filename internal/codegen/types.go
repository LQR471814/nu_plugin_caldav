package main

import (
	"fmt"
	"reflect"

	"github.com/zeebo/xxh3"
)

// Type declarations

// GoNuBridgeType encapsulates all the logic necessary to convert a type
// between golang and nushell
type GoNuBridgeType interface {
	// GoType returns the golang type that will be bridged to nushell
	GoType() reflect.Type
	// TypeExpr returns a golang expression for the nushell type representation
	TypeExpr() string
	// FromBody returns the body of the conversion function from nu -> go
	FromBody() string
	// ToBody returns the body of the conversion function from go -> nu
	ToBody() string
}

// Convenience syntax functions

func TypeDeclSyntaxID(t reflect.Type) string {
	// xxh3 is extremely fast, I wouldn't say this is much slower than doing a
	// hash map lookup
	return fmt.Sprintf("type_%d", xxh3.Hash([]byte(t.String())))
}

func FromDeclSyntaxID(t reflect.Type) string {
	return fmt.Sprintf("type_%d_FromNu", xxh3.Hash([]byte(t.String())))
}

func ToDeclSyntaxID(t reflect.Type) string {
	return fmt.Sprintf("type_%d_ToNu", xxh3.Hash([]byte(t.String())))
}

// Bridge type routing

type CachedBridgeType struct {
	underlying GoNuBridgeType
	goType     reflect.Type
	typeExpr   string
	toBody     string
	fromBody   string
}

func (r CachedBridgeType) GoType() reflect.Type {
	return r.goType
}

func (r CachedBridgeType) TypeExpr() string {
	return r.typeExpr
}

func (r CachedBridgeType) ToBody() string {
	return r.toBody
}

func (r CachedBridgeType) FromBody() string {
	return r.fromBody
}

// BridgeTypeRoute returns a concrete bridge type adapter given a golang type
// or nil if the type is not supported
type BridgeTypeRoute func(router *BridgeTypeRouter, t reflect.Type) GoNuBridgeType

type BridgeTypeRouter struct {
	// Routes is the list of types that are supported, earlier routes take
	// precedence
	Routes []BridgeTypeRoute
	// KnownTypes keeps track of all the types instantiated by the
	// BridgeTypeRouter
	KnownTypes map[string]CachedBridgeType
}

func (r *BridgeTypeRouter) Lookup(t reflect.Type) GoNuBridgeType {
	typestr := t.String()
	existing, ok := r.KnownTypes[typestr]
	if ok {
		return existing.underlying
	}
	for _, route := range r.Routes {
		accepted := route(r, t)
		if accepted != nil {
			// we call TypeExpr(), ToBody(), and FromBody() ahead of time so
			// that all the generated types necessary for this code to work are
			// all known
			r.KnownTypes[typestr] = CachedBridgeType{
				underlying: accepted,
				goType:     t,
				typeExpr:   accepted.TypeExpr(),
				toBody:     accepted.ToBody(),
				fromBody:   accepted.FromBody(),
			}
			return accepted
		}
	}
	panic(fmt.Errorf("unsupported type used! (%v)", t))
}
