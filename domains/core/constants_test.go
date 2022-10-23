package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringification(t *testing.T) {
	table := []struct {
		name string
		item Item
		exp  string
	}{
		{
			name: "Water",
			item: ItemWater,
			exp:  "Water",
		},
		{
			name: "Food",
			item: ItemFood,
			exp:  "Food",
		},
		{
			name: "Medication",
			item: ItemMedication,
			exp:  "Medication",
		},
		{
			name: "Ammunition",
			item: ItemAmmunition,
			exp:  "Ammunition",
		},
		{
			name: "Unknown",
			item: ItemUnknown,
			exp:  "unknown",
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.item.String()
			require.Equal(t, tt.exp, got)
		})
	}
}

func TestGetItemFromString(t *testing.T) {
	table := []struct {
		name string
		item string
		exp  Item
	}{
		{
			name: "Water",
			item: "Water",
			exp:  ItemWater,
		},
		{
			name: "Ammunition",
			item: "Ammunition",
			exp:  ItemAmmunition,
		},
		{
			name: "Food",
			item: "Food",
			exp:  ItemFood,
		},
		{
			name: "Unknown",
			item: "Unknown",
			exp:  ItemUnknown,
		},
		{
			name: "Medication",
			item: "Medication",
			exp:  ItemMedication,
		},
		{
			name: "Hummer Jeep",
			item: "Hummer Jeep",
			exp:  ItemUnknown,
		},
		{
			name: "Bugatti",
			item: "Bugatti",
			exp:  ItemUnknown,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := ItemFromString(tt.item)
			require.Equal(t, tt.exp, got)
		})
	}
}
