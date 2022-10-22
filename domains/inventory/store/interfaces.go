package store

import "context"

// IInventoryStore inventory store interface
type IInventoryStore interface {
	Create(ctx context.Context) error
}
