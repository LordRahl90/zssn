package store

import (
	"context"
	"time"
)

// ITradeStorage interface for the trade service
type ITradeStorage interface {
	Execute(ctx context.Context, seller, buyer *TradeItems) error
	Details(ctx context.Context, ref string) ([]*Transaction, error)
	History(ctx context.Context, userID string, start, endDate time.Time) ([]*Transaction, error)
}
