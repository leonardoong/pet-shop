package product

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockLogType string

const (
	StockLogPurchase   StockLogType = "purchase"
	StockLogSale       StockLogType = "sale"
	StockLogAdjustment StockLogType = "adjustment"
)

type StockLog struct {
	ID        uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID uuid.UUID    `gorm:"type:uuid;not null;index" json:"product_id"`
	Type      StockLogType `gorm:"type:varchar(20);not null" json:"type"`
	Quantity  int          `gorm:"not null" json:"quantity"`
	CostPrice float64      `gorm:"type:numeric(12,2);default:0" json:"cost_price"`
	TotalCost float64      `gorm:"type:numeric(12,2);default:0" json:"total_cost"`
	Note      string       `json:"note"`
	CreatedAt time.Time    `json:"created_at"`
}

func (sl *StockLog) BeforeCreate(tx *gorm.DB) error {
	if sl.ID == uuid.Nil {
		sl.ID = uuid.New()
	}
	return nil
}
