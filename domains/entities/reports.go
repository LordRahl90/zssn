package entities

import "zssn/domains/core"

// Survivor struct for survivors rate
type Survivor struct {
	Total      uint32  `json:"total_survivors" gorm:"total"`
	Clean      uint32  `json:"clean" gorm:"clean"`
	Percentage float64 `json:"percentage_clean" gorm:"percentage_clean"`
}

// Infected struct for infected rate
type Infected struct {
	Total      uint32  `json:"total_survivors" gorm:"total"`
	Infected   uint32  `json:"infected_survivors" gorm:"infected"`
	Percentage float64 `json:"percentage_infected" gorm:"percentage_infected"`
}

// Resource report for the remaining resource per Item
type Resource struct {
	Item    core.Item `json:"item" gorm:"item"`
	Balance uint32    `json:"balance" gorm:"balance"`
}

// LostPoint report highlighting the total points per item
type LostPoint struct {
	Item    core.Item `json:"item" gorm:"item"`
	Balance uint32    `json:"balance" gorm:"balance"`
	Point   uint32    `json:"point" gorm:"point"`
}

// ResourceSharing contains the rate of resource per survivor to the nearest whole number
type ResourceSharing struct {
	Item        string `json:"item" gorm:"item"`
	Balance     uint32 `json:"balance" gorm:"balance"`
	PerSurvivor uint32 `json:"per_survivor"`
}
