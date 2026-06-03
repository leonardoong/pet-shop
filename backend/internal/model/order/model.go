package order

import (
	"time"

	productmodel "petshop/internal/model/product"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusConfirmed  Status = "confirmed"
	StatusProcessing Status = "processing"
	StatusShipped    Status = "shipped"
	StatusDelivered  Status = "delivered"
	StatusCancelled  Status = "cancelled"
)

type Order struct {
	ID                   uuid.UUID        `gorm:"type:uuid;primaryKey"               json:"id"`
	CustomerID           uuid.UUID        `gorm:"type:uuid;not null;index"           json:"customer_id"`
	Status               Status           `gorm:"type:order_status;default:'pending'" json:"status"`
	TotalAmount          float64          `gorm:"type:numeric(12,2);not null"        json:"total_amount"`
	ShipName             string           `gorm:"column:ship_name"                   json:"ship_name"`
	ShipPhone            string           `gorm:"column:ship_phone"                  json:"ship_phone"`
	ShipStreet           string           `gorm:"column:ship_street"                 json:"ship_street"`
	ShipCity             string           `gorm:"column:ship_city"                   json:"ship_city"`
	ShipProvince         string           `gorm:"column:ship_province"               json:"ship_province"`
	ShipPostal           string           `gorm:"column:ship_postal"                 json:"ship_postal"`
	Notes                string           `json:"notes"`
	Items                []OrderItem      `gorm:"foreignKey:OrderID"                 json:"items,omitempty"`
	PaymentTransactionID string           `json:"payment_transaction_id"`
	PaymentURL           string           `json:"payment_url"`
	PaymentStatus        PaymentStatus    `gorm:"default:'pending'"                  json:"payment_status"`
	PaidAt               *time.Time       `json:"paid_at"`
	CreatedAt            time.Time        `json:"created_at"`
	UpdatedAt            time.Time        `json:"updated_at"`
}

type OrderItem struct {
	ID        uuid.UUID              `gorm:"type:uuid;primaryKey"        json:"id"`
	OrderID   uuid.UUID              `gorm:"type:uuid;not null;index"    json:"-"`
	ProductID uuid.UUID              `gorm:"type:uuid;not null;index"    json:"product_id"`
	Product   productmodel.Product   `gorm:"foreignKey:ProductID"        json:"product,omitempty"`
	Quantity  int                    `gorm:"not null"                    json:"quantity"`
	UnitPrice float64                `gorm:"type:numeric(12,2);not null" json:"unit_price"`
	Subtotal  float64                `gorm:"type:numeric(12,2);not null" json:"subtotal"`
	CreatedAt time.Time              `json:"created_at"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	return nil
}
