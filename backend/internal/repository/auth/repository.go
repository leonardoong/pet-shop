package auth

import (
	"errors"

	usermodel "petshop/internal/model/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateCustomer(c *usermodel.Customer) error
	FindCustomerByEmail(email string) (*usermodel.Customer, error)
	FindCustomerByID(id uuid.UUID) (*usermodel.Customer, error)
	UpdateCustomer(c *usermodel.Customer) error

	FindAdminByEmail(email string) (*usermodel.Admin, error)
	FindAdminByID(id uuid.UUID) (*usermodel.Admin, error)
	UpdateAdmin(a *usermodel.Admin) error

	CreateRefreshToken(rt *usermodel.RefreshToken) error
	FindRefreshToken(tokenHash string) (*usermodel.RefreshToken, error)
	RevokeRefreshToken(id uuid.UUID) error
	RevokeAllUserTokens(userID uuid.UUID, userType usermodel.UserType) error

	CreateResetToken(rt *usermodel.ResetToken) error
	FindResetToken(tokenHash string) (*usermodel.ResetToken, error)
	MarkResetTokenUsed(id uuid.UUID) error
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) CreateCustomer(c *usermodel.Customer) error {
	return r.db.Create(c).Error
}

func (r *repository) FindCustomerByEmail(email string) (*usermodel.Customer, error) {
	var c usermodel.Customer
	err := r.db.
		Select("id, email, password_hash, full_name, phone, is_active, created_at, updated_at").
		Where("email = ?", email).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &c, err
}

func (r *repository) FindCustomerByID(id uuid.UUID) (*usermodel.Customer, error) {
	var c usermodel.Customer
	err := r.db.
		Select("id, email, full_name, phone, is_active, created_at, updated_at").
		Where("id = ?", id).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &c, err
}

func (r *repository) FindAdminByEmail(email string) (*usermodel.Admin, error) {
	var a usermodel.Admin
	err := r.db.
		Select("admins.id, admins.email, admins.password_hash, admins.full_name, admins.is_active, admins.created_at, admins.updated_at").
		Preload("Roles.Permissions").
		Where("admins.email = ?", email).First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &a, err
}

func (r *repository) FindAdminByID(id uuid.UUID) (*usermodel.Admin, error) {
	var a usermodel.Admin
	err := r.db.
		Select("admins.id, admins.email, admins.full_name, admins.is_active, admins.created_at, admins.updated_at").
		Preload("Roles.Permissions").
		Where("admins.id = ?", id).First(&a).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &a, err
}

func (r *repository) CreateRefreshToken(rt *usermodel.RefreshToken) error {
	return r.db.Create(rt).Error
}

func (r *repository) FindRefreshToken(tokenHash string) (*usermodel.RefreshToken, error) {
	var rt usermodel.RefreshToken
	err := r.db.Where("token_hash = ? AND revoked = false", tokenHash).First(&rt).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rt, err
}

func (r *repository) RevokeRefreshToken(id uuid.UUID) error {
	return r.db.Model(&usermodel.RefreshToken{}).
		Where("id = ?", id).Update("revoked", true).Error
}

func (r *repository) RevokeAllUserTokens(userID uuid.UUID, userType usermodel.UserType) error {
	return r.db.Model(&usermodel.RefreshToken{}).
		Where("user_id = ? AND user_type = ? AND revoked = false", userID, userType).
		Update("revoked", true).Error
}

func (r *repository) UpdateCustomer(c *usermodel.Customer) error {
	return r.db.Save(c).Error
}

func (r *repository) UpdateAdmin(a *usermodel.Admin) error {
	return r.db.Save(a).Error
}

func (r *repository) CreateResetToken(rt *usermodel.ResetToken) error {
	return r.db.Create(rt).Error
}

func (r *repository) FindResetToken(tokenHash string) (*usermodel.ResetToken, error) {
	var rt usermodel.ResetToken
	err := r.db.Where("token_hash = ? AND used = false", tokenHash).First(&rt).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rt, err
}

func (r *repository) MarkResetTokenUsed(id uuid.UUID) error {
	return r.db.Model(&usermodel.ResetToken{}).Where("id = ?", id).Update("used", true).Error
}
