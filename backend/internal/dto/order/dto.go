package order

type CreateRequest struct {
	AddressID string `json:"address_id" binding:"required,uuid"`
	Notes     string `json:"notes"      binding:"omitempty,max=500"`
}

type ListFilter struct {
	Status string
	Page   int
	Limit  int
}
