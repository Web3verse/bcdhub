package ast

import (
	"strings"

	"github.com/baking-bad/bcdhub/internal/bcd/base"
	"github.com/baking-bad/bcdhub/internal/bcd/consts"
	"github.com/pkg/errors"
)

//  TODO: pack/unpack

// Pair -
type Pair struct {
	Default
	Args []Node
}

// NewPair -
func NewPair(depth int) *Pair {
	return &Pair{
		Default: NewDefault(consts.PAIR, -1, depth),
	}
}

// String -
func (p *Pair) String() string {
	var s strings.Builder
	s.WriteString(p.Default.String())
	for i := range p.Args {
		s.WriteString(strings.Repeat(base.DefaultIndent, p.depth))
		s.WriteString(p.Args[i].String())
	}
	return s.String()
}

// MarshalJSON -
func (p *Pair) MarshalJSON() ([]byte, error) {
	return marshalJSON(consts.PAIR, p.annots, p.Args...)
}

// ParseType -
func (p *Pair) ParseType(node *base.Node, id *int) error {
	if err := p.Default.ParseType(node, id); err != nil {
		return err
	}

	p.Args = make([]Node, 0)
	if len(node.Args) == 2 {
		for _, arg := range node.Args {
			child, err := typingNode(arg, p.depth, id)
			if err != nil {
				return err
			}
			p.Args = append(p.Args, child)
		}
	} else if len(node.Args) > 2 {
		child, err := typingNode(node.Args[0], p.depth, id)
		if err != nil {
			return err
		}
		p.Args = append(p.Args, child)

		newUntyped := &base.Node{
			Prim: consts.PAIR,
			Args: node.Args[1:],
		}
		pairChild, err := typingNode(newUntyped, p.depth+1, id)
		if err != nil {
			return err
		}
		p.Args = append(p.Args, pairChild)
	}

	return nil
}

// ParseValue -
func (p *Pair) ParseValue(node *base.Node) error {
	switch {
	case len(node.Args) == 2:
		for i := range p.Args {
			if err := p.Args[i].ParseValue(node.Args[i]); err != nil {
				return err
			}
		}
		return nil
	case len(node.Args) > 2:
		if err := p.Args[0].ParseValue(node.Args[0]); err != nil {
			return err
		}

		newUntyped := &base.Node{
			Prim: consts.PAIR,
			Args: node.Args[1:],
		}
		return p.Args[1].ParseValue(newUntyped)
	default:
		return errors.Wrap(base.ErrInvalidArgsCount, "Pair.ParseValue")
	}
}

// ToMiguel -
func (p *Pair) ToMiguel() (*MiguelNode, error) {
	node := &MiguelNode{
		Prim:     p.Prim,
		Type:     consts.TypeNamedTuple,
		Name:     p.GetName(),
		Children: make([]*MiguelNode, 0),
	}

	for i := range p.Args {
		child, err := p.Args[i].ToMiguel()
		if err != nil {
			return nil, err
		}

		if p.Prim == p.Args[i].GetPrim() {
			node.Children = append(node.Children, child.Children...)
		} else {
			node.Children = append(node.Children, child)
		}
	}

	return node, nil
}

// ToBaseNode -
func (p *Pair) ToBaseNode(optimized bool) (*base.Node, error) {
	node := new(base.Node)
	node.Prim = p.Prim
	node.Args = make([]*base.Node, 0)
	for i := range p.Args {
		arg, err := p.Args[i].ToBaseNode(optimized)
		if err != nil {
			return nil, err
		}
		node.Args = append(node.Args, arg)
	}
	return node, nil
}

// ToJSONSchema -
func (p *Pair) ToJSONSchema() (*JSONSchema, error) {
	s := &JSONSchema{
		Type:       JSONSchemaTypeObject,
		Properties: map[string]*JSONSchema{},
	}
	for _, arg := range p.Args {
		if err := setChildSchema(arg, false, s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// FromJSONSchema -
func (p *Pair) FromJSONSchema(data map[string]interface{}) error {
	for i := range p.Args {
		if err := p.Args[i].FromJSONSchema(data); err != nil {
			return err
		}
	}
	return nil
}