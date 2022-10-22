package core

import "strings"

// Item special type created for items
type Item int

const (
	// ItemUnknown unknown item
	ItemUnknown Item = iota
	// ItemWater water item
	ItemWater
	// ItemFood food Item
	ItemFood
	// ItemMedication medication item
	ItemMedication
	// ItemAmmunition ammunition item
	ItemAmmunition
)

var (
	// ItemPoints keeps track of the points awarded for each item
	ItemPoints = map[Item]uint32{
		ItemWater:      4,
		ItemFood:       3,
		ItemMedication: 2,
		ItemAmmunition: 1,
	}
)

// String returns the stringified version of the item
func (i Item) String() string {
	switch i {
	case ItemWater:
		return "Water"
	case ItemFood:
		return "Food"
	case ItemMedication:
		return "Medication"
	case ItemAmmunition:
		return "Ammunition"
	default:
		return "unknown"
	}
}

// ItemFromString returns an Item from string representation
func ItemFromString(item string) Item {
	switch strings.ToLower(item) {
	case strings.ToLower(ItemWater.String()):
		return ItemWater
	case strings.ToLower(ItemFood.String()):
		return ItemFood
	case strings.ToLower(ItemMedication.String()):
		return ItemMedication
	case strings.ToLower(ItemAmmunition.String()):
		return ItemAmmunition
	default:
		return ItemUnknown
	}
}
