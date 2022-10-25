package trade

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"zssn/domains/core"
	"zssn/domains/entities"
	"zssn/domains/inventory"
	invStore "zssn/domains/inventory/store"
	"zssn/domains/trade/mocks"
	"zssn/domains/trade/store"
	"zssn/domains/users"
	usrStore "zssn/domains/users/store"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	userService      users.IUserService
	inventoryService inventory.IInventoryService
	storage          store.ITradeStorage
	tradeService     ITradeService
)

type testUser struct {
	user        *entities.User
	inventories []*entities.Inventory
}

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		os.Exit(code)
	}()

	us, err := mocks.NewUserMock()
	if err != nil {
		panic(err)
	}

	storage = mocks.NewStoreMock()
	userService = us
	inventoryService = mocks.NewInventoryMock()
	tradeService = New(storage, userService, inventoryService)

	code = m.Run()
}

func TestVerifyTransaction_EqualItem(t *testing.T) {
	fp := uuid.NewString()
	op := uuid.NewString()

	res := equalTrade(t, fp, op)
	err := tradeService.IsTransactionAmountEqual(res[fp], res[op])
	require.NoError(t, err)
}

func TestVerifyTransaction_NotEqualItem(t *testing.T) {
	fp := uuid.NewString()
	op := uuid.NewString()

	fTrade := newTradeItems(t, fp)
	oTrade := newTradeItems(t, op)

	err := tradeService.IsTransactionAmountEqual(fTrade, oTrade)
	require.EqualError(t, err, "value of the trade doesn't match")
}

func TestAnyInfectedParticipant(t *testing.T) {
	infectedUser := newUser(t)
	infectedUser.Infected = true
	table := []struct {
		name      string
		users     []*entities.User
		expectErr bool
		errMsg    string
	}{
		{
			name:      "empty users",
			users:     []*entities.User{},
			expectErr: true,
			errMsg:    "invalid users provided",
		}, {
			name:      "uninfected users",
			users:     []*entities.User{newUser(t), newUser(t)},
			expectErr: false,
		},
		{
			name:      "one infected user",
			users:     []*entities.User{newUser(t), infectedUser},
			expectErr: true,
			errMsg:    fmt.Sprintf("participant %s is infected, cannot proceed with transaction", infectedUser.Name),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := tradeService.AnyParticipantInfected(tt.users...)
			if tt.expectErr {
				require.EqualError(t, got, tt.errMsg)
			} else {
				require.NoError(t, got)
			}
		})
	}
}

func TestEnoughStock(t *testing.T) {
	table := []struct {
		name      string
		stock     entities.Stock
		tradeItem *entities.TradeItems
		expErr    bool
		errMsg    string
	}{
		{
			name: "enough_stock",
			stock: entities.Stock{
				core.ItemAmmunition: &entities.Inventory{
					Item:     core.ItemAmmunition,
					Quantity: 300,
				},
			},
			tradeItem: &entities.TradeItems{
				Items: []entities.TradeItem{
					{
						Item:     core.ItemAmmunition,
						Quantity: 299,
					},
				},
			},
			expErr: false,
		},
		{
			name:  "empty_stock",
			stock: entities.Stock{},
			tradeItem: &entities.TradeItems{
				Items: []entities.TradeItem{
					{
						Item:     core.ItemAmmunition,
						Quantity: 299,
					},
				},
			},
			expErr: true,
			errMsg: "invalid stock provided",
		},
		{
			name: "empty_items",
			stock: entities.Stock{
				core.ItemAmmunition: &entities.Inventory{
					Item:     core.ItemAmmunition,
					Quantity: 300,
				},
			},
			tradeItem: &entities.TradeItems{
				Items: []entities.TradeItem{},
			},
			expErr: true,
			errMsg: "invalid items in trade items",
		},
		{
			name: "non-existing-item",
			stock: entities.Stock{
				core.ItemAmmunition: &entities.Inventory{
					Item:     core.ItemAmmunition,
					Quantity: 300,
				},
			},
			tradeItem: &entities.TradeItems{
				Items: []entities.TradeItem{
					{
						Item:     core.ItemFood,
						Quantity: 299,
					},
				},
			},
			expErr: true,
			errMsg: "user doesn't have the item in stock Food",
		},
		{
			name: "invalid-quantity",
			stock: entities.Stock{
				core.ItemAmmunition: &entities.Inventory{
					Item:     core.ItemAmmunition,
					Quantity: 300,
				},
			},
			tradeItem: &entities.TradeItems{
				Items: []entities.TradeItem{
					{
						Item:     core.ItemAmmunition,
						Quantity: 399,
					},
				},
			},
			expErr: true,
			errMsg: "user doesn't have enough to fulfill transaction",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := tradeService.EnoughStock(tt.stock, tt.tradeItem)
			if tt.expErr {
				require.EqualError(t, got, tt.errMsg)
			} else {
				require.NoError(t, got)
			}
		})
	}
}

