package cart

import (
	"time"

	productmodel "petshop/internal/model/product"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"          json:"id"`
	CustomerID uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null" json:"customer_id"`
	Items      []CartItem `gorm:"foreignKey:CartID"              json:"items"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type CartItem struct {
	ID        uuid.UUID              `gorm:"type:uuid;primaryKey"     json:"id"`
	CartID    uuid.UUID              `gorm:"type:uuid;not null;index" json:"-"`
	ProductID uuid.UUID              `gorm:"type:uuid;not null;index" json:"product_id"`
	Product   productmodel.Product   `gorm:"foreignKey:ProductID"     json:"product"`
	Quantity  int                    `gorm:"not null"                 json:"quantity"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (ci *CartItem) BeforeCreate(tx *gorm.DB) error {
	if ci.ID == uuid.Nil {
		ci.ID = uuid.New()
	}
	return nil
}

func (c *Cart) Total() float64 {
	var total float64
	for _, item := range c.Items {
		total += item.Product.Price * float64(item.Quantity)
	}
	return total
}
