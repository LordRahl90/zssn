package entities

import (
	"testing"

	"zssn/domains/core"

	"github.com/stretchr/testify/assert"
)

func TestCalculate(t *testing.T) {
	items := TradeItems{
		Items: []TradeItem{
			{
				Item:     core.ItemWater,
				Quantity: 10,
			},
			{
				Item:     core.ItemAmmunition,
				Quantity: 20,
			},
		},
	}

	res := items.Calculate()
	assert.Equal(t, uint32(60), res)
}
