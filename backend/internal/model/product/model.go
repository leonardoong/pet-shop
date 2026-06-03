package product

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"       json:"id"`
	Name      string    `gorm:"uniqueIndex;not null"        json:"name"`
	Slug      string    `gorm:"uniqueIndex;not null;index"  json:"slug"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"         json:"id"`
	CategoryID  uuid.UUID `gorm:"type:uuid;not null;index"     json:"category_id"`
	Category    Category  `gorm:"foreignKey:CategoryID"        json:"category,omitempty"`
	Name        string    `gorm:"not null"                     json:"name"`
	Slug        string    `gorm:"uniqueIndex;not null"         json:"slug"`
	Description string    `json:"description"`
	Price       float64   `gorm:"type:numeric(12,2);not null"  json:"price"`
	CostPrice   float64   `gorm:"type:numeric(12,2);default:0" json:"cost_price"`
	Stock       int       `gorm:"not null;default:0"           json:"stock"`
	SKU         string    `gorm:"uniqueIndex;not null"         json:"sku"`
	ImageURL    string    `json:"image_url"`
	IsActive    bool      `gorm:"default:true"                 json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
