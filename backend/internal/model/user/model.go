package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"  json:"id"`
	Email        string    `gorm:"uniqueIndex;not null"  json:"email"`
	PasswordHash string    `gorm:"not null"              json:"-"`
	FullName     string    `gorm:"not null"              json:"full_name"`
	Phone        string    `json:"phone"`
	AvatarURL    *string   `json:"avatar_url"`
	IsActive     bool      `gorm:"default:true"          json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (c *Customer) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type Admin struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"  json:"id"`
	Email        string    `gorm:"uniqueIndex;not null"  json:"email"`
	PasswordHash string    `gorm:"not null"              json:"-"`
	FullName     string    `gorm:"not null"              json:"full_name"`
	IsActive     bool      `gorm:"default:true"          json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Roles        []Role    `gorm:"many2many:admin_roles;" json:"roles,omitempty"`
}

func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// EffectivePermissions returns deduplicated permission names across all roles.
func (a *Admin) EffectivePermissions() []string {
	seen := make(map[string]struct{})
	result := []string{}
	for _, role := range a.Roles {
		for _, perm := range role.Permissions {
			if _, ok := seen[perm.Name]; !ok {
				seen[perm.Name] = struct{}{}
				result = append(result, perm.Name)
			}
		}
	}
	return result
}

type Role struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey"  json:"id"`
	Name        string       `gorm:"uniqueIndex;not null"  json:"name"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"  json:"id"`
	Name        string    `gorm:"uniqueIndex;not null"  json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type Address struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"      json:"id"`
	CustomerID    uuid.UUID `gorm:"type:uuid;not null;index"  json:"customer_id"`
	Label         string    `json:"label"`
	RecipientName string    `gorm:"not null"                  json:"recipient_name"`
	Phone         string    `gorm:"not null"                  json:"phone"`
	Street        string    `gorm:"not null"                  json:"street"`
	City          string    `gorm:"not null"                  json:"city"`
	Province      string    `gorm:"not null"                  json:"province"`
	PostalCode    string    `gorm:"not null"                  json:"postal_code"`
	IsDefault     bool      `gorm:"default:false"             json:"is_default"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// UserType discriminates customer vs admin refresh tokens.
type UserType string

const (
	UserTypeCustomer UserType = "customer"
	UserTypeAdmin    UserType = "admin"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_refresh_user"`
	UserType  UserType  `gorm:"type:varchar(10);not null;index:idx_refresh_user"`
	TokenHash string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Revoked   bool      `gorm:"default:false"`
	CreatedAt time.Time
}

func (r *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (r *RefreshToken) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

type ResetToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email     string    `gorm:"not null;index"`
	TokenHash string    `gorm:"uniqueIndex;not null"`
	UserType  UserType  `gorm:"type:varchar(10);not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
	CreatedAt time.Time
}

func (r *ResetToken) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (r *ResetToken) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}
