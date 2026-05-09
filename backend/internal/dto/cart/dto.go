package cart

type AddItemRequest struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
	Quantity  int    `json:"quantity"   binding:"required,min=1"`
}

type UpdateItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}