func TestVerifyTransaction(t *testing.T) {
	ctx := context.Background()
	fUser := setupUser(t)
	sUser := setupUser(t)

	fut := &entities.TradeItems{
		UserID: fUser.user.ID,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemWater,
				Quantity: 1,
			}, {
				Item:     core.ItemMedication,
				Quantity: 1,
			},
		},
	}

	sut := &entities.TradeItems{
		UserID: sUser.user.ID,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemAmmunition,
				Quantity: 6,
			},
		},
	}

	balances, err := inventoryService.FindMultipleInventory(ctx, fut.UserID, sut.UserID)
	require.NoError(t, err)

	err = tradeService.VerifyTransaction(ctx, balances, fut, sut)
	require.NoError(t, err)
}

func TestExecuteTransaction(t *testing.T) {
	ctx := context.Background()
	fUser := setupUser(t)
	sUser := setupUser(t)

	fut := &entities.TradeItems{
		UserID: fUser.user.ID,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemWater,
				Quantity: 1,
			}, {
				Item:     core.ItemMedication,
				Quantity: 1,
			},
		},
	}

	sut := &entities.TradeItems{
		UserID: sUser.user.ID,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemAmmunition,
				Quantity: 6,
			},
		},
	}

	err := tradeService.Execute(ctx, fut, sut)
	require.NoError(t, err)
	require.NotEmpty(t, sut.Reference)
}

func TestFailedExecutionDuetoFailedVerificcation(t *testing.T) {
	ctx := context.Background()
	fUser := setupUser(t)
	sUser := setupUser(t)

	fut := &entities.TradeItems{
		UserID: fUser.user.ID,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemWater,
				Quantity: 1,
			}, {
				Item:     core.ItemMedication,
				Quantity: 1,
			},
		},
	}

	sut := &entities.TradeItems{
		UserID: sUser.user.ID,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemAmmunition,
				Quantity: 7,
			},
		},
	}

	err := tradeService.Execute(ctx, fut, sut)
	require.EqualError(t, err, "value of the trade doesn't match")
}

func TestExecutionError(t *testing.T) {
	mockStore := &mocks.MockTradeStore{}
	mockStore.ExecuteFunc = func(ctx context.Context, seller, buyer *store.TradeItems) error {
		return fmt.Errorf("cannot complete transaction")
	}
	ts := New(mockStore, userService, inventoryService)

	ctx := context.Background()
	fUser := setupUser(t)
	sUser := setupUser(t)

	fut := &entities.TradeItems{
		UserID: fUser.user.ID,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemWater,
				Quantity: 1,
			}, {
				Item:     core.ItemMedication,
				Quantity: 1,
			},
		},
	}

	sut := &entities.TradeItems{
		UserID: sUser.user.ID,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemAmmunition,
				Quantity: 6,
			},
		},
	}

	err := ts.Execute(ctx, fut, sut)
	require.EqualError(t, err, "cannot complete transaction")
}

