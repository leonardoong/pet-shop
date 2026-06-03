package product

type AdjustStockRequest struct {
	Operation string  `json:"operation"  binding:"required,oneof=add subtract set"`
	Quantity  int     `json:"quantity"   binding:"required,min=1"`
	CostPrice float64 `json:"cost_price"`
	Note      string  `json:"note"       binding:"omitempty,max=500"`
}

type InventoryItem struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	SKU        string `json:"sku"`
	Category   string `json:"category"`
	Stock      int    `json:"stock"`
	Price      float64 `json:"price"`
	IsActive   bool   `json:"is_active"`
	ImageURL   string `json:"image_url"`
}

type InventoryFilter struct {
	Search     string
	CategoryID string
	LowStock   *bool
	Page       int
	Limit      int
}
