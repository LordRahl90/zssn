package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TradeStore implementation of ITradeStorage
type TradeStorage struct {
	DB *gorm.DB
}

// New returns a new implementation of trade storage
func New(db *gorm.DB) (ITradeStorage, error) {
	if db == nil {
		return nil, fmt.Errorf("invalid connection passed")
	}
	if err := db.AutoMigrate(&Transaction{}); err != nil {
		return nil, err
	}
	return &TradeStorage{
		DB: db,
	}, nil
}

// Execute implements ITradeStorage
func (ts *TradeStorage) Execute(ctx context.Context, item *TradeItems) error {
	ref := uuid.NewString()
	var trans []Transaction
	for _, v := range item.Items {
		// create a transaction record for every item
		trans = append(trans, Transaction{
			ID:        uuid.NewString(),
			Reference: ref,
			SellerID:  item.Seller,
			BuyerID:   item.Buyer,
			Item:      v.Item,
			Quantity:  v.Quantity,
		})
	}
	item.Reference = ref
	return ts.DB.Create(trans).Error
}

// History returns the trade history for a particular user within a timeframe, we still want accountability even in an apocalypse :)
func (ts *TradeStorage) History(ctx context.Context, userID string, start time.Time, endDate time.Time) ([]*Transaction, error) {
	var result []*Transaction
	err := ts.DB.Debug().Where("(seller_id = ? OR buyer_id = ?) AND DATE(created_at) BETWEEN DATE(?) AND DATE(?)", userID, userID, start, endDate).Find(&result).Error

	return result, err
}

// Details returns the details of a given transaction
func (ts *TradeStorage) Details(ctx context.Context, ref string) ([]*Transaction, error) {
	var result []*Transaction
	err := ts.DB.Where("reference = ?", ref).Find(&result).Error

	return result, err
}
