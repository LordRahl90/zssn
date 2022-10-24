package users

import (
	"context"
	"fmt"
	"testing"
	"zssn/domains/users/store"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMockStorageService(t *testing.T) {
	ctx := context.Background()
	mockSvc := NewMockStore()
	u := newStoreUser(t)

	storage = mockSvc
	err := storage.Create(ctx, u)
	require.NoError(t, err)
	assert.NotEmpty(t, u.ID)
}

func TestMockStorageCreateError(t *testing.T) {
	ctx := context.Background()
	mockSvc := NewMockStore()
	u := newStoreUser(t)
	mockSvc.CreateFunc = func(ctx context.Context, user *store.User) error {
		u.ID = ""
		return fmt.Errorf("cannot connect to db right now")
	}
	err := mockSvc.Create(ctx, u)
	require.EqualError(t, err, "cannot connect to db right now")
	assert.Empty(t, u.ID)
}

func TestMockCreateFunctionNotDefined(t *testing.T) {
	ctx := context.Background()

	u := newStoreUser(t)
	mockSVC := &MockUserStorage{}

	err := mockSVC.Create(ctx, u)
	require.EqualError(t, err, "mock function not defined")
}

func TestMockMethodNotDefined(t *testing.T) {
	failureMockSvc := &MockUserStorage{}
	failureMockSvc.CreateFunc = func(ctx context.Context, user *store.User) error {
		user.ID = uuid.NewString()
		return nil
	}
	storage := failureMockSvc
	ctx := context.Background()

	user := newStoreUser(t)
	require.NoError(t, storage.Create(ctx, user))
	assert.NotEmpty(t, user.ID)

	infectedUser := newStoreUser(t)
	require.NoError(t, storage.Create(ctx, infectedUser))
	assert.NotEmpty(t, infectedUser.ID)

	err := failureMockSvc.FlagUser(ctx, user.ID, infectedUser.ID)
	require.EqualError(t, err, errMockNotDefined.Error())
}

func TestMockNotDefinedFindMethod(t *testing.T) {
	failureMockSvc := &MockUserStorage{}
	failureMockSvc.CreateFunc = func(ctx context.Context, user *store.User) error {
		user.ID = uuid.NewString()
		return nil
	}
	storage := failureMockSvc
	ctx := context.Background()

	infectedUser := newStoreUser(t)
	require.NoError(t, storage.Create(ctx, infectedUser))
	assert.NotEmpty(t, infectedUser.ID)

	res, err := failureMockSvc.Find(ctx, infectedUser.ID)
	require.EqualError(t, err, errMockNotDefined.Error())
	assert.Empty(t, res)
}

func TestMockCannotFindRecordByID(t *testing.T) {
	mockSVC := NewMockStore()
	storage := mockSVC
	ctx := context.Background()

	user := newStoreUser(t)
	require.NoError(t, storage.Create(ctx, user))
	assert.NotEmpty(t, user.ID)

	res, err := storage.Find(ctx, uuid.NewString())
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	assert.Empty(t, res)
}

func TestMockCannotFindRecordByEmail(t *testing.T) {
	mockSVC := NewMockStore()
	storage := mockSVC
	ctx := context.Background()

	user := newStoreUser(t)
	require.NoError(t, storage.Create(ctx, user))
	assert.NotEmpty(t, user.ID)

	res, err := storage.FindByEmail(ctx, gofakeit.Email())
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	assert.Empty(t, res)
}

func TestMockFlagUser(t *testing.T) {
	mockSvc := NewMockStore()
	storage := mockSvc
	ctx := context.Background()
	user := newStoreUser(t)
	require.NoError(t, storage.Create(ctx, user))
	assert.NotEmpty(t, user.ID)
	infectedUser := newStoreUser(t)
	require.NoError(t, storage.Create(ctx, infectedUser))
	assert.NotEmpty(t, infectedUser.ID)

	err := storage.FlagUser(ctx, user.ID, infectedUser.ID)
	assert.NoError(t, err)

	res, err := storage.Find(ctx, infectedUser.ID)
	require.NoError(t, err)
	assert.NotNil(t, res)

	err = storage.FlagUser(ctx, user.ID, uuid.NewString())
	assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func TestMockNotDefinedUpdateInfectedStatus(t *testing.T) {
	mockSvc := NewMockStore()
	failureMock := &MockUserStorage{}
	storage := mockSvc
	ctx := context.Background()
	infectedUser := newStoreUser(t)
	require.NoError(t, storage.Create(ctx, infectedUser))
	assert.NotEmpty(t, infectedUser.ID)

	err := failureMock.UpdateInfectedStatus(ctx, infectedUser.ID)
	require.EqualError(t, err, errMockNotDefined.Error())
}

func TestMockUpdateInfectedStatus(t *testing.T) {
	mockSvc := NewMockStore()
	storage := mockSvc
	ctx := context.Background()
	infectedUser := newStoreUser(t)
	require.NoError(t, storage.Create(ctx, infectedUser))
	assert.NotEmpty(t, infectedUser.ID)

	err := storage.UpdateInfectedStatus(ctx, infectedUser.ID)
	require.NoError(t, err)

	res, err := storage.Find(ctx, infectedUser.ID)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.True(t, res.Infected)

	err = storage.UpdateInfectedStatus(ctx, uuid.NewString())
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func TestMockNotDefinedUpdateLocation(t *testing.T) {
	failureMock := &MockUserStorage{}
	mockSvc := NewMockStore()
	storage := mockSvc
	ctx := context.Background()
	user := newStoreUser(t)
	require.NoError(t, storage.Create(ctx, user))
	assert.NotEmpty(t, user.ID)

	err := failureMock.UpdateLocation(ctx, user.ID, gofakeit.Latitude(), gofakeit.Longitude())
	require.EqualError(t, err, errMockNotDefined.Error())
}

func TestMockUpdateLocation(t *testing.T) {
	mockSvc := NewMockStore()
	storage := mockSvc
	ctx := context.Background()
	user := newStoreUser(t)
	require.NoError(t, storage.Create(ctx, user))
	assert.NotEmpty(t, user.ID)

	newLat := gofakeit.Latitude()
	newLong := gofakeit.Longitude()

	err := storage.UpdateLocation(ctx, user.ID, newLat, newLong)
	require.NoError(t, err)

	res, err := storage.Find(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, newLat, res.Latitude)
	assert.Equal(t, newLong, res.Longitude)

	err = storage.UpdateLocation(ctx, uuid.NewString(), newLat, newLong)
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func newStoreUser(t *testing.T) *store.User {
	t.Helper()
	return &store.User{
		Email:     gofakeit.Email(),
		Name:      gofakeit.LastName() + " " + gofakeit.FirstName(),
		Age:       20,
		Gender:    store.GenderMale,
		Latitude:  gofakeit.Latitude(),
		Longitude: gofakeit.Longitude(),
	}
}
