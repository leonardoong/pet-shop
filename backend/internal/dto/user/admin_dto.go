package user

import "time"

type AdminCustomerFilter struct {
	Search string
	Page   int
	Limit  int
}

type CustomerResponse struct {
	ID         string    `json:"id"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	IsActive   bool      `json:"is_active"`
	OrderCount int64     `json:"order_count"`
	TotalSpent float64   `json:"total_spent"`
	CreatedAt  time.Time `json:"created_at"`
}

type ToggleActiveRequest struct {
	IsActive bool `json:"is_active"`
}
