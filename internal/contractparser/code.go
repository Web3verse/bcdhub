package contractparser

import (
	"fmt"

	"github.com/aopoltorzhicky/bcdhub/internal/contractparser/consts"
	"github.com/aopoltorzhicky/bcdhub/internal/contractparser/node"
	"github.com/aopoltorzhicky/bcdhub/internal/helpers"
	"github.com/aopoltorzhicky/bcdhub/internal/tlsh"
	"github.com/tidwall/gjson"
)

// Code -
type Code struct {
	*parser

	Parameter Parameter
	Storage   Storage
	Code      gjson.Result

	Hash string
	hash []byte

	Tags        helpers.Set
	Language    string
	FailStrings helpers.Set
	Primitives  helpers.Set
	Annotations helpers.Set
}

func newCode(script gjson.Result) (Code, error) {
	res := Code{
		parser:      &parser{},
		Language:    consts.LangUnknown,
		hash:        make([]byte, 0),
		FailStrings: make(helpers.Set),
		Primitives:  make(helpers.Set),
		Tags:        make(helpers.Set),
		Annotations: make(helpers.Set),
	}
	res.primHandler = res.handlePrimitive
	res.arrayHandler = res.handleArray

	code := script.Get("code").Array()
	if len(code) != 3 {
		return res, fmt.Errorf("Invalid tag 'code' length")
	}

	for i := range code {
		n := node.NewNodeJSON(code[i])
		if err := res.parseStruct(n); err != nil {
			return res, err
		}
	}
	return res, nil
}

func (c *Code) parseStruct(n node.Node) error {
	switch n.Prim {
	case consts.CODE:
		c.Code = n.Args
		if err := c.parseCode(n.Args); err != nil {
			return err
		}
	case consts.STORAGE:
		store, err := newStorage(n.Args)
		if err != nil {
			return err
		}
		c.Storage = store
	case consts.PARAMETER:
		parameter, err := newParameter(n.Args)
		if err != nil {
			return err
		}
		c.Parameter = parameter
	}

	return nil
}

func (c *Code) parseCode(args gjson.Result) error {
	for _, val := range args.Array() {
		if err := c.parse(val); err != nil {
			return err
		}
	}

	if len(c.hash) == 0 {
		c.hash = append(c.hash, 0)
	}
	h, err := tlsh.HashBytes(c.hash)
	if err != nil {
		return err
	}
	c.Hash = h.String()
	return nil
}

func (c *Code) handleArray(arr []gjson.Result) error {
	if fail := parseFail(arr); fail != nil {
		if fail.With != "" {
			c.FailStrings.Append(fail.With)
		}
	}
	return nil
}

func (c *Code) handlePrimitive(n node.Node) (err error) {
	c.Primitives.Append(n.Prim)
	c.hash = append(c.hash, []byte(n.Prim)...)

	if n.HasAnnots() {
		c.Annotations.Append(n.Annotations...)
	}

	c.detectLanguage(n)
	c.Tags.Append(primTags(n))

	return nil
}

func (c *Code) detectLanguage(n node.Node) {
	if c.Language != consts.LangUnknown {
		return
	}
	if detectLiquidity(n, c.Parameter.Entrypoints()) {
		c.Language = consts.LangLiquidity
		return
	}
	if detectPython(n) {
		c.Language = consts.LangPython
		return
	}
	if detectLIGO(n) {
		c.Language = consts.LangLigo
		return
	}
	if c.Language == "" {
		c.Language = consts.LangUnknown
	}
}