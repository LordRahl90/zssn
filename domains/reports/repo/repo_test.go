package repo

import (
	"context"
	"os"
	"testing"
	"zssn/domains/core"
	invStore "zssn/domains/inventory/store"
	usrStore "zssn/domains/users/store"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db          *gorm.DB
	repo        IReportRepository
	userStorage usrStore.IUserStorage
	invStorage  invStore.IInventoryStorage
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
	us, err := usrStore.New(db)
	if err != nil {
		panic(err)
	}
	userStorage = us
	is, err := invStore.New(db)
	if err != nil {
		panic(err)
	}
	invStorage = is

	repo = New(db)
	code = m.Run()
}

func TestGetTotalUsers(t *testing.T) {
	ids := createSomeInfectedUser(t, 10, 2)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id IN ?", ids)
	})
	ctx := context.Background()

	total, err := repo.Total(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint32(10), total)
}

func TestGetEmptyTotal(t *testing.T) {
	ctx := context.Background()
	total, err := repo.Total(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint32(0), total)
}

func TestGetInfectedRate(t *testing.T) {
	ids := createSomeInfectedUser(t, 10, 2)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id IN ?", ids)
	})
	ctx := context.Background()

	inf, err := repo.Infected(ctx)
	require.NoError(t, err)
	require.NotNil(t, inf)
	assert.Equal(t, uint32(10), inf.Total)
	assert.Equal(t, uint32(5), inf.Infected)
	assert.Equal(t, float64(50.0), inf.Percentage)
}

func TestGetInfectedRateNoUsers(t *testing.T) {
	ctx := context.Background()

	inf, err := repo.Infected(ctx)
	require.NoError(t, err)
	require.NotNil(t, inf)
	assert.Equal(t, uint32(0), inf.Total)
	assert.Equal(t, uint32(0), inf.Infected)
	assert.Equal(t, float64(0.0), inf.Percentage)
}

func TestGetSurvivorRate(t *testing.T) {
	ids := createSomeInfectedUser(t, 10, 2)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id IN ?", ids)
	})
	ctx := context.Background()

	inf, err := repo.Survivors(ctx)
	require.NoError(t, err)
	require.NotNil(t, inf)
	assert.Equal(t, uint32(10), inf.Total)
	assert.Equal(t, uint32(5), inf.Clean)
	assert.Equal(t, float64(50.0), inf.Percentage)
}

func TestGetSurvivorRateNoUsers(t *testing.T) {
	ctx := context.Background()

	inf, err := repo.Survivors(ctx)
	require.NoError(t, err)
	require.NotNil(t, inf)
	assert.Equal(t, uint32(0), inf.Total)
	assert.Equal(t, uint32(0), inf.Clean)
	assert.Equal(t, float64(0.0), inf.Percentage)
}

func TestResources(t *testing.T) {
	ids, userIDs := createInaccessibleInventory(t, 10, 2)
	t.Cleanup(func() {
		db.Exec("DELETE FROM inventories WHERE id IN ?", ids)
		db.Exec("DELETE FROM users WHERE id IN ?", userIDs)
	})

	ctx := context.Background()
	res, err := repo.Resources(ctx)
	require.NoError(t, err)
	require.NotNil(t, res)

	water, ok := res[core.ItemWater]
	assert.True(t, ok)
	assert.Equal(t, uint32(100), water.Balance)

	medic, ok := res[core.ItemMedication]
	assert.True(t, ok)
	assert.Equal(t, uint32(150), medic.Balance)
}

func TestPoints(t *testing.T) {
	ids, userIDs := createInaccessibleInventory(t, 10, 2)
	t.Cleanup(func() {
		db.Exec("DELETE FROM inventories WHERE id IN ?", ids)
		db.Exec("DELETE FROM users WHERE id IN ?", userIDs)
	})

	ctx := context.Background()
	res, err := repo.Points(ctx)
	require.NoError(t, err)
	require.NotNil(t, res)

	water, ok := res[core.ItemWater]
	assert.True(t, ok)
	assert.Equal(t, uint32(100), water.Balance)

	medic, ok := res[core.ItemMedication]
	assert.True(t, ok)
	assert.Equal(t, uint32(150), medic.Balance)
}

func fullName() string {
	return gofakeit.FirstName() + " " + gofakeit.LastName()
}

func createSomeInfectedUser(t *testing.T, size, rate int) (res []string) {
	ctx := context.Background()
	for i := 1; i <= size; i++ {
		u := newUser(t)
		if i%rate == 0 {
			u.Infected = true
		}

		require.NoError(t, userStorage.Create(ctx, u))
		res = append(res, u.ID)
	}
	return
}

func createInaccessibleInventory(t *testing.T, size, rate int) (res, users []string) {
	t.Helper()
	ctx := context.Background()
	for i := 1; i <= size; i++ {
		u := newUser(t)
		require.NoError(t, userStorage.Create(ctx, u))
		users = append(users, u.ID)

		inv := newInventory(t, u.ID)
		require.NoError(t, invStorage.Create(ctx, inv))
		for k := range inv {
			res = append(res, inv[k].ID)
		}

		if i%rate == 0 {
			err := invStorage.UpdateUserInventoryAccessibility(ctx, u.ID)
			require.NoError(t, err)
		}
	}
	return
}

func newUser(t *testing.T) *usrStore.User {
	t.Helper()
	return &usrStore.User{
		Email:     gofakeit.Email(),
		Name:      fullName(),
		Age:       20,
		Gender:    usrStore.GenderMale,
		Latitude:  gofakeit.Latitude(),
		Longitude: gofakeit.Longitude(),
	}
}

func newInventory(t *testing.T, userID string) []*invStore.Inventory {
	t.Helper()
	return []*invStore.Inventory{
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
	db.Exec("DELETE FROM transactions")
	db.Exec("DELETE FROM inventories")
	db.Exec("DELETE FROM flag_monitors")
	db.Exec("DELETE FROM users")
}
