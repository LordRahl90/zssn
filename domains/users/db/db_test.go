package db

import (
	"context"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db      *gorm.DB
	storage IUserStorage
)

func TestMain(m *testing.M) {
	cmd := 1
	defer func() {
		os.Exit(cmd)
	}()
	d, err := setupTestDB()
	if err != nil {
		panic(err)
	}
	db = d
	storage, err = New(db)
	if err != nil {
		panic(err)
	}
	cmd = m.Run()
	cleanup()
}

func TestNewStorageService(t *testing.T) {
	svc, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, svc)
	assert.IsType(t, &UserStorage{}, svc)

	mockSvc := MockUserStorage{
		CreateFunc: func(ctx context.Context, user *User) error {
			user.ID = uuid.NewString()
			return nil
		},
	}
	require.NotNil(t, mockSvc)
	svc = &mockSvc // they implement the same interface and should be swappable
	require.NotNil(t, svc)
}

func TestCreateNewUser(t *testing.T) {
	ctx := context.Background()
	u := &User{
		Email:     gofakeit.Email(),
		Name:      fullName(),
		Age:       20,
		Gender:    GenderMale,
		Latitude:  gofakeit.Latitude(),
		Longitude: gofakeit.Longitude(),
	}
	err := storage.Create(ctx, u)
	require.NoError(t, err)
	assert.NotEmpty(t, u.ID)
}

func TestFindUserByID(t *testing.T) {
	ctx := context.Background()
	ids := []string{}
	users := make(map[string]*User)
	for i := 1; i <= 3; i++ {
		u := newUser(t)
		require.NoError(t, storage.Create(ctx, u))
		ids = append(ids, u.ID)
		users[u.ID] = u
	}

	id := ids[2]
	res, err := storage.Find(ctx, id)
	require.NoError(t, err)
	assert.NotNil(t, res)
	mapUser, ok := users[id]
	assert.True(t, ok)
	assert.Equal(t, mapUser.Name, res.Name)
	assert.Equal(t, mapUser.Email, res.Email)
}

func TestFindByUserWithInvalidID(t *testing.T) {
	ctx := context.Background()
	storage, err := New(db)
	require.NoError(t, err)
	res, err := storage.Find(ctx, uuid.NewString())
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	assert.Empty(t, res)
}

func TestFindByUserWithInvalidEmail(t *testing.T) {
	ctx := context.Background()
	res, err := storage.Find(ctx, gofakeit.Email())
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	assert.Empty(t, res)
}

func TestFindUserByEmail(t *testing.T) {
	ctx := context.Background()
	storage, err := New(db)
	require.NoError(t, err)
	emails := []string{}
	users := make(map[string]*User)
	for i := 1; i <= 3; i++ {
		u := newUser(t)
		require.NoError(t, storage.Create(ctx, u))
		emails = append(emails, u.Email)
		users[u.Email] = u
	}

	email := emails[2]
	res, err := storage.FindByEmail(ctx, email)
	require.NoError(t, err)
	assert.NotNil(t, res)
	mapUser, ok := users[email]
	assert.True(t, ok)
	assert.Equal(t, mapUser.Name, res.Name)
	assert.Equal(t, mapUser.Email, res.Email)
}

func fullName() string {
	return gofakeit.FirstName() + " " + gofakeit.LastName()
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
	err := db.Exec("DELETE FROM users")
	if err != nil {
		panic(err)
	}
	err = db.Exec("DELETE FROM flag_monitor")
	if err != nil {
		panic(err)
	}
}

func newUser(t *testing.T) *User {
	t.Helper()
	return &User{
		Email:     gofakeit.Email(),
		Name:      fullName(),
		Age:       20,
		Gender:    GenderMale,
		Latitude:  gofakeit.Latitude(),
		Longitude: gofakeit.Longitude(),
	}
}
