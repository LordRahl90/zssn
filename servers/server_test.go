package servers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"zssn/domains/core"
	"zssn/requests"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	server *Server
	db     *gorm.DB
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
	server = s
	code = m.Run()
}

func TestLandingPage(t *testing.T) {
	res := handleReqest(t, http.MethodGet, "/", "", nil)
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "Welcome to Zombie Survival Social Network API", string(body))
}

func TestUnHandledMethod(t *testing.T) {
	res := handleReqest(t, http.MethodPost, "/", "", nil)
	require.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
}

func TestNonExistingLink(t *testing.T) {
	res := handleReqest(t, http.MethodGet, "/home", "", nil)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestNewServer(t *testing.T) {
	s, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, s)
	require.NotNil(t, s.Router)
}

func newSurvivor(t *testing.T) *requests.Survivor {
	t.Helper()
	return &requests.Survivor{
		Name:      gofakeit.FirstName() + " " + gofakeit.LastName(),
		Email:     gofakeit.Email(),
		Age:       uint32(20),
		Gender:    "Male",
		Latitude:  gofakeit.Latitude(),
		Longitude: gofakeit.Longitude(),
		Inventory: []requests.Inventory{
			{
				Item:     core.ItemWater,
				Quantity: uint32(gofakeit.Number(10, 1000)),
			},
			{
				Item:     core.ItemFood,
				Quantity: uint32(gofakeit.Number(10, 1000)),
			},
			{
				Item:     core.ItemMedication,
				Quantity: uint32(gofakeit.Number(100, 1000)),
			},
			{
				Item:     core.ItemAmmunition,
				Quantity: uint32(gofakeit.Number(1000, 2000)),
			},
		},
	}
}

func handleReqest(t *testing.T, method, path, token string, body []byte) *http.Response {
	t.Helper()
	var req *http.Request
	if len(body) == 0 {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, bytes.NewBuffer(body))
	}

	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	res, err := server.Router.Test(req)
	require.NoError(t, err)
	return res
}

func cleanup() {
	db.Exec("DELETE FROM transactions")
	db.Exec("DELETE FROM inventories")
	db.Exec("DELETE FROM flag_monitors")
	db.Exec("DELETE FROM users")
}

func setupTestDB() (*gorm.DB, error) {
	env := os.Getenv("ENVIRONMENT")
	dsn := "root:@tcp(127.0.0.1:3306)/zssn?charset=utf8mb4&parseTime=True&loc=Local"
	if env == "cicd" {
		dsn = "zssn_user:password@tcp(127.0.0.1:33306)/zssn?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
