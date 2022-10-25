package responses

// Inventory contains the inventory request
type Inventory struct {
	Item     string `json:"item"`
	Quantity uint32 `json:"quantity"`
	Balance  uint32 `json:"balance"`
}
