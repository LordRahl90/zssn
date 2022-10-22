package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockStorageService(t *testing.T) {
	ctx := context.Background()
	u := newUser(t)
	storage, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, storage)
	mockSvc := &MockUserStorage{
		CreateFunc: func(ctx context.Context, user *User) error {
			user.ID = uuid.NewString()
			return nil
		},
	}
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
