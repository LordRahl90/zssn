package servers

import (
	"encoding/json"
	"net/http"
	"testing"
	"zssn/domains/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSurvivorReport(t *testing.T) {
	user1 := createDemoUser(t)
	user2 := createDemoUser(t)

	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id IN ?",
			[]string{user1.ID, user2.ID})
	})

	res := handleReqest(t, http.MethodGet, "/reports/survivors", user1.Token, nil)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var result *entities.Survivor

	err := json.NewDecoder(res.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, uint32(2), result.Total)
	assert.Equal(t, uint32(2), result.Clean)
	assert.Equal(t, float64(100), result.Percentage)
}
