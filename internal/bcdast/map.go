package bcdast

import (
	"fmt"
	"strings"

	"github.com/baking-bad/bcdhub/internal/contractparser/consts"
	"github.com/pkg/errors"
)

// Map -
type Map struct {
	Default
	KeyType   AstNode
	ValueType AstNode

	Data map[AstNode]AstNode
}

// NewMap -
func NewMap(depth int) *Map {
	return &Map{
		Default: NewDefault(consts.MAP, 2, depth),
	}
}

// String -
func (m *Map) String() string {
	var s strings.Builder

	s.WriteString(m.Default.String())
	if len(m.Data) > 0 {
		for key, val := range m.Data {
			s.WriteString(strings.Repeat(indent, m.depth))
			s.WriteByte('{')
			s.WriteByte('\n')
			s.WriteString(strings.Repeat(indent, m.depth+1))
			s.WriteString(key.String())
			s.WriteString(strings.Repeat(indent, m.depth+1))
			s.WriteString(val.String())
			s.WriteString(strings.Repeat(indent, m.depth))
			s.WriteByte('}')
			s.WriteByte('\n')
		}
	} else {
		s.WriteString(strings.Repeat(indent, m.depth))
		s.WriteString(m.KeyType.String())
		s.WriteString(strings.Repeat(indent, m.depth))
		s.WriteString(m.ValueType.String())
	}

	return s.String()
}

// MarshalJSON -
func (m *Map) MarshalJSON() ([]byte, error) {
	return marshalJSON(consts.MAP, m.annots, m.KeyType, m.ValueType)
}

// ParseType -
func (m *Map) ParseType(untyped Untyped, id *int) error {
	if err := m.Default.ParseType(untyped, id); err != nil {
		return err
	}

	keyType, err := typingNode(untyped.Args[0], m.depth, id)
	if err != nil {
		return err
	}
	m.KeyType = keyType

	valType, err := typingNode(untyped.Args[1], m.depth, id)
	if err != nil {
		return err
	}
	m.ValueType = valType

	return nil
}

// ParseValue -
func (m *Map) ParseValue(untyped Untyped) error {
	if untyped.Prim != PrimArray {
		return errors.Wrap(ErrInvalidPrim, "Map.ParseValue")
	}

	data, err := createMapFromElts(untyped.Args, m.KeyType, m.ValueType)
	if err != nil {
		return err
	}
	m.Data = data

	return nil
}

// ToMiguel -
func (m *Map) ToMiguel() (*MiguelNode, error) {
	node, err := m.Default.ToMiguel()
	if err != nil {
		return nil, err
	}

	node.Children = make([]*MiguelNode, 0)
	for key, value := range m.Data {
		keyChild, err := key.ToMiguel()
		if err != nil {
			return nil, err
		}
		if keyChild != nil {
			child, err := value.ToMiguel()
			if err != nil {
				return nil, err
			}

			name, err := getMapKeyName(keyChild)
			if err != nil {
				return nil, err
			}
			child.Name = name
			node.Children = append(node.Children, child)
		}
	}

	return node, nil
}

func createMapFromElts(args []Untyped, keyType, valueType AstNode) (map[AstNode]AstNode, error) {
	data := make(map[AstNode]AstNode)

	for i := range args {
		elt := args[i]
		if elt.Prim != consts.Elt {
			return nil, errors.Wrap(ErrInvalidPrim, "BigMap.ParseValue")
		}
		if len(elt.Args) != 2 {
			return nil, errors.Wrap(ErrInvalidArgsCount, "BigMap.ParseValue")
		}
		key, err := createByType(keyType)
		if err != nil {
			return nil, err
		}
		if err := key.ParseValue(elt.Args[0]); err != nil {
			return nil, err
		}
		val, err := createByType(valueType)
		if err != nil {
			return nil, err
		}
		if err := val.ParseValue(elt.Args[1]); err != nil {
			return nil, err
		}

		data[key] = val
	}
	return data, nil
}

func getMapKeyName(node *MiguelNode) (s string, err error) {
	switch kv := node.Value.(type) {
	case string:
		if kv == "" {
			kv = `""`
		}
		s = kv
	case int, int64:
		s = fmt.Sprintf("%d", kv)
	case bool:
		s = fmt.Sprintf("%t", kv)
	case map[string]interface{}:
		s = fmt.Sprintf("%v", kv["miguel_value"])
	case []interface{}:
		s = ""
		for i, item := range kv {
			val := item.(map[string]interface{})
			if i != 0 {
				s += "@"
			}
			s += fmt.Sprintf("%v", val["miguel_value"])
		}
	default:
		err = errors.Errorf("Invalid map key type: %v %T", node, node)
	}
	return
}
