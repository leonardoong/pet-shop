package order

import ordermodel "petshop/internal/model/order"

// OrdersPage is a concrete paginated type used only for Swagger docs.
type OrdersPage struct {
	Items      []ordermodel.Order `json:"items"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}
