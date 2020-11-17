package operations

import (
	"strings"
	"time"

	"github.com/baking-bad/bcdhub/internal/contractparser/consts"
	"github.com/baking-bad/bcdhub/internal/contractparser/meta"
	"github.com/baking-bad/bcdhub/internal/elastic"
	"github.com/baking-bad/bcdhub/internal/helpers"
	"github.com/baking-bad/bcdhub/internal/models"
	"github.com/tidwall/gjson"
)

// Transaction -
type Transaction struct {
	*ParseParams
}

// NewTransaction -
func NewTransaction(params *ParseParams) Transaction {
	return Transaction{params}
}

// Parse -
func (p Transaction) Parse(data gjson.Result) ([]elastic.Model, error) {
	tx := models.Operation{
		ID:            helpers.GenerateID(),
		Network:       p.network,
		Hash:          p.hash,
		Protocol:      p.head.Protocol,
		Level:         p.head.Level,
		Timestamp:     p.head.Timestamp,
		Kind:          data.Get("kind").String(),
		Initiator:     data.Get("source").String(),
		Source:        data.Get("source").String(),
		Fee:           data.Get("fee").Int(),
		Counter:       data.Get("counter").Int(),
		GasLimit:      data.Get("gas_limit").Int(),
		StorageLimit:  data.Get("storage_limit").Int(),
		Amount:        data.Get("amount").Int(),
		Destination:   data.Get("destination").String(),
		PublicKey:     data.Get("public_key").String(),
		ManagerPubKey: data.Get("manager_pubkey").String(),
		Delegate:      data.Get("delegate").String(),
		Parameters:    data.Get("parameters").String(),
		IndexedTime:   time.Now().UnixNano() / 1000,
		ContentIndex:  p.contentIdx,
	}

	p.fillInternal(&tx)

	if data.Get("nonce").Exists() {
		nonce := data.Get("nonce").Int()
		tx.Nonce = &nonce
	}

	txMetadata := parseMetadata(data, tx)
	tx.Result = &txMetadata.Result
	tx.Status = tx.Result.Status
	tx.Errors = tx.Result.Errors

	tx.SetBurned(p.constants)
	txModels := []elastic.Model{&tx}

	if tx.IsApplied() {
		for i := range txMetadata.BalanceUpdates {
			txModels = append(txModels, txMetadata.BalanceUpdates[i])
		}

		appliedModels, err := p.appliedHandler(data, &tx)
		if err != nil {
			return nil, err
		}
		txModels = append(txModels, appliedModels...)
	}

	if err := p.tagTransaction(&tx); err != nil {
		return nil, err
	}

	p.stackTrace.Add(tx)

	transfers, err := p.transferParser.Parse(tx)
	if err != nil {
		return nil, err
	}
	for i := range transfers {
		txModels = append(txModels, transfers[i])
	}

	if err := p.createTokenBalanceUpdates(transfers); err != nil {
		return nil, err
	}

	return txModels, nil
}

func (p Transaction) fillInternal(tx *models.Operation) {
	if p.main == nil {
		p.main = tx
		return
	}

	tx.Counter = p.main.Counter
	tx.Hash = p.main.Hash
	tx.Level = p.main.Level
	tx.Timestamp = p.main.Timestamp
	tx.Internal = true
	tx.Initiator = p.main.Source
}

func (p Transaction) appliedHandler(item gjson.Result, op *models.Operation) ([]elastic.Model, error) {
	if !helpers.IsContract(op.Destination) || !op.IsApplied() {
		return nil, nil
	}

	metadata, err := meta.GetContractMetadata(p.es, op.Destination)
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			return nil, nil
		}
		return nil, err
	}

	resultModels := make([]elastic.Model, 0)

	rs, err := p.storageParser.Parse(item, metadata, op)
	if err != nil {
		return nil, err
	}
	if rs.Empty {
		return nil, err
	}
	op.DeffatedStorage = rs.DeffatedStorage

	resultModels = append(resultModels, rs.Models...)

	migration := NewMigration(op).Parse(item)
	if migration != nil {
		resultModels = append(resultModels, migration)
	}

	bu := NewBalanceUpdate("metadata", *op).Parse(item)
	for i := range bu {
		resultModels = append(resultModels, bu[i])
	}
	return resultModels, p.getEntrypoint(item, metadata, op)
}

func (p Transaction) getEntrypoint(item gjson.Result, metadata *meta.ContractMetadata, op *models.Operation) error {
	m, err := metadata.Get(consts.PARAMETER, op.Protocol)
	if err != nil {
		return err
	}

	params := item.Get("parameters")
	if params.Exists() {
		ep, err := m.GetByPath(params)
		if err != nil && op.Errors == nil {
			return err
		}
		op.Entrypoint = ep
	} else {
		op.Entrypoint = consts.DefaultEntrypoint
	}

	return nil
}

func (p Transaction) tagTransaction(tx *models.Operation) error {
	if !helpers.IsContract(tx.Destination) {
		return nil
	}

	contract := models.NewEmptyContract(tx.Network, tx.Destination)
	if err := p.es.GetByID(&contract); err != nil {
		if elastic.IsRecordNotFound(err) {
			return nil
		}
		return err
	}
	tx.Tags = make([]string, 0)
	for _, tag := range contract.Tags {
		if helpers.StringInArray(tag, []string{
			consts.FA12Tag, consts.FA2Tag,
		}) {
			tx.Tags = append(tx.Tags, tag)
		}
	}
	return nil
}

func (p Transaction) createTokenBalanceUpdates(transfers []*models.Transfer) error {
	exists := make(map[string]*models.TokenBalance)
	updates := make([]*models.TokenBalance, 0)
	for i := range transfers {
		idFrom := transfers[i].GetFromTokenBalanceID()
		if idFrom != "" {
			if update, ok := exists[idFrom]; ok {
				update.Balance -= int64(transfers[i].Amount)
			} else {
				upd := transfers[i].MakeTokenBalanceUpdate(true, false)
				updates = append(updates, upd)
				exists[idFrom] = upd
			}
		}
		idTo := transfers[i].GetToTokenBalanceID()
		if idTo != "" {
			if update, ok := exists[idTo]; ok {
				update.Balance += int64(transfers[i].Amount)
			} else {
				upd := transfers[i].MakeTokenBalanceUpdate(false, false)
				updates = append(updates, upd)
				exists[idTo] = upd
			}
		}
	}

	return p.es.UpdateTokenBalances(updates)
}