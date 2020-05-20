package parsers

import (
	"time"

	"github.com/baking-bad/bcdhub/internal/contractparser/consts"
	"github.com/baking-bad/bcdhub/internal/elastic"
	"github.com/baking-bad/bcdhub/internal/helpers"
	"github.com/baking-bad/bcdhub/internal/metrics"
	"github.com/baking-bad/bcdhub/internal/models"
	"github.com/baking-bad/bcdhub/internal/noderpc"
	"github.com/tidwall/gjson"
)

// MigrationParser -
type MigrationParser struct {
	rpc            noderpc.Pool
	es             *elastic.Elastic
	filesDirectory string
}

// NewMigrationParser -
func NewMigrationParser(rpc noderpc.Pool, es *elastic.Elastic, filesDirectory string) *MigrationParser {
	return &MigrationParser{
		rpc:            rpc,
		es:             es,
		filesDirectory: filesDirectory,
	}
}

// Parse -
func (p *MigrationParser) Parse(data gjson.Result, old models.Contract, prevProtocol, newProtocol models.Protocol) ([]elastic.Model, error) {
	metadata := models.Metadata{ID: old.Address}
	if err := p.es.GetByID(&metadata); err != nil {
		return nil, err
	}

	if err := updateMetadata(data, newProtocol.SymLink, &old, &metadata); err != nil {
		return nil, err
	}

	migrationBlock, err := p.rpc.GetHeader(prevProtocol.EndLevel)
	if err != nil {
		return nil, err
	}

	newFingerprint, err := metrics.GetFingerprint(data)
	if err != nil {
		return nil, err
	}
	if newFingerprint.Compare(old.Fingerprint) {
		return []elastic.Model{&metadata}, nil
	}

	migration := models.Migration{
		ID:          helpers.GenerateID(),
		IndexedTime: time.Now().UnixNano() / 1000,

		Network:      old.Network,
		Level:        prevProtocol.EndLevel,
		Protocol:     newProtocol.Hash,
		PrevProtocol: prevProtocol.Hash,
		Address:      old.Address,
		Timestamp:    migrationBlock.Timestamp,
		Kind:         consts.MigrationUpdate,
	}

	return []elastic.Model{&metadata, &migration}, nil
}
