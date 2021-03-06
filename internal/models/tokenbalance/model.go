package tokenbalance

import (
	"github.com/baking-bad/bcdhub/internal/models/types"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TokenBalance -
type TokenBalance struct {
	ID            int64           `json:"-" gorm:"autoIncrement:true"`
	Network       types.Network   `json:"network" gorm:"not null;primaryKey;index:token_balances_token_idx;default:0"`
	Address       string          `json:"address" gorm:"not null;primaryKey"`
	Contract      string          `json:"contract" gorm:"not null;primaryKey;index:token_balances_token_idx"`
	TokenID       uint64          `json:"token_id" gorm:"type:numeric(50,0);default:0;primaryKey;autoIncrement:false;index:token_balances_token_idx"`
	Balance       decimal.Decimal `json:"balance" gorm:"type:numeric(100,0);default:0"`
	BalanceString string          `json:"balance_string"`

	IsLedger bool `json:"-" gorm:"-"`
}

// GetID -
func (tb *TokenBalance) GetID() int64 {
	return tb.ID
}

// GetIndex -
func (tb *TokenBalance) GetIndex() string {
	return "token_balances"
}

// Constraint -
func (tb *TokenBalance) Save(tx *gorm.DB) error {
	var s clause.Set

	if tb.IsLedger {
		s = clause.Assignments(map[string]interface{}{
			"balance":        tb.Balance,
			"balance_string": tb.Balance.String(),
		})
	} else {
		s = clause.Assignments(map[string]interface{}{
			"balance":        gorm.Expr("token_balances.balance + ?", tb.Balance),
			"balance_string": gorm.Expr("(token_balances.balance + ?)::text", tb.Balance),
		})
	}

	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "network"},
			{Name: "contract"},
			{Name: "address"},
			{Name: "token_id"},
		},
		DoUpdates: s,
	}).Save(tb).Error
}

// GetQueues -
func (tb *TokenBalance) GetQueues() []string {
	return nil
}

// MarshalToQueue -
func (tb *TokenBalance) MarshalToQueue() ([]byte, error) {
	return nil, nil
}

// LogFields -
func (tb *TokenBalance) LogFields() logrus.Fields {
	return logrus.Fields{
		"network":  tb.Network.String(),
		"address":  tb.Address,
		"contract": tb.Contract,
		"token_id": tb.TokenID,
	}
}

// BeforeCreate -
func (tb *TokenBalance) BeforeCreate(tx *gorm.DB) (err error) {
	return tb.marshal()
}

// BeforeUpdate -
func (tb *TokenBalance) BeforeUpdate(tx *gorm.DB) (err error) {
	return tb.marshal()
}

func (tb *TokenBalance) marshal() error {
	tb.BalanceString = tb.Balance.String()
	return nil
}
