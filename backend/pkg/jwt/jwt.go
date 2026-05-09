package jwt

import (
	"errors"
	"time"

	"petshop/pkg/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	IssuerCustomer = "customer"
	IssuerAdmin    = "admin"
)

type CustomerClaims struct {
	jwt.RegisteredClaims
	CustomerID string `json:"cid"`
}

type AdminClaims struct {
	jwt.RegisteredClaims
	AdminID     string   `json:"aid"`
	Permissions []string `json:"perms"`
}

type Manager struct {
	cfg *config.JWTConfig
}

func NewManager(cfg *config.JWTConfig) *Manager {
	return &Manager{cfg: cfg}
}

// --- Customer tokens ---

func (m *Manager) GenerateCustomerAccess(customerID string) (string, error) {
	claims := CustomerClaims{
		CustomerID: customerID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    IssuerCustomer,
			Subject:   customerID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(m.cfg.AccessExpiryMinutes) * time.Minute)),
			ID:        uuid.NewString(),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(m.cfg.AccessSecret))
}

func (m *Manager) ValidateCustomerAccess(tokenStr string) (*CustomerClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomerClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(m.cfg.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*CustomerClaims)
	if !ok || claims.Issuer != IssuerCustomer {
		return nil, errors.New("invalid customer token")
	}
	return claims, nil
}

// --- Admin tokens ---

func (m *Manager) GenerateAdminAccess(adminID string, permissions []string) (string, error) {
	claims := AdminClaims{
		AdminID:     adminID,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    IssuerAdmin,
			Subject:   adminID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(m.cfg.AccessExpiryMinutes) * time.Minute)),
			ID:        uuid.NewString(),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(m.cfg.AccessSecret))
}

func (m *Manager) ValidateAdminAccess(tokenStr string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AdminClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(m.cfg.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*AdminClaims)
	if !ok || claims.Issuer != IssuerAdmin {
		return nil, errors.New("invalid admin token")
	}
	return claims, nil
}

// --- Refresh token (shared, distinguished by user_type in DB) ---

func (m *Manager) GenerateRefresh(subjectID, issuer string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   subjectID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(m.cfg.RefreshExpiryDays) * 24 * time.Hour)),
		ID:        uuid.NewString(),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(m.cfg.RefreshSecret))
}

func (m *Manager) ValidateRefresh(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(m.cfg.RefreshSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("invalid refresh token")
	}
	return claims, nil
}
