package trade

import (
	"context"
	"errors"
	"time"
	"zssn/domains/trade/store"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	mockDB = make(map[string][]*store.Transaction)

	errMockNotInitialized = errors.New("mock not initialized")
)

type MockTradeStore struct {
	DetailsFunc func(ctx context.Context, ref string) ([]*store.Transaction, error)
	ExecuteFunc func(ctx context.Context, item *store.TradeItems) error
	HistoryFunc func(ctx context.Context, userID string, start time.Time, endDate time.Time) ([]*store.Transaction, error)
}

// NewStoreMock returns a new mock for storage trade
func NewStoreMock() store.ITradeStorage {
	return &MockTradeStore{
		ExecuteFunc: func(ctx context.Context, item *store.TradeItems) error {
			item.Reference = uuid.NewString()
			ref := uuid.NewString()
			var trans []*store.Transaction
			for _, v := range item.Items {
				// create a transaction record for every item
				trans = append(trans, &store.Transaction{
					ID:        uuid.NewString(),
					Reference: ref,
					SellerID:  item.Seller,
					BuyerID:   item.Buyer,
					Item:      v.Item,
					Quantity:  v.Quantity,
				})
			}
			item.Reference = ref
			mockDB[ref] = trans
			return nil
		},
		DetailsFunc: func(ctx context.Context, ref string) ([]*store.Transaction, error) {
			v, ok := mockDB[ref]
			if !ok {
				return nil, gorm.ErrRecordNotFound
			}
			return v, nil
		},
	}
}

// Details implements store.ITradeStorage
func (m *MockTradeStore) Details(ctx context.Context, ref string) ([]*store.Transaction, error) {
	if m.DetailsFunc == nil {
		return nil, errMockNotInitialized
	}
	return m.DetailsFunc(ctx, ref)
}

// Execute implements store.ITradeStorage
func (m *MockTradeStore) Execute(ctx context.Context, item *store.TradeItems) error {
	if m.ExecuteFunc == nil {
		return errMockNotInitialized
	}
	return m.ExecuteFunc(ctx, item)
}

// History implements store.ITradeStorage
func (m *MockTradeStore) History(ctx context.Context, userID string, start time.Time, endDate time.Time) ([]*store.Transaction, error) {
	if m.HistoryFunc == nil {
		return nil, errMockNotInitialized
	}
	return m.HistoryFunc(ctx, userID, start, endDate)
}
