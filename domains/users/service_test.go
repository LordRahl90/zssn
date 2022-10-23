package users

import (
	"context"
	"os"
	"testing"
	"zssn/domains/entities"
	"zssn/domains/users/store"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	storage = NewMockStore()
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		os.Exit(code)
	}()
	code = m.Run()
	// cleanup resources
}

func TestNewService(t *testing.T) {
	svc, err := New(storage)
	require.NoError(t, err)
	require.NotNil(t, svc)
}

func TestCreateNewUser(t *testing.T) {
	ctx := context.Background()
	svc, err := New(storage)
	require.NoError(t, err)
	require.NotNil(t, svc)

	u := newUser(t)
	require.NoError(t, svc.Create(ctx, u))
	assert.NotEmpty(t, u.ID)
}

func TestCreateNewUserWithError(t *testing.T) {
	ctx := context.Background()
	failureStore := &MockUserStorage{}
	svc, err := New(failureStore)
	require.NoError(t, err)
	require.NotNil(t, svc)

	failureStore.CreateFunc = func(ctx context.Context, user *store.User) error {
		user.ID = ""
		return gorm.ErrRecordNotFound
	}

	u := newUser(t)
	require.EqualError(t, svc.Create(ctx, u), gorm.ErrRecordNotFound.Error())
	assert.Empty(t, u.ID)
}

func TestFindUserByID(t *testing.T) {
	ctx := context.Background()
	svc, err := New(storage)
	require.NoError(t, err)
	require.NotNil(t, svc)

	var (
		ids   []string
		users []*entities.User
	)

	for i := 0; i <= 3; i++ {
		user := newUser(t)
		require.NoError(t, svc.Create(ctx, user))
		ids = append(ids, user.ID)
		users = append(users, user)
	}

	res, err := svc.Find(ctx, ids[0])
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, users)
	assert.Equal(t, users[0].ID, res.ID)
}

func TestFindUserByIDFailure(t *testing.T) {
	ctx := context.Background()
	svc, err := New(storage)
	require.NoError(t, err)
	require.NotNil(t, svc)

	res, err := svc.Find(ctx, uuid.NewString())
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	assert.Empty(t, res)
}

func TestFindUserByEmail(t *testing.T) {
	ctx := context.Background()
	svc, err := New(storage)
	require.NoError(t, err)
	require.NotNil(t, svc)

	var (
		emails []string
		users  []*entities.User
	)

	for i := 0; i <= 3; i++ {
		user := newUser(t)
		require.NoError(t, svc.Create(ctx, user))
		emails = append(emails, user.Email)
		users = append(users, user)
	}

	res, err := svc.FindByEmail(ctx, emails[0])
	require.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, users)
	assert.Equal(t, users[0].ID, res.ID)
}

func TestFindUserByEmailFailure(t *testing.T) {
	ctx := context.Background()
	svc, err := New(storage)
	require.NoError(t, err)
	require.NotNil(t, svc)

	res, err := svc.FindByEmail(ctx, gofakeit.Email())
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	assert.Empty(t, res)
}

func TestUpdateLocation(t *testing.T) {
	ctx := context.Background()
	svc, err := New(storage)
	require.NoError(t, err)
	require.NotNil(t, svc)

	user := newUser(t)
	require.NoError(t, svc.Create(ctx, user))
	require.NotEmpty(t, user.ID)

	newLat := gofakeit.Latitude()
	newLong := gofakeit.Longitude()

	err = svc.UpdateLocation(ctx, user.ID, newLat, newLong)
	require.NoError(t, err)

	res, err := svc.Find(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, newLat, res.Latitude)
	assert.Equal(t, newLong, res.Longitude)
}

func TestUpdateLocationError(t *testing.T) {
	ctx := context.Background()
	svc, err := New(storage)
	require.NoError(t, err)
	require.NotNil(t, svc)

	newLat := gofakeit.Latitude()
	newLong := gofakeit.Longitude()

	user := newUser(t)

	err = svc.UpdateLocation(ctx, user.ID, newLat, newLong)
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func TestFlagUser(t *testing.T) {
	ctx := context.Background()
	svc, err := New(storage)
	require.NoError(t, err)
	require.NotNil(t, svc)

	infectedUser := newUser(t)
	require.NoError(t, svc.Create(ctx, infectedUser))
	require.NotEmpty(t, infectedUser.ID)

	flagger := newUser(t)
	require.NoError(t, svc.Create(ctx, flagger))
	require.NotEmpty(t, flagger.ID)

	err = svc.FlagUser(ctx, flagger.ID, infectedUser.ID)
	require.NoError(t, err)

	res, err := svc.Find(ctx, infectedUser.ID)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, res.FlagMonitor, 1)
}

func TestFlagNonExistingUser(t *testing.T) {
	ctx := context.Background()
	svc, err := New(storage)
	require.NoError(t, err)
	require.NotNil(t, svc)

	infectedUser := newUser(t)

	flagger := newUser(t)
	require.NoError(t, svc.Create(ctx, flagger))
	require.NotEmpty(t, flagger.ID)

	err = svc.FlagUser(ctx, flagger.ID, infectedUser.ID)
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func TestFlagUserWithNonExistingUserAccount(t *testing.T) {
	ctx := context.Background()
	svc, err := New(storage)
	require.NoError(t, err)
	require.NotNil(t, svc)

	infectedUser := newUser(t)
	require.NoError(t, svc.Create(ctx, infectedUser))
	require.NotEmpty(t, infectedUser.ID)

	flagger := newUser(t)

	err = svc.FlagUser(ctx, flagger.ID, infectedUser.ID)
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
}

func TestIsInfected(t *testing.T) {
	ctx := context.Background()
	svc, err := New(storage)
	require.NoError(t, err)
	require.NotNil(t, svc)

	infectedUser := newUser(t)
	require.NoError(t, svc.Create(ctx, infectedUser))
	require.NotEmpty(t, infectedUser.ID)

	ok, err := svc.IsInfected(ctx, infectedUser.ID)
	require.NoError(t, err)
	assert.False(t, ok)

	err = storage.UpdateInfectedStatus(ctx, infectedUser.ID)
	require.NoError(t, err)

	ok, err = svc.IsInfected(ctx, infectedUser.ID)
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestIsInfectedWithFindError(t *testing.T) {
	ctx := context.Background()
	svc, err := New(storage)
	require.NoError(t, err)
	require.NotNil(t, svc)

	infectedUser := newUser(t)
	require.NoError(t, svc.Create(ctx, infectedUser))
	require.NotEmpty(t, infectedUser.ID)

	ok, err := svc.IsInfected(ctx, uuid.NewString())
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	assert.False(t, ok)
}

func newUser(t *testing.T) *entities.User {
	t.Helper()
	return &entities.User{
		Email:     gofakeit.Email(),
		Name:      gofakeit.Name(),
		Age:       20,
		Gender:    "Male",
		Latitude:  gofakeit.Latitude(),
		Longitude: gofakeit.Longitude(),
	}
}
