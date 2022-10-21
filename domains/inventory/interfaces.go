package inventory

import "context"

type InventoryService interface {
	Create(ctx context.Context)
	Get(ctx context.Context)
	Update(ctx context.Context)
}
