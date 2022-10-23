package store

import (
	"context"
	"os"
	"testing"
	"zssn/domains/core"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db      *gorm.DB
	storage IInventoryStorage
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
	storage = s
	code = m.Run()
}

func TestNewStoreImplementation(t *testing.T) {
	st, err := New(db)
	require.NoError(t, err)
	assert.NotNil(t, st)
}

func TestStoreWithNilDB(t *testing.T) {
	var emptyDB *gorm.DB
	st, err := New(emptyDB)
	require.NotNil(t, err)
	require.EqualError(t, err, "invalid db provided")
	assert.Nil(t, st)
}

func TestCreateNewInventoryRecord(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	invs := newInventory(t, userID)
	err := storage.Create(ctx, invs)
	require.NoError(t, err)
	for _, v := range invs {
		require.NotEmpty(t, v.ID)
	}
}

func TestCreateNewDuplicateRecord(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	invs := newInventory(t, userID)
	err := storage.Create(ctx, invs)
	require.NoError(t, err)
	for _, v := range invs {
		require.NotEmpty(t, v.ID)
		require.NotEmpty(t, v.Balance)
		require.True(t, v.Accessible)
	}

	addition := []*Inventory{
		{
			UserID:   userID,
			Item:     core.ItemWater,
			Quantity: 15,
		},
	}
	err = storage.Create(ctx, addition)
	require.NotNil(t, err)
}

func TestFindUserInventory(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	invs := newInventory(t, userID)

	anotherUser := uuid.NewString()
	aInv := newInventory(t, anotherUser)

	err := storage.Create(ctx, invs)
	require.NoError(t, err)

	err = storage.Create(ctx, aInv)
	require.NoError(t, err)

	res, err := storage.FindUserInventory(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, res)

	for _, v := range invs {
		inv, ok := res[v.Item]
		require.True(t, ok)
		require.NotNil(t, inv)
		assert.Equal(t, inv.ID, v.ID)
		assert.Equal(t, inv.Quantity, v.Quantity)
		assert.Equal(t, inv.Balance, v.Balance)
		assert.Equal(t, inv.Accessible, v.Accessible)
	}
}

func TestFindUserInventoryForNonExistingUser(t *testing.T) {
	ctx := context.Background()

	res, err := storage.FindUserInventory(ctx, uuid.NewString())
	require.NoError(t, err)
	assert.Empty(t, res)
}

func TestFindUserInventoryForManyUsers(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	invs := newInventory(t, userID)

	anotherUserID := uuid.NewString()
	aInv := newInventory(t, anotherUserID)

	err := storage.Create(ctx, invs)
	require.NoError(t, err)

	err = storage.Create(ctx, aInv)
	require.NoError(t, err)

	res, err := storage.FindUsersInventory(ctx, userID, anotherUserID)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Len(t, res, 2)

	userRes, ok := res[userID]
	assert.True(t, ok)

	anotherUserRes, ok := res[anotherUserID]
	assert.True(t, ok)

	for _, v := range invs {
		inv, ok := userRes[v.Item]
		require.True(t, ok)
		require.NotNil(t, inv)
		assert.Equal(t, inv.ID, v.ID)
		assert.Equal(t, inv.Quantity, v.Quantity)
		assert.Equal(t, inv.Balance, v.Balance)
		assert.Equal(t, inv.Accessible, v.Accessible)
	}

	for _, v := range aInv {
		inv, ok := anotherUserRes[v.Item]
		require.True(t, ok)
		require.NotNil(t, inv)
		assert.Equal(t, inv.ID, v.ID)
		assert.Equal(t, inv.Quantity, v.Quantity)
		assert.Equal(t, inv.Balance, v.Balance)
		assert.Equal(t, inv.Accessible, v.Accessible)
	}
}

func TestUpdateBalance(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	invs := newInventory(t, userID)

	err := storage.Create(ctx, invs)
	require.NoError(t, err)

	err = storage.UpdateBalance(ctx, userID, core.ItemWater, 50)
	require.NoError(t, err)

	err = storage.UpdateBalance(ctx, userID, core.ItemMedication, 500)
	require.NoError(t, err)

	res, err := storage.FindUserInventory(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, uint32(50), res[core.ItemWater].Balance)
	assert.Equal(t, uint32(500), res[core.ItemMedication].Balance)
}

func TestUpdateInventoryAccessibility(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	invs := newInventory(t, userID)

	err := storage.Create(ctx, invs)
	require.NoError(t, err)

	err = storage.UpdateUserInventoryAccessibility(ctx, userID)
	require.NoError(t, err)

	res, err := storage.FindUserInventory(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, res)

	for _, v := range res {
		assert.False(t, v.Accessible)
	}
}

func newInventory(t *testing.T, userID string) []*Inventory {
	t.Helper()
	return []*Inventory{
		{
			UserID:   userID,
			Item:     core.ItemWater,
			Quantity: 20,
		},
		{
			UserID:   userID,
			Item:     core.ItemFood,
			Quantity: 20,
		},
		{
			UserID:   userID,
			Item:     core.ItemMedication,
			Quantity: 30,
		},
		{
			UserID:   userID,
			Item:     core.ItemAmmunition,
			Quantity: 50,
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
	db.Exec("DELETE FROM inventories")
}
