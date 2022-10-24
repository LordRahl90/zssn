package inventory

import (
	"context"
	"zssn/domains/core"
	"zssn/domains/entities"
	"zssn/domains/inventory/store"
)

// InventoryService contains an implementation of IInventoryService
type InventoryService struct {
	store store.IInventoryStorage
}

// New returns a new implementation of IInventoryService
func New(storage store.IInventoryStorage) IInventoryService {
	return &InventoryService{
		store: storage,
	}
}

// Create implements IInventoryService
func (iv *InventoryService) Create(ctx context.Context, items []*entities.Inventory) error {
	var dbItems []*store.Inventory
	for _, v := range items {
		dbItems = append(dbItems, v.ToInventoryDBEntity())
	}
	err := iv.store.Create(ctx, dbItems)
	if err != nil {
		return err
	}
	for i, v := range dbItems {
		items[i].ID = v.ID
	}
	return nil
}

// FindMultipleInventory implements IInventoryService
func (iv *InventoryService) FindMultipleInventory(ctx context.Context, userIDs ...string) (entities.UserStock, error) {
	result := make(entities.UserStock)
	res, err := iv.store.FindUsersInventory(ctx, userIDs...)
	if err != nil {
		return nil, err
	}
	for k, v := range res {
		ent := make(entities.Stock)
		for i, j := range v {
			ent[i] = entities.FromInventoryDBEntity(j)
		}
		result[k] = ent
	}
	return result, nil
}

// FindUserInventory implements IInventoryService
func (iv *InventoryService) FindUserInventory(ctx context.Context, userID string) (map[core.Item]*entities.Inventory, error) {
	res, err := iv.store.FindUserInventory(ctx, userID)
	result := make(map[core.Item]*entities.Inventory)
	if err != nil {
		return nil, err
	}
	for k, v := range res {
		result[k] = entities.FromInventoryDBEntity(v)
	}
	return result, nil
}

// UpdateBalance implements IInventoryService
func (iv *InventoryService) UpdateBalance(ctx context.Context, userID string, item core.Item, newBalance uint32) error {
	return iv.store.UpdateBalance(ctx, userID, item, newBalance)
}

// BlockUserInventory implements IInventoryService
func (iv *InventoryService) BlockUserInventory(ctx context.Context, userID string) error {
	return iv.store.UpdateUserInventoryAccessibility(ctx, userID)
}
