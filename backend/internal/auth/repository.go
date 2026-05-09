package auth

import (
	"errors"

	"petshop/internal/admin"
	"petshop/internal/customer"
	"petshop/internal/token"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	// Customer
	CreateCustomer(c *customer.Customer) error
	FindCustomerByEmail(email string) (*customer.Customer, error)
	FindCustomerByID(id uuid.UUID) (*customer.Customer, error)

	// Admin
	FindAdminByEmail(email string) (*admin.Admin, error)
	FindAdminByID(id uuid.UUID) (*admin.Admin, error)

	// Refresh tokens
	CreateRefreshToken(rt *token.RefreshToken) error
	FindRefreshToken(tokenHash string) (*token.RefreshToken, error)
	RevokeRefreshToken(id uuid.UUID) error
	RevokeAllUserTokens(userID uuid.UUID, userType token.UserType) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateCustomer(c *customer.Customer) error {
	return r.db.Create(c).Error
}

func (r *repository) FindCustomerByEmail(email string) (*customer.Customer, error) {
	var c customer.Customer
	err := r.db.
		Select("id, email, password_hash, full_name, phone, is_active, created_at, updated_at").
		Where("email = ?", email).
		First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &c, err
}

func (r *repository) FindCustomerByID(id uuid.UUID) (*customer.Customer, error) {
	var c customer.Customer
	err := r.db.
		Select("id, email, full_name, phone, is_active, created_at, updated_at").
		Where("id = ?", id).
		First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &c, err
}

func (r *repository) FindAdminByEmail(email string) (*admin.Admin, error) {
	var a admin.Admin
	err := r.db.
		Select("admins.id, admins.email, admins.password_hash, admins.full_name, admins.is_active, admins.created_at, admins.updated_at").
		Preload("Roles.Permissions").
		Where("admins.email = ?", email).
		First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &a, err
}

func (r *repository) FindAdminByID(id uuid.UUID) (*admin.Admin, error) {
	var a admin.Admin
	err := r.db.
		Select("admins.id, admins.email, admins.full_name, admins.is_active, admins.created_at, admins.updated_at").
		Preload("Roles.Permissions").
		Where("admins.id = ?", id).
		First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &a, err
}

func (r *repository) CreateRefreshToken(rt *token.RefreshToken) error {
	return r.db.Create(rt).Error
}

func (r *repository) FindRefreshToken(tokenHash string) (*token.RefreshToken, error) {
	var rt token.RefreshToken
	err := r.db.
		Where("token_hash = ? AND revoked = false", tokenHash).
		First(&rt).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rt, err
}

func (r *repository) RevokeRefreshToken(id uuid.UUID) error {
	return r.db.Model(&token.RefreshToken{}).
		Where("id = ?", id).
		Update("revoked", true).Error
}

func (r *repository) RevokeAllUserTokens(userID uuid.UUID, userType token.UserType) error {
	return r.db.Model(&token.RefreshToken{}).
		Where("user_id = ? AND user_type = ? AND revoked = false", userID, userType).
		Update("revoked", true).Error
}
