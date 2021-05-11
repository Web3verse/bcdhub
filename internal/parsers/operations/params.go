package operations

import (
	"sync"

	"github.com/baking-bad/bcdhub/internal/cache"
	"github.com/baking-bad/bcdhub/internal/config"
	"github.com/baking-bad/bcdhub/internal/logger"
	"github.com/baking-bad/bcdhub/internal/models/operation"
	"github.com/baking-bad/bcdhub/internal/models/protocol"
	"github.com/baking-bad/bcdhub/internal/noderpc"
	"github.com/baking-bad/bcdhub/internal/parsers/contract"
	"github.com/baking-bad/bcdhub/internal/parsers/stacktrace"
	"github.com/baking-bad/bcdhub/internal/parsers/transfer"
)

// ParseParams -
type ParseParams struct {
	ctx *config.Context

	rpc      noderpc.INode
	shareDir string

	constants protocol.Constants

	contractParser *contract.Parser
	transferParser *transfer.Parser

	storageParser *RichStorage

	stackTrace *stacktrace.StackTrace

	ipfs []string

	network    string
	hash       string
	head       noderpc.Header
	contentIdx int64
	main       *operation.Operation
	cache      *cache.Cache

	once *sync.Once
}

// ParseParamsOption -
type ParseParamsOption func(*ParseParams)

// WithIPFSGateways -
func WithIPFSGateways(ipfs []string) ParseParamsOption {
	return func(dp *ParseParams) {
		dp.ipfs = ipfs
	}
}

// WithConstants -
func WithConstants(constants protocol.Constants) ParseParamsOption {
	return func(dp *ParseParams) {
		dp.constants = constants
	}
}

// WithShareDirectory -
func WithShareDirectory(shareDir string) ParseParamsOption {
	return func(dp *ParseParams) {
		dp.shareDir = shareDir
	}
}

// WithNetwork -
func WithNetwork(network string) ParseParamsOption {
	return func(dp *ParseParams) {
		dp.network = network
	}
}

// WithHash -
func WithHash(hash string) ParseParamsOption {
	return func(dp *ParseParams) {
		dp.hash = hash
	}
}

// WithHead -
func WithHead(head noderpc.Header) ParseParamsOption {
	return func(dp *ParseParams) {
		dp.head = head
	}
}

// WithContentIndex -
func WithContentIndex(index int64) ParseParamsOption {
	return func(dp *ParseParams) {
		dp.contentIdx = index
	}
}

// WithMainOperation -
func WithMainOperation(main *operation.Operation) ParseParamsOption {
	return func(dp *ParseParams) {
		dp.main = main
	}
}

// WithCache -
func WithCache(cache *cache.Cache) ParseParamsOption {
	return func(dp *ParseParams) {
		dp.cache = cache
	}
}

// NewParseParams -
func NewParseParams(rpc noderpc.INode, ctx *config.Context, opts ...ParseParamsOption) *ParseParams {
	params := &ParseParams{
		ctx:        ctx,
		rpc:        rpc,
		once:       &sync.Once{},
		stackTrace: stacktrace.New(),
	}
	for i := range opts {
		opts[i](params)
	}

	transferParser, err := transfer.NewParser(
		params.rpc,
		ctx.TZIP, ctx.Blocks, ctx.Storage, params.shareDir,
		transfer.WithStackTrace(params.stackTrace),
		transfer.WithNetwork(params.network),
		transfer.WithChainID(params.head.ChainID),
		transfer.WithGasLimit(params.constants.HardGasLimitPerOperation),
	)
	if err != nil {
		logger.Error(err)
	}
	params.transferParser = transferParser

	params.contractParser = contract.NewParser(
		params.ctx,
		contract.WithShareDir(params.shareDir),
	)
	storageParser, err := NewRichStorage(ctx.BigMapDiffs, rpc, params.head.Protocol)
	if err != nil {
		logger.Error(err)
	}
	params.storageParser = storageParser

	if params.cache == nil {
		params.cache = cache.NewCache()
	}
	return params
}
