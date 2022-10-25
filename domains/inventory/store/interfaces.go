package store

import (
	"context"
	"zssn/domains/core"
)

// IInventoryStore inventory store interface
type IInventoryStorage interface {
	Create(ctx context.Context, items []*Inventory) error
	FindUserInventory(ctx context.Context, userID string) (Response, error)
	FindUsersInventory(ctx context.Context, userIDs ...string) (map[string]Response, error)
	UpdateBalance(ctx context.Context, userID string, item core.Item, newBalance uint32) error
	UpdateMultipleBalance(ctx context.Context, userID string, items map[core.Item]uint32) error
	UpdateUserInventoryAccessibility(ctx context.Context, userID string) error
}
