package mocks

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

	_ store.ITradeStorage = (*MockTradeStore)(nil)
)

type MockTradeStore struct {
	DetailsFunc func(ctx context.Context, ref string) ([]*store.Transaction, error)
	ExecuteFunc func(ctx context.Context, seller, buyer *store.TradeItems) error
	HistoryFunc func(ctx context.Context, userID string, start time.Time, endDate time.Time) ([]*store.Transaction, error)
}

// NewStoreMock returns a new mock for storage trade
func NewStoreMock() store.ITradeStorage {
	return &MockTradeStore{
		ExecuteFunc: func(ctx context.Context, seller, buyer *store.TradeItems) error {
			ref := uuid.NewString()
			var trans []*store.Transaction
			for _, v := range seller.Items {
				// create a transaction record for every item
				trans = append(trans, &store.Transaction{
					ID:        uuid.NewString(),
					Reference: ref,
					SellerID:  seller.UserID,
					BuyerID:   buyer.UserID,
					Item:      v.Item,
					Quantity:  v.Quantity,
				})
			}
			for _, v := range buyer.Items {
				// seller becomes buyer at this point and the total is calculated into inventory service
				trans = append(trans, &store.Transaction{
					ID:        uuid.NewString(),
					Reference: ref,
					SellerID:  buyer.UserID,
					BuyerID:   seller.UserID,
					Item:      v.Item,
					Quantity:  v.Quantity,
				})
			}
			seller.Reference = ref
			buyer.Reference = ref
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
		HistoryFunc: func(ctx context.Context, userID string, start, endDate time.Time) ([]*store.Transaction, error) {
			var result []*store.Transaction
			for _, trans := range mockDB {
				for _, v := range trans {
					if v.BuyerID == userID || v.SellerID == userID {
						result = append(result, v)
					}
				}
			}

			return result, nil
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
func (m *MockTradeStore) Execute(ctx context.Context, seller, buyer *store.TradeItems) error {
	if m.ExecuteFunc == nil {
		return errMockNotInitialized
	}
	return m.ExecuteFunc(ctx, seller, buyer)
}

// History implements store.ITradeStorage
func (m *MockTradeStore) History(ctx context.Context, userID string, start time.Time, endDate time.Time) ([]*store.Transaction, error) {
	if m.HistoryFunc == nil {
		return nil, errMockNotInitialized
	}
	return m.HistoryFunc(ctx, userID, start, endDate)
}
