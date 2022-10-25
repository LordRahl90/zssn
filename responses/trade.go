package responses

// Trade response struct for trade
type Trade struct {
	Reference string       `json:"reference"`
	Balance   []*Inventory `json:"balance"`
}
