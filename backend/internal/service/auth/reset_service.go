package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	authrepo "petshop/internal/repository/auth"
	usermodel "petshop/internal/model/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrTokenExpired = errors.New("reset token expired")
)

type ResetService interface {
	ForgotPassword(email string) (string, string, error)
	ResetPassword(token, newPassword string) error
}

type resetService struct {
	repo authrepo.Repository
}

func NewResetService(repo authrepo.Repository) ResetService {
	return &resetService{repo: repo}
}

func (s *resetService) ForgotPassword(email string) (string, string, error) {
	typeStr := usermodel.UserTypeCustomer

	name := ""
	c, err := s.repo.FindCustomerByEmail(email)
	if err != nil {
		return "", "", err
	}
	if c == nil {
		a, err := s.repo.FindAdminByEmail(email)
		if err != nil {
			return "", "", err
		}
		if a == nil {
			return "", "", nil
		}
		name = a.FullName
		typeStr = usermodel.UserTypeAdmin
	} else {
		name = c.FullName
	}

	token := uuid.NewString()
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	rt := &usermodel.ResetToken{
		Email:     email,
		TokenHash: tokenHash,
		UserType:  typeStr,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	if err := s.repo.CreateResetToken(rt); err != nil {
		return "", "", err
	}

	return token, name, nil
}

func (s *resetService) ResetPassword(token, newPassword string) error {
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	rt, err := s.repo.FindResetToken(tokenHash)
	if err != nil {
		return err
	}
	if rt == nil || rt.IsExpired() || rt.Used {
		return ErrTokenExpired
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if rt.UserType == usermodel.UserTypeAdmin {
		a, err := s.repo.FindAdminByEmail(rt.Email)
		if err != nil || a == nil {
			return errors.New("admin not found")
		}
		a.PasswordHash = string(passwordHash)
		if err := s.repo.UpdateAdmin(a); err != nil {
			return err
		}
	} else {
		c, err := s.repo.FindCustomerByEmail(rt.Email)
		if err != nil || c == nil {
			return errors.New("customer not found")
		}
		c.PasswordHash = string(passwordHash)
		if err := s.repo.UpdateCustomer(c); err != nil {
			return err
		}
	}

	return s.repo.MarkResetTokenUsed(rt.ID)
}
