package inventory

import (
	"context"
	"errors"

	"zssn/domains/core"
	"zssn/domains/inventory/store"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	_ store.IInventoryStorage = (*MockInventoryStore)(nil)

	mockStore = make(map[string]store.Response)

	errMockNotInitialized = errors.New("mock not initialized")
)

// MockInventoryStore inventory store mock
type MockInventoryStore struct {
	CreateFunc                           func(ctx context.Context, items []*store.Inventory) error
	FindUserInventoryFunc                func(ctx context.Context, userID string) (store.Response, error)
	FindUsersInventoryFunc               func(ctx context.Context, userIDs ...string) (map[string]store.Response, error)
	UpdateBalanceFunc                    func(ctx context.Context, userID string, item core.Item, newBalance uint32) error
	UpdateUserInventoryAccessibilityFunc func(ctx context.Context, userID string) error
	ReduceBalanceFunc                    func(ctx context.Context, userID string, item core.Item, qty uint32) error
	UpdateMultipleBalanceFunc            func(ctx context.Context, userID string, items map[core.Item]uint32) error
}

// NewMockStore return a new mock store with prefilled functions using mockStore
func NewMockStore() *MockInventoryStore {
	return &MockInventoryStore{
		CreateFunc: func(ctx context.Context, items []*store.Inventory) error {
			var res = make(store.Response)
			for _, v := range items {
				v.ID = uuid.NewString()
				v.Accessible = true
				v.Balance = v.Quantity
				res[v.Item] = v
			}
			mockStore[items[0].UserID] = res
			return nil
		},
		FindUserInventoryFunc: func(ctx context.Context, userID string) (store.Response, error) {
			res, ok := mockStore[userID]
			if !ok {
				return nil, gorm.ErrRecordNotFound
			}
			return res, nil
		},
		FindUsersInventoryFunc: func(ctx context.Context, userIDs ...string) (map[string]store.Response, error) {
			result := make(map[string]store.Response)
			for _, v := range userIDs {
				data, ok := mockStore[v]
				if !ok {
					continue
				}
				result[v] = data
			}

			return result, nil
		},
		UpdateBalanceFunc: func(ctx context.Context, userID string, item core.Item, newBalance uint32) error {
			data, ok := mockStore[userID]
			if !ok {
				return nil
			}
			data[item].Balance = newBalance
			mockStore[userID] = data
			return nil
		},
		UpdateUserInventoryAccessibilityFunc: func(ctx context.Context, userID string) error {
			data, ok := mockStore[userID]
			if !ok {
				return gorm.ErrRecordNotFound
			}
			for _, v := range data {
				v.Accessible = false
			}
			mockStore[userID] = data
			return nil
		},
	}
}

// Create implements store.IInventoryStorage
func (m *MockInventoryStore) Create(ctx context.Context, items []*store.Inventory) error {
	if m.CreateFunc == nil {
		return errMockNotInitialized
	}
	return m.CreateFunc(ctx, items)
}

// FindUserInventory implements store.IInventoryStorage
func (m *MockInventoryStore) FindUserInventory(ctx context.Context, userID string) (store.Response, error) {
	if m.FindUserInventoryFunc == nil {
		return nil, errMockNotInitialized
	}

	return m.FindUserInventoryFunc(ctx, userID)
}

// FindUsersInventory implements store.IInventoryStorage
func (m *MockInventoryStore) FindUsersInventory(ctx context.Context, userIDs ...string) (map[string]store.Response, error) {
	if m.FindUsersInventoryFunc == nil {
		return nil, errMockNotInitialized
	}

	return m.FindUsersInventoryFunc(ctx, userIDs...)
}

// UpdateBalance implements store.IInventoryStorage
func (m *MockInventoryStore) UpdateBalance(ctx context.Context, userID string, item core.Item, newBalance uint32) error {
	if m.UpdateBalanceFunc == nil {
		return errMockNotInitialized
	}
	return m.UpdateBalanceFunc(ctx, userID, item, newBalance)
}

// ReduceBalance implements store.IInventoryStorage
func (m *MockInventoryStore) ReduceBalance(ctx context.Context, userID string, item core.Item, qty uint32) error {
	if m.ReduceBalanceFunc == nil {
		return errMockNotInitialized
	}

	return m.ReduceBalanceFunc(ctx, userID, item, qty)
}

func (m *MockInventoryStore) UpdateMultipleBalance(ctx context.Context, userID string, items map[core.Item]uint32) error {
	if m.UpdateMultipleBalanceFunc == nil {
		return errMockNotInitialized
	}
	return m.UpdateMultipleBalanceFunc(ctx, userID, items)
}

// UpdateUserInventoryAccessibility implements store.IInventoryStorage
func (m *MockInventoryStore) UpdateUserInventoryAccessibility(ctx context.Context, userID string) error {
	if m.UpdateUserInventoryAccessibilityFunc == nil {
		return errMockNotInitialized
	}
	return m.UpdateUserInventoryAccessibilityFunc(ctx, userID)
}
