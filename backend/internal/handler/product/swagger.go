package product

import productmodel "petshop/internal/model/product"

// ProductsPage is a concrete paginated type used only for Swagger docs.
type ProductsPage struct {
	Items      []productmodel.Product `json:"items"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
}
