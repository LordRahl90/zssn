package requests

import (
	"zssn/domains/core"
	"zssn/domains/entities"
)

// TradeItems collection of trade details for users
type TradeItems struct {
	UserID    string      `json:"userID"`
	Reference string      `json:"reference"`
	Items     []TradeItem `json:"items"`
}

// TradeItem represents a single trading unit
type TradeItem struct {
	Item     core.Item `json:"item"`
	Quantity uint32    `json:"quantity"`
}

// TradeRequest is a sample construct trade for both parties
type TradeRequest struct {
	Owner       *TradeItems `json:"originator"`
	SecondParty *TradeItems `json:"second_party"`
}

// ToServiceEntities convert request entities to service entities
func (t *TradeItems) ToServiceEntities() *entities.TradeItems {
	res := &entities.TradeItems{
		UserID: t.UserID,
	}
	for _, v := range t.Items {
		res.Items = append(res.Items, entities.TradeItem{
			Item:     v.Item,
			Quantity: v.Quantity,
		})
	}
	return res
}
