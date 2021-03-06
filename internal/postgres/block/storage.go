package block

import (
	"github.com/baking-bad/bcdhub/internal/models"
	"github.com/baking-bad/bcdhub/internal/models/block"
	"github.com/baking-bad/bcdhub/internal/models/types"
	"github.com/baking-bad/bcdhub/internal/postgres/core"
)

// Storage -
type Storage struct {
	*core.Postgres
}

// NewStorage -
func NewStorage(pg *core.Postgres) *Storage {
	return &Storage{pg}
}

// Get -
func (storage *Storage) Get(network types.Network, level int64) (block block.Block, err error) {
	err = storage.DB.Table(models.DocBlocks).Scopes(core.Network(network)).Where("level = ?", level).Limit(1).Scan(&block).Error
	return
}

// Last - returns current indexer state for network
func (storage *Storage) Last(network types.Network) (block block.Block, err error) {
	err = storage.DB.Table(models.DocBlocks).Scopes(core.Network(network)).Order("id DESC").Limit(1).Scan(&block).Error
	if storage.IsRecordNotFound(err) {
		err = nil
		block.Network = network
		return
	}
	return
}

// LastByNetworks - return last block for all networks
func (storage *Storage) LastByNetworks() (response []block.Block, err error) {
	subQuery := storage.DB.Table(models.DocBlocks).Select("MAX(id) as id").Group("network")
	err = storage.DB.Table(models.DocBlocks).Where("id IN (?)", subQuery).Find(&response).Error
	return
}

// GetNetworkAlias -
func (storage *Storage) GetNetworkAlias(chainID string) (string, error) {
	var network types.Network
	err := storage.DB.Table(models.DocBlocks).
		Select("network").
		Where("chain_id = ?", chainID).
		Limit(1).
		Scan(&network).Error

	return network.String(), err
}
