package servers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"zssn/domains/core"
	"zssn/domains/entities"
	"zssn/requests"
	"zssn/responses"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateNewUser(t *testing.T) {
	id := ""
	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id = ?", id)
	})
	u := newSurvivor(t)
	b, err := json.Marshal(u)
	require.NoError(t, err)

	res := handleReqest(t, http.MethodPost, "/users", "", b)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var result *responses.User
	err = json.NewDecoder(res.Body).Decode(&result)
	require.NoError(t, err)

	id = result.ID

	assert.Equal(t, u.Name, result.Name)
	assert.Equal(t, u.Email, result.Email)
	assert.Equal(t, u.Latitude, result.Latitude)
	assert.Equal(t, u.Longitude, result.Longitude)
}

func TestUserUpdatesLocation(t *testing.T) {
	ctx := context.Background()
	user := createDemoUser(t)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id = ?", user.ID)
	})

	req := &requests.UpdateLocation{
		Latitude:  gofakeit.Latitude(),
		Longitude: gofakeit.Longitude(),
	}
	b, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotNil(t, b)

	res := handleReqest(t, http.MethodPatch, "/users/location", user.Token, b)
	require.Equal(t, http.StatusOK, res.StatusCode)

	data, err := userService.Find(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, data)

	assert.Equal(t, data.Latitude, req.Latitude)
	assert.Equal(t, data.Longitude, req.Longitude)
}

func TestNewTokenForUser(t *testing.T) {
	// ctx := context.Background()
	user := createDemoUser(t)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id = ?", user.ID)
	})

	req := &requests.NewToken{
		Email: user.Email,
	}

	b, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotNil(t, b)

	// infected user tries to update location
	res := handleReqest(t, http.MethodPost, "/users/new-token", "", b)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var result *responses.User
	err = json.NewDecoder(res.Body).Decode(&result)
	require.NoError(t, err)

}

func TestCreateUserWithInvalidData(t *testing.T) {
	u := newSurvivor(t)
	u.Name = ""
	b, err := json.Marshal(u)
	require.NoError(t, err)

	res := handleReqest(t, http.MethodPost, "/users", "", b)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestCreateNewUserInvalidJSON(t *testing.T) {
	b := `
	{
		"email": "cordiajacobi@carroll.net",
		"name": "Bart Beatty",
		"age": 20,
		"gender": "Male",
		"latitude": -78.18533654085428,
		"longitude": -123.65306829619516,
		"inventories": [
			{
				"item": 1,
				"quantity": 645
			},
		],
	}
	`
	res := handleReqest(t, http.MethodPost, "/users", "", []byte(b))
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestGetUserDetails(t *testing.T) {
	result := createDemoUser(t)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id = ?", result.ID)
	})

	res := handleReqest(t, http.MethodGet, "/users/me", result.Token, nil)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var data *responses.User
	err = json.Unmarshal(body, &data)
	require.NoError(t, err)

	assert.Equal(t, result.ID, data.ID)
	assert.NotEmpty(t, data.Token)
	assert.Equal(t, result.Email, data.Email)
}

func TestGetUserThatDoesntExist(t *testing.T) {
	td := core.TokenData{
		UserID: uuid.NewString(),
		Email:  gofakeit.Email(),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	res := handleReqest(t, http.MethodGet, "/users/me", token, nil)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestFlagUserAsInfected(t *testing.T) {
	ctx := context.Background()
	user := createDemoUser(t)
	infectedUser := createDemoUser(t)
	user2 := createDemoUser(t)
	user3 := createDemoUser(t)

	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id IN ?",
			[]string{user.ID, infectedUser.ID, user2.ID, user3.ID})
	})

	flag := requests.FlagUser{
		InfectedUserID: infectedUser.ID,
	}
	b, err := json.Marshal(flag)
	require.NoError(t, err)

	res := handleReqest(t, http.MethodPost, "/users/flag", user.Token, b)
	require.Equal(t, http.StatusOK, res.StatusCode)

	res = handleReqest(t, http.MethodPost, "/users/flag", user2.Token, b)
	require.Equal(t, http.StatusOK, res.StatusCode)

	res = handleReqest(t, http.MethodPost, "/users/flag", user3.Token, b)
	require.Equal(t, http.StatusOK, res.StatusCode)

	data, err := userService.Find(ctx, infectedUser.ID)
	require.NoError(t, err)
	require.True(t, data.Infected)

	invData, err := inventoryService.FindUserInventory(ctx, flag.InfectedUserID)
	require.NoError(t, err)
	for _, v := range invData {
		require.False(t, v.Accessible)
	}

	req := &requests.UpdateLocation{
		Latitude:  gofakeit.Latitude(),
		Longitude: gofakeit.Longitude(),
	}
	b, err = json.Marshal(req)
	require.NoError(t, err)

	// infected user tries to update location
	res = handleReqest(t, http.MethodPatch, "/users/location", infectedUser.Token, b)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	// get survivor report.
	// should be 75% survivor, 25% infected
	res = handleReqest(t, http.MethodGet, "/reports/survivors", "", nil)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var result *entities.Survivor

	err = json.NewDecoder(res.Body).Decode(&result)
	require.NoError(t, err)

	require.Equal(t, float64(75.0), result.Percentage)
	require.Equal(t, uint32(3), result.Clean)
	require.Equal(t, uint32(4), result.Total)

	res = handleReqest(t, http.MethodGet, "/reports/infected", "", nil)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var infResult *entities.Infected

	err = json.NewDecoder(res.Body).Decode(&infResult)
	require.NoError(t, err)

	require.Equal(t, float64(25.0), infResult.Percentage)
	require.Equal(t, uint32(1), infResult.Infected)
	require.Equal(t, uint32(4), infResult.Total)
}

func createDemoUser(t *testing.T) responses.User {
	u := newSurvivor(t)
	b, err := json.Marshal(u)
	require.NoError(t, err)

	res := handleReqest(t, http.MethodPost, "/users", "", b)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	var result *responses.User
	err = json.NewDecoder(res.Body).Decode(&result)
	require.NoError(t, err)

	return *result
}
