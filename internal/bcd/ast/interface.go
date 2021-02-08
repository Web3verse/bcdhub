package ast

import (
	"fmt"

	"github.com/baking-bad/bcdhub/internal/bcd/base"
)

// Node -
type Node interface {
	fmt.Stringer
	Type
	Value
	Base
}

// Type -
type Type interface {
	ParseType(node *base.Node, id *int) error
	GetPrim() string
	GetName() string
	IsPrim(prim string) bool
	IsNamed() bool
	GetEntrypoints() []string
	ToJSONSchema() (*JSONSchema, error)
	Docs(inferredName string) ([]Typedef, string, error)
	FindByName(name string) Node
}

// Value -
type Value interface {
	ParseValue(node *base.Node) error
	GetValue() interface{}
	ToMiguel() (*MiguelNode, error)
	FromJSONSchema(data map[string]interface{}) error
	EnrichBigMap(bmd []*base.BigMapDiff) error
	ToParameters() ([]byte, error)
}

// Base -
type Base interface {
	ToBaseNode(optimized bool) (*base.Node, error)
}
