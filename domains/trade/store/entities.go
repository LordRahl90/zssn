package store

import (
	"zssn/domains/core"

	"gorm.io/gorm"
)

// TradeItem represents a unit of trade item
type TradeItem struct {
	Item     core.Item `json:"item"`
	Quantity uint32    `json:"quantity"`
}

// TradItems collection of trade item
type TradeItems struct {
	UserID    string      `json:"userID"`
	Reference string      `json:"reference"`
	Items     []TradeItem `json:"items"`
}

// Transactions a ledger type of table that keeps a log of all the transactions
type Transaction struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Reference string    `json:"reference"`
	SellerID  string    `json:"user_id"`
	BuyerID   string    `json:"buyer_id"`
	Item      core.Item `json:"item"`
	Quantity  uint32    `json:"credit"`
	gorm.Model
}
