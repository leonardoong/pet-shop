package admin

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null"  json:"email"`
	PasswordHash string    `gorm:"not null"              json:"-"`
	FullName     string    `gorm:"not null"              json:"full_name"`
	IsActive     bool      `gorm:"default:true"          json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Roles []Role `gorm:"many2many:admin_roles;" json:"roles,omitempty"`
}

type Role struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string       `gorm:"uniqueIndex;not null"  json:"name"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`

	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null"  json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
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
