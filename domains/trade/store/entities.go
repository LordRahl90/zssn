package store

import "zssn/domains/core"

// TradeItem represents a unit of trade item
type TradeItem struct {
	ID       string    `json:"id"`
	Item     core.Item `json:"item"`
	Quantity uint32    `json:"quantity"`
}

// TradItems collection of trade item
type TradItems []TradeItem

// Calculate calculates a collection of trade items based on their points and quantity
func (t TradItems) Calculate() (result uint32) {
	for _, v := range t {
		pts, ok := core.ItemPoints[v.Item]
		if !ok {
			continue
		}
		result += (pts * v.Quantity)
	}
	return
}
