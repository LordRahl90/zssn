package servers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"zssn/domains/core"
	"zssn/requests"
	"zssn/responses"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTrade(t *testing.T) {
	ctx := context.Background()
	user1 := createDemoUser(t)
	user2 := createDemoUser(t)

	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id IN ?", []string{user1.ID, user2.ID})
	})

	tr1 := &requests.TradeItems{
		UserID: user1.ID,
		Items: []requests.TradeItem{
			{
				Item:     core.ItemWater,
				Quantity: 1,
			},
			{
				Item:     core.ItemMedication,
				Quantity: 1,
			},
		},
	}

	tr2 := &requests.TradeItems{
		UserID: user2.ID,
		Items: []requests.TradeItem{
			{
				Item:     core.ItemAmmunition,
				Quantity: 6,
			},
		},
	}

	tr := &requests.TradeRequest{
		Owner:       tr1,
		SecondParty: tr2,
	}
	b, err := json.Marshal(tr)
	require.NoError(t, err)

	res := handleReqest(t, http.MethodPost, "/trades", user1.Token, b)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	var result *responses.Trade
	err = json.NewDecoder(res.Body).Decode(&result)
	require.NoError(t, err)

	participantsInventory, err := inventoryService.FindMultipleInventory(ctx, tr1.UserID, tr2.UserID)
	require.NoError(t, err)
	require.NotNil(t, participantsInventory)

	rst := participantsInventory[user1.ID]

	waterRec := rst[core.ItemWater]
	assert.NotEqual(t, waterRec.Balance, waterRec.Quantity)
	assert.Equal(t, waterRec.Balance+1, waterRec.Quantity)

	medicRec := rst[core.ItemMedication]
	assert.NotEqual(t, medicRec.Balance, medicRec.Quantity)
	assert.Equal(t, medicRec.Balance+1, medicRec.Quantity)

	ammoRec := rst[core.ItemAmmunition]
	assert.NotEqual(t, ammoRec.Balance, ammoRec.Quantity)
	assert.Equal(t, ammoRec.Balance-6, ammoRec.Quantity)

	rst = participantsInventory[tr2.UserID]

	ammoRec = rst[core.ItemAmmunition]
	assert.Equal(t, ammoRec.Balance+6, ammoRec.Quantity)

	waterRec = rst[core.ItemWater]
	assert.NotEqual(t, waterRec.Balance, waterRec.Quantity)
	assert.Equal(t, waterRec.Balance-1, waterRec.Quantity)

	medicRec = rst[core.ItemMedication]
	assert.NotEqual(t, medicRec.Balance, medicRec.Quantity)
	assert.Equal(t, medicRec.Balance-1, medicRec.Quantity)
}

func TestUnequalTrade(t *testing.T) {
	user1 := createDemoUser(t)
	user2 := createDemoUser(t)

	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id IN ?", []string{user1.ID, user2.ID})
	})

	tr1 := &requests.TradeItems{
		UserID: user1.ID,
		Items: []requests.TradeItem{
			{
				Item:     core.ItemWater,
				Quantity: 1,
			},
			{
				Item:     core.ItemMedication,
				Quantity: 1,
			},
		},
	}

	tr2 := &requests.TradeItems{
		UserID: user2.ID,
		Items: []requests.TradeItem{
			{
				Item:     core.ItemAmmunition,
				Quantity: 7,
			},
		},
	}

	tr := &requests.TradeRequest{
		Owner:       tr1,
		SecondParty: tr2,
	}
	b, err := json.Marshal(tr)
	require.NoError(t, err)

	res := handleReqest(t, http.MethodPost, "/trades", user1.Token, b)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestExecuteTradeInvalidJSON(t *testing.T) {
	b := []byte(`{
		"originator": {
		  "userID": "32d748f0-9864-46d4-9506-a7d5aee6e33a",
		  "reference": "",
		  "items": [
			{
			  "item": 1,
			  "quantity": 1
			},
			{
			  "item": 3,
			  "quantity": 1
			}
		  ]
		},
		"second_party": {
		  "userID": "2a53e866-7977-4763-a26f-f7724da49231",
		  "reference": "",
		  "items": [
			{
			  "item": 4,
			  "quantity": 7
			}
		  ]
		},
	  }`)

	user2 := createDemoUser(t)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users WHERE id IN ?", []string{user2.ID})
	})
	res := handleReqest(t, http.MethodPost, "/trades", user2.Token, b)
	require.Equal(t, http.StatusBadRequest, res.StatusCode)
}
