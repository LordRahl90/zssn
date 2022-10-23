package store

import (
	"zssn/domains/core"

	"gorm.io/gorm"
)

// Response type for search responses
type Response map[core.Item]*Inventory

// Inventory contains the mapping for inventory storage
type Inventory struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	UserID     string    `json:"user_id" gorm:"size:50;index:idx_user_item,unique"`
	Item       core.Item `json:"item" gorm:"index:idx_user_item,unique"`
	Quantity   uint32    `json:"quantity"`
	Balance    uint32    `json:"balance"`
	Accessible bool      `json:"-"`
	gorm.Model
}
