package entities

import (
	"zssn/domains/core"
	"zssn/domains/trade/store"
)

// TradeItem represents a single trading unit
type TradeItem struct {
	Item     core.Item `json:"item"`
	Quantity uint32    `json:"quantity"`
}

// TradeItems to specify trading item the user is providing
type TradeItems struct {
	UserID    string      `json:"userID"`
	Reference string      `json:"reference"`
	Items     []TradeItem `json:"items"`
}

// Transactions service layer entity
type Transaction struct {
	ID        string    `json:"id"`
	Reference string    `json:"reference"`
	SellerID  string    `json:"user_id"`
	BuyerID   string    `json:"buyer_id"`
	Item      core.Item `json:"item"`
	Quantity  uint32    `json:"credit"`
}

// ToDBTradeItemEntities converts service entities to db entities
func (ti *TradeItems) ToDBTradeItemEntities() *store.TradeItems {
	st := &store.TradeItems{
		UserID: ti.UserID,
	}

	for _, v := range ti.Items {
		st.Items = append(st.Items, store.TradeItem{
			Item:     v.Item,
			Quantity: v.Quantity,
		})
	}
	return st
}

// FromDBTransactionEntity converts repo/store entities to service entities
func FromDBTransactionEntity(m *store.Transaction) *Transaction {
	return &Transaction{
		ID:        m.ID,
		Reference: m.Reference,
		SellerID:  m.SellerID,
		BuyerID:   m.BuyerID,
		Item:      m.Item,
		Quantity:  m.Quantity,
	}
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
