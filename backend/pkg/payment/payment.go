package payment

type ItemDetail struct {
	ID       string
	Name     string
	Price    float64
	Quantity int
}

type CreateRequest struct {
	OrderID     string
	GrossAmount float64
	CustomerName string
	CustomerEmail string
	CustomerPhone string
	Items        []ItemDetail
}

type Response struct {
	TransactionID string `json:"transaction_id"`
	PaymentURL    string `json:"payment_url"`
	PaymentToken  string `json:"payment_token"`
	RedirectURL   string `json:"redirect_url"`
	Status        string `json:"status"`
}

type Driver interface {
	Name() string
	CreateTransaction(req CreateRequest) (*Response, error)
}
