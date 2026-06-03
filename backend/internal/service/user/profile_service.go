package user

import (
	"errors"

	userdto "petshop/internal/dto/user"
	usermodel "petshop/internal/model/user"
	authrepo "petshop/internal/repository/auth"
	userrepo "petshop/internal/repository/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrWrongPassword   = errors.New("wrong password")
	ErrEmailTaken      = errors.New("email already in use")
)

type ProfileService interface {
	GetProfile(customerID uuid.UUID) (*usermodel.Customer, error)
	UpdateProfile(customerID uuid.UUID, req userdto.UpdateProfileRequest) (*usermodel.Customer, error)
	ChangeEmail(customerID uuid.UUID, req userdto.ChangeEmailRequest) error
	ChangePassword(customerID uuid.UUID, req userdto.ChangePasswordRequest) error
}

type profileService struct {
	repo     userrepo.AdminRepository
	authRepo authrepo.Repository
}

func NewProfileService(repo userrepo.AdminRepository, authRepo authrepo.Repository) ProfileService {
	return &profileService{repo: repo, authRepo: authRepo}
}

func (s *profileService) GetProfile(customerID uuid.UUID) (*usermodel.Customer, error) {
	c, err := s.repo.FindCustomerByID(customerID.String())
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrCustomerNotFound
	}
	return c, nil
}

func (s *profileService) UpdateProfile(customerID uuid.UUID, req userdto.UpdateProfileRequest) (*usermodel.Customer, error) {
	c, err := s.repo.FindCustomerByID(customerID.String())
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrCustomerNotFound
	}

	if req.FullName != "" {
		c.FullName = req.FullName
	}
	if req.Phone != "" {
		c.Phone = req.Phone
	}

	if err := s.repo.UpdateCustomer(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *profileService) ChangeEmail(customerID uuid.UUID, req userdto.ChangeEmailRequest) error {
	c, err := s.repo.FindCustomerByID(customerID.String())
	if err != nil {
		return err
	}
	if c == nil {
		return ErrCustomerNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(c.PasswordHash), []byte(req.Password)); err != nil {
		return ErrWrongPassword
	}

	existing, err := s.authRepo.FindCustomerByEmail(req.Email)
	if err != nil {
		return err
	}
	if existing != nil && existing.ID != c.ID {
		return ErrEmailTaken
	}

	c.Email = req.Email
	return s.repo.UpdateCustomer(c)
}

func (s *profileService) ChangePassword(customerID uuid.UUID, req userdto.ChangePasswordRequest) error {
	c, err := s.repo.FindCustomerByID(customerID.String())
	if err != nil {
		return err
	}
	if c == nil {
		return ErrCustomerNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(c.PasswordHash), []byte(req.OldPassword)); err != nil {
		return ErrWrongPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	c.PasswordHash = string(hash)
	return s.repo.UpdateCustomer(c)
}
