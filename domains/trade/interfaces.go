package trade

import "context"

type TradeInterface interface {
	InitiateTrade(ctx context.Context) error
	Match(ctx context.Context) error
}
