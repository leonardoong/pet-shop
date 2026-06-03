package order

import ordermodel "petshop/internal/model/order"

type CreateRequest struct {
	AddressID     string `json:"address_id"     binding:"required,uuid"`
	CustomerEmail string `json:"customer_email" binding:"omitempty,email"`
	Notes         string `json:"notes"          binding:"omitempty,max=500"`
}

type CheckoutResponse struct {
	Order        ordermodel.Order `json:"order"`
	PaymentURL   string           `json:"payment_url"`
	PaymentToken string           `json:"payment_token"`
}

type ListFilter struct {
	Status string
	Page   int
	Limit  int
}
