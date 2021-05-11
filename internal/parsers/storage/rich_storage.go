package storage

import (
	"github.com/baking-bad/bcdhub/internal/models/operation"
	"github.com/baking-bad/bcdhub/internal/noderpc"
	"github.com/baking-bad/bcdhub/internal/parsers"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// RichStorage -
type RichStorage struct {
	DeffatedStorage []byte
	Result          *parsers.Result
	Empty           bool
}

// Parser -
type Parser interface {
	ParseTransaction(content noderpc.Operation, operation operation.Operation) (RichStorage, error)
	ParseOrigination(content noderpc.Operation, operation operation.Operation) (RichStorage, error)
}
