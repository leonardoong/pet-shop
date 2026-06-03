package order

import "time"

type AdminOrderFilter struct {
	Status     string
	Search     string
	DateFrom   string
	DateTo     string
	Page       int
	Limit      int
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Note   string `json:"note"   binding:"omitempty,max=500"`
}

type AdminOrderResponse struct {
	ID           string              `json:"id"`
	CustomerName string              `json:"customer_name"`
	CustomerID   string              `json:"customer_id"`
	Status       string              `json:"status"`
	TotalAmount  float64             `json:"total_amount"`
	ShipName     string              `json:"ship_name"`
	ShipPhone    string              `json:"ship_phone"`
	ShipStreet   string              `json:"ship_street"`
	ShipCity     string              `json:"ship_city"`
	ShipProvince string              `json:"ship_province"`
	ShipPostal   string              `json:"ship_postal"`
	Notes        string              `json:"notes"`
	Items        []AdminOrderItemResponse `json:"items,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type AdminOrderItemResponse struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Product   string  `json:"product_name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Subtotal  float64 `json:"subtotal"`
}
