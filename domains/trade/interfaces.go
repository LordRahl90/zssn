package trade

import (
	"context"
	"time"
	"zssn/domains/entities"
)

// ITradeService contract for the business logic around trade management
type ITradeService interface {
	Execute(ctx context.Context, seller, buyer *entities.TradeItems) error
	History(ctx context.Context, id string, startDate, endDate time.Time) ([]*entities.Transaction, error)
	IsTransactionAmountEqual(sellerItem, buyerItem *entities.TradeItems) error
	AnyParticipantInfected(users ...*entities.User) error
	EnoughStock(stock entities.Stock, items *entities.TradeItems) error
	VerifyTransaction(ctx context.Context, sellerItem, buyerItem *entities.TradeItems) error
}
