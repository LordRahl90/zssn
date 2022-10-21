package servers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	server *Server
	db     *gorm.DB
)

func TestMain(m *testing.M) {
	server = New(db)
	os.Exit(m.Run())
}

func TestLandinPage(t *testing.T) {
	res := handleReqest(t, http.MethodGet, "/", "")
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "Welcome to Zombie Survival Social Network API", string(body))
}

func handleReqest(t *testing.T, method, path, body string) *http.Response {
	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, bytes.NewBuffer([]byte(body)))
	}

	res, err := server.Router.Test(req)
	require.NoError(t, err)
	return res
}
