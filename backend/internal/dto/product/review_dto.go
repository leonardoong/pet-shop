package product

type CreateReviewRequest struct {
	Rating  int    `json:"rating"  binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"omitempty,max=1000"`
}

type ReviewResponse struct {
	ID           string `json:"id"`
	CustomerName string `json:"customer_name"`
	Rating       int    `json:"rating"`
	Comment      string `json:"comment"`
	IsApproved   bool   `json:"is_approved"`
	CreatedAt    string `json:"created_at"`
}

type AdminReviewFilter struct {
	IsApproved *bool
	Page       int
	Limit      int
}
