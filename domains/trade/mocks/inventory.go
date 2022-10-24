package mocks

import (
	"context"

	"zssn/domains/core"
	"zssn/domains/entities"
	"zssn/domains/inventory"
)

var _ inventory.IInventoryService = (*MockInventoryService)(nil)

// MockInventoryService mock for IInventoryService
type MockInventoryService struct {
	BlockUserInventoryFunc    func(ctx context.Context, userID string) error
	CreateFunc                func(ctx context.Context, item []*entities.Inventory) error
	FindMultipleInventoryFunc func(ctx context.Context, userIDs ...string) (entities.UserStock, error)
	FindUserInventoryFunc     func(ctx context.Context, userID string) (map[string]*entities.Inventory, error)
	UpdateBalanceFunc         func(ctx context.Context, userID string, item core.Item, newBalance uint32) error
}

// NewInventoryMock returns a new mock implementation using in-memory db
func NewInventoryMock() inventory.IInventoryService {
	return inventory.New(inventory.NewMockStore())
}

// BlockUserInventory implements inventory.IInventoryService
func (m *MockInventoryService) BlockUserInventory(ctx context.Context, userID string) error {
	if m.BlockUserInventoryFunc == nil {
		return errMockNotDefined
	}
	return m.BlockUserInventoryFunc(ctx, userID)
}

// Create implements inventory.IInventoryService
func (m *MockInventoryService) Create(ctx context.Context, item []*entities.Inventory) error {
	if m.CreateFunc == nil {
		return errMockNotDefined
	}

	return m.CreateFunc(ctx, item)
}

// FindMultipleInventory implements inventory.IInventoryService
func (m *MockInventoryService) FindMultipleInventory(ctx context.Context, userIDs ...string) (entities.UserStock, error) {
	if m.FindMultipleInventoryFunc == nil {
		return nil, errMockNotDefined
	}
	return m.FindMultipleInventoryFunc(ctx, userIDs...)
}

// FindUserInventory implements inventory.IInventoryService
func (m *MockInventoryService) FindUserInventory(ctx context.Context, userID string) (map[string]*entities.Inventory, error) {
	if m.FindUserInventoryFunc == nil {
		return nil, errMockNotDefined
	}
	return m.FindUserInventoryFunc(ctx, userID)
}

// UpdateBalance implements inventory.IInventoryService
func (m *MockInventoryService) UpdateBalance(ctx context.Context, userID string, item core.Item, newBalance uint32) error {
	if m.UpdateBalanceFunc == nil {
		return errMockNotDefined
	}
	return m.UpdateBalanceFunc(ctx, userID, item, newBalance)
}
