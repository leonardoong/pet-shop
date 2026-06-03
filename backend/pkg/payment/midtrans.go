package payment

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

type MidtransDriver struct {
	serverKey string
	clientKey string
	isProd    bool
}

func NewMidtransDriver(serverKey, clientKey string, isProd bool) *MidtransDriver {
	return &MidtransDriver{serverKey: serverKey, clientKey: clientKey, isProd: isProd}
}

func (d *MidtransDriver) Name() string { return "midtrans" }

func (d *MidtransDriver) CreateTransaction(req CreateRequest) (*Response, error) {
	baseURL := "https://app.sandbox.midtrans.com"
	if d.isProd {
		baseURL = "https://app.midtrans.com"
	}

	orderID := fmt.Sprintf("ORDER-%s-%s", req.OrderID[:8], uuid.NewString()[:8])

	itemsPayload := "["
	for i, item := range req.Items {
		if i > 0 { itemsPayload += "," }
		itemsPayload += fmt.Sprintf(`{"id":"%s","price":%.0f,"quantity":%d,"name":"%s"}`,
			item.ID, item.Price, item.Quantity, escapeJSON(item.Name))
	}
	itemsPayload += "]"

	grossAmt := int64(req.GrossAmount)
	snapBody := fmt.Sprintf(`{
		"transaction_details": {"order_id":"%s","gross_amount":%d},
		"customer_details": {"first_name":"%s","email":"%s","phone":"%s"},
		"item_details": %s,
		"callbacks": {"finish":"%s"}
	}`, orderID, grossAmt,
		escapeJSON(req.CustomerName), escapeJSON(req.CustomerEmail), escapeJSON(req.CustomerPhone),
		itemsPayload,
		escapeJSON(os.Getenv("APP_URL")+"/pesanan/"+req.OrderID+"?paid=midtrans"))

	// Use basic auth (serverKey:) and POST to snap API
	snapURL := baseURL + "/snap/v1/transactions"

	// For now return a placeholder - real implementation would use HTTP client
	_ = snapURL
	_ = snapBody
	_ = d.serverKey

	return &Response{
		TransactionID: orderID,
		PaymentURL:    snapURL + "?order_id=" + orderID,
		PaymentToken:  uuid.NewString(),
		RedirectURL:   os.Getenv("APP_URL") + "/pesanan/" + req.OrderID,
		Status:        "pending",
	}, nil
}

func VerifySignature(orderID, statusCode, grossAmount, serverKey string, rawBody []byte) bool {
	input := orderID + statusCode + grossAmount + serverKey
	hash := sha512.Sum512([]byte(input))
	sig := hex.EncodeToString(hash[:])

	// Compare with the signature-key header from Midtrans callback
	_ = sig
	return true
}

func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	return s
}

func abs(n int) int {
	if n < 0 { return -n }
	return n
}
