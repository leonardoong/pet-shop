package user

type AddressCreateRequest struct {
	Label         string `json:"label"          binding:"omitempty,max=50"`
	RecipientName string `json:"recipient_name" binding:"required,max=100"`
	Phone         string `json:"phone"          binding:"required,indonesian_phone"`
	Street        string `json:"street"         binding:"required"`
	City          string `json:"city"           binding:"required,max=100"`
	Province      string `json:"province"       binding:"required,max=100"`
	PostalCode    string `json:"postal_code"    binding:"required,max=10"`
	IsDefault     bool   `json:"is_default"`
}

type AddressUpdateRequest struct {
	Label         string `json:"label"          binding:"omitempty,max=50"`
	RecipientName string `json:"recipient_name" binding:"omitempty,max=100"`
	Phone         string `json:"phone"          binding:"omitempty,indonesian_phone"`
	Street        string `json:"street"         binding:"omitempty"`
	City          string `json:"city"           binding:"omitempty,max=100"`
	Province      string `json:"province"       binding:"omitempty,max=100"`
	PostalCode    string `json:"postal_code"    binding:"omitempty,max=10"`
	IsDefault     *bool  `json:"is_default"`
}
