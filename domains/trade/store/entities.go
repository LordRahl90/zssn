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
	Seller    string      `json:"seller"`
	Buyer     string      `json:"buyer"`
	Reference string      `json:"reference"`
	Items     []TradeItem `json:"items"`
}

// Calculate calculates a collection of trade items based on their points and quantity
func (t TradeItems) Calculate() (result uint32) {
	for _, v := range t.Items {
		pts, ok := core.ItemPoints[v.Item]
		if !ok {
			continue
		}
		result += (pts * v.Quantity)
	}
	return
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
