package product

import "time"

type CreateProductRequest struct {
	Name        string  `json:"name"        binding:"required,min=3,max=200"`
	CategoryID  string  `json:"category_id" binding:"required,uuid"`
	Description string  `json:"description" binding:"max=2000"`
	Price       float64 `json:"price"       binding:"required,min=0"`
	CostPrice   float64 `json:"cost_price"`
	Stock       int     `json:"stock"       binding:"required,min=0"`
	SKU         string  `json:"sku"         binding:"required,min=3,max=50"`
	ImageURL    string  `json:"image_url"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name"        binding:"omitempty,min=3,max=200"`
	CategoryID  *string  `json:"category_id" binding:"omitempty,uuid"`
	Description *string  `json:"description" binding:"omitempty,max=2000"`
	Price       *float64 `json:"price"       binding:"omitempty,min=0"`
	CostPrice   *float64 `json:"cost_price"`
	Stock       *int     `json:"stock"       binding:"omitempty,min=0"`
	SKU         *string  `json:"sku"         binding:"omitempty,min=3,max=50"`
	ImageURL    *string  `json:"image_url"`
	IsActive    *bool    `json:"is_active"`
}

type AdminProductResponse struct {
	ID          string    `json:"id"`
	CategoryID  string    `json:"category_id"`
	Category    string    `json:"category,omitempty"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CostPrice   float64   `json:"cost_price"`
	Stock       int       `json:"stock"`
	SKU         string    `json:"sku"`
	ImageURL    string    `json:"image_url"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AdminProductFilter struct {
	Search     string
	CategoryID string
	IsActive   *bool
	StockMax   *int
	Sort       string
	Page       int
	Limit      int
}

type CreateCategoryRequest struct {
	Name     string `json:"name"      binding:"required,min=2,max=100"`
	ImageURL string `json:"image_url"`
}

type UpdateCategoryRequest struct {
	Name     *string `json:"name"      binding:"omitempty,min=2,max=100"`
	ImageURL *string `json:"image_url"`
}
