package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockStorageService(t *testing.T) {
	ctx := context.Background()
	mockSvc := &MockUserStorage{
		CreateFunc: func(ctx context.Context, user *User) error {
			user.ID = uuid.NewString()
			return nil
		},
	}
	u := newUser(t)
	storage, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, storage)

	storage = mockSvc
	err = storage.Create(ctx, u)
	require.NoError(t, err)
	assert.NotEmpty(t, u.ID)

	mockSvc.CreateFunc = func(ctx context.Context, user *User) error {
		u.ID = ""
		return fmt.Errorf("cannot connect to db right now")
	}
	err = storage.Create(ctx, u)
	require.EqualError(t, err, "cannot connect to db right now")
	assert.Empty(t, u.ID)
}

func TestMockCreateFunctionNotDefined(t *testing.T) {
	ctx := context.Background()
	u := newUser(t)
	mockSVC := &MockUserStorage{}
	err := mockSVC.Create(ctx, u)
	require.EqualError(t, err, "mock function not defined")
}

func TestMockFlagUser(t *testing.T) {
	mockSvc := &MockUserStorage{
		CreateFunc: func(ctx context.Context, user *User) error {
			user.ID = uuid.NewString()
			return nil
		},
	}
	storage := mockSvc
	ctx := context.Background()
	user := newUser(t)
	require.NoError(t, storage.Create(ctx, user))
	assert.NotEmpty(t, user.ID)
	infectedUser := newUser(t)
	require.NoError(t, storage.Create(ctx, infectedUser))
	assert.NotEmpty(t, infectedUser.ID)

	err := storage.FlagUser(ctx, user.ID, infectedUser.ID)
	require.EqualError(t, err, errMockNotDefined.Error())

	mockSvc.FlagUserFunc = func(ctx context.Context, id, infectedUser string) error {
		return nil
	}
	res, err := storage.Find(ctx, infectedUser.ID)
	require.EqualError(t, err, errMockNotDefined.Error())
	assert.Empty(t, res)

	mockSvc.FindFunc = func(ctx context.Context, id string) (*User, error) {
		return &User{
			ID:       infectedUser.ID,
			Infected: true,
			FlagMonitor: []FlagMonitor{
				{
					InfectedUserID: id,
					UserID:         user.ID,
				},
			},
		}, nil
	}
	err = storage.FlagUser(ctx, user.ID, infectedUser.ID)
	assert.NoError(t, err)

	res, err = storage.Find(ctx, infectedUser.ID)
	require.NoError(t, err)
	assert.NotNil(t, res)
}

func TestMockUpdateInfectedStatus(t *testing.T) {
	mockSvc := &MockUserStorage{
		CreateFunc: func(ctx context.Context, user *User) error {
			user.ID = uuid.NewString()
			return nil
		},
	}
	storage := mockSvc
	ctx := context.Background()
	infectedUser := newUser(t)
	require.NoError(t, storage.Create(ctx, infectedUser))
	assert.NotEmpty(t, infectedUser.ID)

	err := storage.UpdateInfectedStatus(ctx, infectedUser.ID)
	require.EqualError(t, err, errMockNotDefined.Error())

	mockSvc.UpdateInfectedStatusFunc = func(ctx context.Context, id string) error {
		return nil
	}

	err = storage.UpdateInfectedStatus(ctx, infectedUser.ID)
	require.NoError(t, err)
}

func TestMockUpdateLocation(t *testing.T) {
	mockSvc := &MockUserStorage{
		CreateFunc: func(ctx context.Context, user *User) error {
			user.ID = uuid.NewString()
			return nil
		},
	}
	storage := mockSvc
	ctx := context.Background()
	user := newUser(t)
	require.NoError(t, storage.Create(ctx, user))
	assert.NotEmpty(t, user.ID)

	err := storage.UpdateLocation(ctx, user.ID, gofakeit.Latitude(), gofakeit.Longitude())
	require.EqualError(t, err, errMockNotDefined.Error())

	mockSvc.UpdateLocationFunc = func(ctx context.Context, id string, lat, long float64) error {
		user.Latitude = lat
		user.Longitude = long
		return nil
	}

	err = storage.UpdateLocation(ctx, user.ID, gofakeit.Latitude(), gofakeit.Longitude())
	require.NoError(t, err)
}
