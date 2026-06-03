package payment

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MockDriver struct{}

func NewMockDriver() *MockDriver { return &MockDriver{} }

func (d *MockDriver) Name() string { return "mock" }

func (d *MockDriver) CreateTransaction(req CreateRequest) (*Response, error) {
	txID := fmt.Sprintf("MOCK-%s-%d", req.OrderID[:8], time.Now().UnixMilli())
	return &Response{
		TransactionID: txID,
		PaymentURL:    fmt.Sprintf("http://localhost:8080/api/v1/payment/mock/pay?order=%s&tx=%s", req.OrderID, txID),
		PaymentToken:  uuid.NewString(),
		RedirectURL:   fmt.Sprintf("http://localhost:3000/pesanan/%s?paid=mock", req.OrderID),
		Status:        "pending",
	}, nil
}
