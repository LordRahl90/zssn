package inventory

import (
	"context"
	"zssn/domains/core"
	"zssn/domains/entities"
)

type IInventoryService interface {
	Create(ctx context.Context, item []*entities.Inventory) error
	FindUserInventory(ctx context.Context, userID string) (map[string]*entities.Inventory, error)
	FindMultipleInventory(ctx context.Context, userIDs ...string) (entities.UserStock, error)
	BlockUserInventory(ctx context.Context, userID string) error
	UpdateBalance(ctx context.Context, userID string, item core.Item, newBalance uint32) error
	UpdateMultipleBalance(ctx context.Context, userID string, items map[core.Item]uint32) error
}
