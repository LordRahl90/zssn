package store

import (
	"context"
	"time"
)

// ITradeService interface to manage trades
type ITradeService interface {
	Execute(ctx context.Context, debitUserID, creditUserID string, items TradeItem) error
	History(ctx context.Context, id string, startDate, endDate time.Time) ([]Transactions, error)
}
