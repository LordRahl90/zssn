package store

import (
	"context"
	"fmt"
	"zssn/domains/core"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InventoryStore inventory store implementing IInventoryStore
type InventoryStore struct {
	DB *gorm.DB
}

// New creates a new instance of IInventoryStorage
func New(db *gorm.DB) (IInventoryStorage, error) {
	if db == nil {
		return nil, fmt.Errorf("invalid db provided")
	}
	if err := db.AutoMigrate(&Inventory{}); err != nil {
		return nil, err
	}
	return &InventoryStore{
		DB: db,
	}, nil
}

// Create implements IInventoryStore
func (inv *InventoryStore) Create(ctx context.Context, items []*Inventory) error {
	for _, v := range items {
		// make sure the original quantity is the same as the balance.
		// Also that the items are accessible by default
		v.ID = uuid.NewString()
		v.Balance = v.Quantity
		v.Accessible = true
	}
	return inv.DB.Create(&items).Error
}

// FindUsersInventory implements IInventoryStorage
func (inv *InventoryStore) FindUsersInventory(ctx context.Context, userIDs ...string) (map[string]Response, error) {
	var (
		res    []*Inventory
		result = make(map[string]Response)
	)
	err := inv.DB.WithContext(ctx).Where("user_id IN (?)", userIDs).Find(&res).Error
	for _, v := range res {
		r, ok := result[v.UserID]
		if !ok {
			r = make(Response)
		}
		_, ok = r[v.Item]
		if !ok {
			r[v.Item] = v
		}
		result[v.UserID] = r
	}
	return result, err
}

// FindUserInventory implements IInventoryStore
func (inv *InventoryStore) FindUserInventory(ctx context.Context, userID string) (Response, error) {
	var (
		res    []*Inventory
		result = make(map[core.Item]*Inventory)
	)
	err := inv.DB.Debug().WithContext(ctx).Where("user_id = ?", userID).Find(&res).Error
	if err != nil {
		return nil, err
	}
	for _, v := range res {
		result[v.Item] = v
	}
	return result, nil
}

// UpdateBalance implements IInventoryStore
func (inv *InventoryStore) UpdateBalance(ctx context.Context, userID string, item core.Item, newBalance uint32) error {
	return inv.DB.Model(&Inventory{}).Where("user_id = ? AND item = ?", userID, item).Update("balance", newBalance).Error
}

// UpdateUserInventoryAccessibility implements IInventoryStore
func (inv *InventoryStore) UpdateUserInventoryAccessibility(ctx context.Context, userID string) error {
	return inv.DB.Model(&Inventory{}).Where("user_id = ?", userID).Update("accessible", false).Error
}
