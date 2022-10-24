package store

import (
	"context"
	"os"
	"testing"
	"time"

	"zssn/domains/core"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db    *gorm.DB
	store ITradeStorage
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		cleanup()
		os.Exit(code)
	}()

	d, err := setupTestDB()
	if err != nil {
		panic(err)
	}
	db = d
	s, err := New(db)
	if err != nil {
		panic(err)
	}
	store = s
	code = m.Run()
}

func TestNewWithNilDB(t *testing.T) {
	var dd *gorm.DB
	st, err := New(dd)
	require.NotNil(t, err)
	assert.EqualError(t, err, "invalid connection passed")
	require.Nil(t, st)
}

func TestCreateTransaction(t *testing.T) {
	ctx := context.Background()
	sellerID, buyerID := uuid.NewString(), uuid.NewString()
	ti := newTradeItems(t, sellerID)
	it := newTradeItems(t, buyerID)

	err := store.Execute(ctx, ti, it)
	require.NoError(t, err)
	require.NotEmpty(t, ti.Reference)

	res, err := store.Details(ctx, ti.Reference)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	assert.Len(t, res, 6)
}

func TestTradeHistory(t *testing.T) {
	ctx := context.Background()
	sellerID, buyerID := uuid.NewString(), uuid.NewString()

	ti := newTradeItems(t, sellerID)
	it := newTradeItems(t, buyerID)

	err := store.Execute(ctx, ti, it)
	require.NoError(t, err)
	require.NotEmpty(t, ti.Reference)

	// flip the buyer and seller
	ti = newTradeItems(t, buyerID)
	it = newTradeItems(t, sellerID)

	err = store.Execute(ctx, ti, it)
	require.NoError(t, err)
	require.NotEmpty(t, ti.Reference)
	start, end := time.Now().Add(-24*time.Hour), time.Now()

	res, err := store.History(ctx, sellerID, start, end)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestTradeDetails(t *testing.T) {
	ctx := context.Background()
	sellerID, buyerID := uuid.NewString(), uuid.NewString()

	ti := newTradeItems(t, sellerID)
	it := newTradeItems(t, buyerID)

	err := store.Execute(ctx, ti, it)
	require.NoError(t, err)
	require.NotEmpty(t, ti.Reference)

	res, err := store.Details(ctx, ti.Reference)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	assert.Len(t, res, 6)

	for _, v := range res {
		assert.Equal(t, ti.Reference, v.Reference)
	}
}

func newTradeItems(t *testing.T, userID string) *TradeItems {
	t.Helper()
	return &TradeItems{
		UserID: userID,
		Items: []TradeItem{
			{
				Item:     core.ItemWater,
				Quantity: 10,
			},
			{
				Item:     core.ItemAmmunition,
				Quantity: 20,
			},
			{
				Item:     core.ItemMedication,
				Quantity: 30,
			},
		},
	}
}

func setupTestDB() (*gorm.DB, error) {
	env := os.Getenv("ENVIRONMENT")
	dsn := "root:@tcp(127.0.0.1:3306)/zssn?charset=utf8mb4&parseTime=True&loc=Local"
	if env == "cicd" {
		dsn = "zssn_user:password@tcp(127.0.0.1:33306)/zssn?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func cleanup() {
	db.Exec("DELETE FROM transactions")
}
