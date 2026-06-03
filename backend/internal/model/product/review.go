package product

import (
	"time"

	usermodel "petshop/internal/model/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ID         uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID  uuid.UUID       `gorm:"type:uuid;not null;index" json:"product_id"`
	CustomerID uuid.UUID       `gorm:"type:uuid;not null;index" json:"customer_id"`
	Customer   usermodel.Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Rating     int             `gorm:"not null" json:"rating"`
	Comment    string          `json:"comment"`
	IsApproved bool            `gorm:"default:false" json:"is_approved"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
