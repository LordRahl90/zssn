package inventory

import (
	"context"
	"os"
	"testing"
	"zssn/domains/core"
	"zssn/domains/entities"
	"zssn/domains/inventory/store"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	storage store.IInventoryStorage
	service IInventoryService
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		os.Exit(code)
	}()
	storage = NewMockStore()
	service = New(storage)
	code = m.Run()
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	inv := newInventory(t, userID)

	err := service.Create(ctx, inv)
	require.NoError(t, err)

	for _, v := range inv {
		assert.NotEmpty(t, v.ID)
	}
}

func TestCreateWithInvalidService(t *testing.T) {
	store := &MockInventoryStore{}
	svc := New(store)

	ctx := context.Background()
	userID := uuid.NewString()
	inv := newInventory(t, userID)

	err := svc.Create(ctx, inv)
	require.EqualError(t, err, errMockNotInitialized.Error())
}

func TestFindUserInventory(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	inv := newInventory(t, userID)

	err := service.Create(ctx, inv)
	require.NoError(t, err)

	anotherUserID := uuid.NewString()
	aInv := newInventory(t, anotherUserID)

	err = service.Create(ctx, aInv)
	require.NoError(t, err)

	res, err := service.FindUserInventory(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, res, 4)
	assert.Equal(t, res[core.ItemWater.String()].UserID, userID)
}

func TestFindMultipleUserInventories(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	inv := newInventory(t, userID)

	err := service.Create(ctx, inv)
	require.NoError(t, err)

	anotherUserID := uuid.NewString()
	aInv := newInventory(t, anotherUserID)

	err = service.Create(ctx, aInv)
	require.NoError(t, err)

	result, err := service.FindMultipleInventory(ctx, userID, anotherUserID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)

	userRec, ok := result[userID]
	assert.True(t, ok)
	assert.NotNil(t, userRec)
	assert.Len(t, userRec, 4)
}

func TestUpdateUserItemBalance(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	inv := newInventory(t, userID)

	err := service.Create(ctx, inv)
	require.NoError(t, err)

	err = service.UpdateBalance(ctx, userID, core.ItemAmmunition, 5000)
	require.NoError(t, err)

	res, err := service.FindUserInventory(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, res, 4)
	assert.Equal(t, uint32(5000), res[core.ItemAmmunition.String()].Balance)
}

func TestUpdateUserItemBalanceWithBadMock(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	inv := newInventory(t, userID)
	emptyStore := &MockInventoryStore{}
	fakeMockSVC := New(emptyStore)

	err := service.Create(ctx, inv)
	require.NoError(t, err)

	err = fakeMockSVC.UpdateBalance(ctx, userID, core.ItemAmmunition, 5000)
	require.EqualError(t, err, errMockNotInitialized.Error())
}

func TestBlockUserInventory(t *testing.T) {
	ctx := context.Background()

	userID := uuid.NewString()
	inv := newInventory(t, userID)
	err := service.Create(ctx, inv)
	require.NoError(t, err)

	anotherUserID := uuid.NewString()
	aInv := newInventory(t, anotherUserID)
	err = service.Create(ctx, aInv)
	require.NoError(t, err)

	err = service.BlockUserInventory(ctx, userID)
	require.NoError(t, err)

	result, err := service.FindMultipleInventory(ctx, userID, anotherUserID)
	res := result[userID]
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, res, 4)
	for _, v := range res {
		assert.False(t, v.Accessible)
	}

	// let's make sure we did not accidentally block another user
	res = result[anotherUserID]
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, res, 4)
	for _, v := range res {
		assert.Equal(t, anotherUserID, v.UserID)
		assert.True(t, v.Accessible)
	}
}

func TestBlockUserAccessWithEmptyMock(t *testing.T) {
	ctx := context.Background()
	userID := uuid.NewString()
	inv := newInventory(t, userID)
	emptyStore := &MockInventoryStore{}
	fakeMockSVC := New(emptyStore)

	err := service.Create(ctx, inv)
	require.NoError(t, err)

	err = fakeMockSVC.BlockUserInventory(ctx, userID)
	require.EqualError(t, err, errMockNotInitialized.Error())
}

func newInventory(t *testing.T, userID string) []*entities.Inventory {
	t.Helper()
	return []*entities.Inventory{
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
