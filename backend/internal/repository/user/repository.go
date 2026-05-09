package user

import (
	"errors"

	usermodel "petshop/internal/model/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	ListAddresses(customerID uuid.UUID) ([]usermodel.Address, error)
	FindAddressByID(id, customerID uuid.UUID) (*usermodel.Address, error)
	CreateAddress(a *usermodel.Address) error
	UpdateAddress(a *usermodel.Address) error
	DeleteAddress(id, customerID uuid.UUID) error
	ClearDefaultAddress(customerID uuid.UUID) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) ListAddresses(customerID uuid.UUID) ([]usermodel.Address, error) {
	var addrs []usermodel.Address
	err := r.db.Where("customer_id = ?", customerID).
		Order("is_default desc, created_at asc").Find(&addrs).Error
	return addrs, err
}

func (r *repository) FindAddressByID(id, customerID uuid.UUID) (*usermodel.Address, error) {
	var a usermodel.Address
	err := r.db.Where("id = ? AND customer_id = ?", id, customerID).First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &a, err
}

func (r *repository) CreateAddress(a *usermodel.Address) error {
	return r.db.Create(a).Error
}

func (r *repository) UpdateAddress(a *usermodel.Address) error {
	return r.db.Save(a).Error
}

func (r *repository) DeleteAddress(id, customerID uuid.UUID) error {
	return r.db.Where("id = ? AND customer_id = ?", id, customerID).Delete(&usermodel.Address{}).Error
}

func (r *repository) ClearDefaultAddress(customerID uuid.UUID) error {
	return r.db.Model(&usermodel.Address{}).
		Where("customer_id = ? AND is_default = true", customerID).
		Update("is_default", false).Error
}
