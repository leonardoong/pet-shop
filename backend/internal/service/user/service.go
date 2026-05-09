package user

import (
	"errors"

	userdto "petshop/internal/dto/user"
	usermodel "petshop/internal/model/user"
	userrepo "petshop/internal/repository/user"

	"github.com/google/uuid"
)

var ErrAddressNotFound = errors.New("address not found")

type AddressService interface {
	ListAddresses(customerID uuid.UUID) ([]usermodel.Address, error)
	CreateAddress(customerID uuid.UUID, req userdto.AddressCreateRequest) (*usermodel.Address, error)
	UpdateAddress(id, customerID uuid.UUID, req userdto.AddressUpdateRequest) (*usermodel.Address, error)
	DeleteAddress(id, customerID uuid.UUID) error
	SetDefaultAddress(id, customerID uuid.UUID) (*usermodel.Address, error)
}

type addressService struct{ repo userrepo.Repository }

func NewAddressService(repo userrepo.Repository) AddressService {
	return &addressService{repo: repo}
}

func (s *addressService) ListAddresses(customerID uuid.UUID) ([]usermodel.Address, error) {
	return s.repo.ListAddresses(customerID)
}

func (s *addressService) CreateAddress(customerID uuid.UUID, req userdto.AddressCreateRequest) (*usermodel.Address, error) {
	if req.IsDefault {
		if err := s.repo.ClearDefaultAddress(customerID); err != nil {
			return nil, err
		}
	}

	a := &usermodel.Address{
		CustomerID:    customerID,
		Label:         req.Label,
		RecipientName: req.RecipientName,
		Phone:         req.Phone,
		Street:        req.Street,
		City:          req.City,
		Province:      req.Province,
		PostalCode:    req.PostalCode,
		IsDefault:     req.IsDefault,
	}

	if err := s.repo.CreateAddress(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *addressService) UpdateAddress(id, customerID uuid.UUID, req userdto.AddressUpdateRequest) (*usermodel.Address, error) {
	a, err := s.repo.FindAddressByID(id, customerID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, ErrAddressNotFound
	}

	if req.Label != "" {
		a.Label = req.Label
	}
	if req.RecipientName != "" {
		a.RecipientName = req.RecipientName
	}
	if req.Phone != "" {
		a.Phone = req.Phone
	}
	if req.Street != "" {
		a.Street = req.Street
	}
	if req.City != "" {
		a.City = req.City
	}
	if req.Province != "" {
		a.Province = req.Province
	}
	if req.PostalCode != "" {
		a.PostalCode = req.PostalCode
	}
	if req.IsDefault != nil {
		if *req.IsDefault && !a.IsDefault {
			if err := s.repo.ClearDefaultAddress(customerID); err != nil {
				return nil, err
			}
		}
		a.IsDefault = *req.IsDefault
	}

	if err := s.repo.UpdateAddress(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *addressService) DeleteAddress(id, customerID uuid.UUID) error {
	a, err := s.repo.FindAddressByID(id, customerID)
	if err != nil {
		return err
	}
	if a == nil {
		return ErrAddressNotFound
	}
	return s.repo.DeleteAddress(id, customerID)
}

func (s *addressService) SetDefaultAddress(id, customerID uuid.UUID) (*usermodel.Address, error) {
	a, err := s.repo.FindAddressByID(id, customerID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, ErrAddressNotFound
	}

	if err := s.repo.ClearDefaultAddress(customerID); err != nil {
		return nil, err
	}
	a.IsDefault = true
	if err := s.repo.UpdateAddress(a); err != nil {
		return nil, err
	}
	return a, nil
}