func TestExecutionFailed_CannotFindUsers(t *testing.T) {
	usrStore := users.MockUserStorage{
		FindUsersFunc: func(ctx context.Context, ids ...string) (map[string]*usrStore.User, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	usrSvc, err := users.New(&usrStore)
	require.NoError(t, err)
	require.NotNil(t, usrSvc)

	ts := New(storage, usrSvc, inventoryService)

	ctx := context.Background()
	fUser := setupUser(t)
	sUser := setupUser(t)

	fut := &entities.TradeItems{
		UserID: fUser.user.ID,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemWater,
				Quantity: 1,
			}, {
				Item:     core.ItemMedication,
				Quantity: 1,
			},
		},
	}

	sut := &entities.TradeItems{
		UserID: sUser.user.ID,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemAmmunition,
				Quantity: 6,
			},
		},
	}

	err = ts.Execute(ctx, fut, sut)
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func TestExecutionFailed_CannotFindInventory(t *testing.T) {
	invStore := inventory.MockInventoryStore{
		FindUsersInventoryFunc: func(ctx context.Context, userIDs ...string) (map[string]invStore.Response, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	invSvc := inventory.New(&invStore)
	require.NotNil(t, invSvc)

	ts := New(storage, userService, invSvc)

	ctx := context.Background()
	fUser := setupUser(t)
	sUser := setupUser(t)

	fut := &entities.TradeItems{
		UserID: fUser.user.ID,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemWater,
				Quantity: 1,
			}, {
				Item:     core.ItemMedication,
				Quantity: 1,
			},
		},
	}

	sut := &entities.TradeItems{
		UserID: sUser.user.ID,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemAmmunition,
				Quantity: 6,
			},
		},
	}

	err := ts.Execute(ctx, fut, sut)
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func TestGetUserTransactionHistory(t *testing.T) {
	ctx := context.Background()
	fUser := setupUser(t)
	sUser := setupUser(t)

	fut := &entities.TradeItems{
		UserID: fUser.user.ID,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemWater,
				Quantity: 1,
			}, {
				Item:     core.ItemMedication,
				Quantity: 1,
			},
		},
	}

	sut := &entities.TradeItems{
		UserID: sUser.user.ID,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemAmmunition,
				Quantity: 6,
			},
		},
	}

	err := tradeService.Execute(ctx, fut, sut)
	require.NoError(t, err)
	require.NotEmpty(t, sut.Reference)

	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	res, err := tradeService.History(ctx, sut.UserID, start, end)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, res, 3)
}

func TestGetHistoryError(t *testing.T) {
	mockStore := &mocks.MockTradeStore{}
	mockStore.HistoryFunc = func(ctx context.Context, userID string, start, endDate time.Time) ([]*store.Transaction, error) {
		return nil, gorm.ErrRecordNotFound
	}
	ts := New(mockStore, userService, inventoryService)

	ctx := context.Background()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()
	id := uuid.NewString()

	res, err := ts.History(ctx, id, start, end)
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Nil(t, res)
}

func setupUser(t *testing.T) *testUser {
	ctx := context.Background()
	t.Helper()
	user := newUser(t)
	// save user
	require.NoError(t, userService.Create(ctx, user))
	inventory := newInventory(t, user.ID)
	require.NoError(t, inventoryService.Create(ctx, inventory))
	return &testUser{
		user:        user,
		inventories: inventory,
	}
}

func newUser(t *testing.T) *entities.User {
	t.Helper()
	return &entities.User{
		Email:     gofakeit.Email(),
		Name:      gofakeit.FirstName() + " " + gofakeit.LastName(),
		Age:       20,
		Gender:    "Male",
		Latitude:  gofakeit.Latitude(),
		Longitude: gofakeit.Longitude(),
	}
}

func newInventory(t *testing.T, userID string) []*entities.Inventory {
	t.Helper()
	return []*entities.Inventory{
		{
			UserID:   userID,
			Item:     core.ItemWater,
			Quantity: uint32(gofakeit.Number(10, 1000)),
		},
		{
			UserID:   userID,
			Item:     core.ItemFood,
			Quantity: uint32(gofakeit.Number(10, 1000)),
		},
		{
			UserID:   userID,
			Item:     core.ItemMedication,
			Quantity: uint32(gofakeit.Number(10, 1000)),
		},
		{
			UserID:   userID,
			Item:     core.ItemAmmunition,
			Quantity: uint32(gofakeit.Number(10, 1000)),
		},
	}
}

func equalTrade(t *testing.T, seller, buyer string) map[string]*entities.TradeItems {
	t.Helper()
	result := make(map[string]*entities.TradeItems)
	b := &entities.TradeItems{
		UserID: seller,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemWater,
				Quantity: 1,
			},
			{
				Item:     core.ItemMedication,
				Quantity: 1,
			},
		},
	}

	s := &entities.TradeItems{
		UserID: buyer,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemAmmunition,
				Quantity: 6,
			},
		},
	}

	result[seller] = s
	result[buyer] = b

	return result
}

func newTradeItems(t *testing.T, userID string) *entities.TradeItems {
	t.Helper()
	return &entities.TradeItems{
		UserID: userID,
		Items: []entities.TradeItem{
			{
				Item:     core.ItemWater,
				Quantity: uint32(gofakeit.Number(10, 50)),
			},
			{
				Item:     core.ItemAmmunition,
				Quantity: uint32(gofakeit.Number(10, 50)),
			},
			{
				Item:     core.ItemMedication,
				Quantity: uint32(gofakeit.Number(10, 50)),
			},
		},
	}
}
