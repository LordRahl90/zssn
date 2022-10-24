package requests

import "zssn/domains/core"

// Inventory contains the inventory request
type Inventory struct {
	Item     core.Item `json:"item"`
	Quantity uint32    `json:"quantity"`
}
