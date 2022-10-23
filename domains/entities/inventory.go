package entities

import (
	"zssn/domains/core"
	"zssn/domains/inventory/store"
)

// Inventory DTO object for transferring invetory items
type Inventory struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	Item       core.Item `json:"item"`
	Quantity   uint32    `json:"quantity"`
	Balance    uint32    `json:"balance"`
	Accessible bool      `json:"-"`
}

// ToInventoryDBEntity converts from service entity to db entity
func (i *Inventory) ToInventoryDBEntity() *store.Inventory {
	return &store.Inventory{
		ID:     i.ID,
		UserID: i.UserID,
		Item:   i.Item,
	}
}

// FromInventoryDBEntity converts from db entity to service entity
func FromInventoryDBEntity(m *store.Inventory) *Inventory {
	return &Inventory{
		ID:         m.ID,
		UserID:     m.UserID,
		Item:       m.Item,
		Quantity:   m.Quantity,
		Balance:    m.Balance,
		Accessible: m.Accessible,
	}
}
