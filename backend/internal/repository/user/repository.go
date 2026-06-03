package user

import (
	"errors"

	userdto "petshop/internal/dto/user"
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

type AdminRepository interface {
	ListCustomers(f userdto.AdminCustomerFilter) ([]usermodel.Customer, int64, error)
	FindCustomerByID(id string) (*usermodel.Customer, error)
	UpdateCustomer(c *usermodel.Customer) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) (Repository, AdminRepository) {
	r := &repository{db: db}
	return r, r
}

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

func (r *repository) ListCustomers(f userdto.AdminCustomerFilter) ([]usermodel.Customer, int64, error) {
	q := r.db.Model(&usermodel.Customer{}).Select("id, email, full_name, phone, is_active, created_at, updated_at")

	if f.Search != "" {
		like := "%" + f.Search + "%"
		q = q.Where("full_name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", like, like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.Limit
	var customers []usermodel.Customer
	err := q.Order("created_at desc").Limit(f.Limit).Offset(offset).Find(&customers).Error
	return customers, total, err
}

func (r *repository) FindCustomerByID(id string) (*usermodel.Customer, error) {
	var c usermodel.Customer
	err := r.db.Select("id, email, full_name, phone, is_active, created_at, updated_at").
		Where("id = ?", id).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &c, err
}

func (r *repository) UpdateCustomer(c *usermodel.Customer) error {
	return r.db.Save(c).Error
}
