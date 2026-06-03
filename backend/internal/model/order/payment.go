package order

type PaymentDriver string

const (
	PaymentDriverNone      PaymentDriver = "none"
	PaymentDriverMock      PaymentDriver = "mock"
	PaymentDriverMidtrans  PaymentDriver = "midtrans"
)

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusSettlement PaymentStatus = "settlement"
	PaymentStatusExpired    PaymentStatus = "expired"
	PaymentStatusCancel     PaymentStatus = "cancel"
	PaymentStatusDeny       PaymentStatus = "deny"
	PaymentStatusRefund     PaymentStatus = "refund"
)

func (s PaymentStatus) IsSuccess() bool {
	return s == PaymentStatusSettlement
}
