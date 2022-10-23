package trade

import (
	"context"
	"time"
	"zssn/domains/entities"
)

// ITradeService contract for the business logic around trade management
type ITradeService interface {
	Execute(ctx context.Context, debitUserID, creditUserID string, items *entities.TradeItem) error
	History(ctx context.Context, id string, startDate, endDate time.Time) ([]*entities.Transaction, error)
	VerifyTransaction(ctx context.Context, sellerItem, buyerItem entities.TradeItems) error
}
