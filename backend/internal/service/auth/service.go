package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	authdto "petshop/internal/dto/auth"
	usermodel "petshop/internal/model/user"
	authrepo "petshop/internal/repository/auth"
	jwtpkg "petshop/pkg/jwt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken      = errors.New("email already in use")
	ErrInvalidCreds    = errors.New("invalid email or password")
	ErrAccountInactive = errors.New("account is inactive")
	ErrInvalidToken    = errors.New("invalid or expired token")
)

type Service interface {
	RegisterCustomer(req authdto.CustomerRegisterRequest) (*authdto.CustomerAuthResponse, error)
	LoginCustomer(req authdto.CustomerLoginRequest) (*authdto.CustomerAuthResponse, error)
	LoginAdmin(req authdto.AdminLoginRequest) (*authdto.AdminAuthResponse, error)
	RefreshCustomerToken(refreshToken string) (*authdto.TokenPair, error)
	RefreshAdminToken(refreshToken string) (*authdto.TokenPair, error)
	Logout(refreshToken string) error
}

type service struct {
	repo       authrepo.Repository
	jwtManager *jwtpkg.Manager
	accessExp  int
}

func NewService(repo authrepo.Repository, jwtManager *jwtpkg.Manager, accessExpiryMinutes int) Service {
	return &service{repo: repo, jwtManager: jwtManager, accessExp: accessExpiryMinutes}
}

func (s *service) RegisterCustomer(req authdto.CustomerRegisterRequest) (*authdto.CustomerAuthResponse, error) {
	existing, err := s.repo.FindCustomerByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	c := &usermodel.Customer{
		Email:        req.Email,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		Phone:        req.Phone,
	}
	if err := s.repo.CreateCustomer(c); err != nil {
		return nil, err
	}

	return s.buildCustomerAuthResponse(c)
}

func (s *service) LoginCustomer(req authdto.CustomerLoginRequest) (*authdto.CustomerAuthResponse, error) {
	c, err := s.repo.FindCustomerByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrInvalidCreds
	}
	if !c.IsActive {
		return nil, ErrAccountInactive
	}
	if err := bcrypt.CompareHashAndPassword([]byte(c.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCreds
	}

	return s.buildCustomerAuthResponse(c)
}

func (s *service) LoginAdmin(req authdto.AdminLoginRequest) (*authdto.AdminAuthResponse, error) {
	a, err := s.repo.FindAdminByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, ErrInvalidCreds
	}
	if !a.IsActive {
		return nil, ErrAccountInactive
	}
	if err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCreds
	}

	perms := a.EffectivePermissions()
	accessToken, err := s.jwtManager.GenerateAdminAccess(a.ID.String(), perms)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefresh(a.ID.String(), jwtpkg.IssuerAdmin)
	if err != nil {
		return nil, err
	}

	if err := s.storeRefreshToken(a.ID.String(), usermodel.UserTypeAdmin, refreshToken); err != nil {
		return nil, err
	}

	return &authdto.AdminAuthResponse{
		Admin:       authdto.AdminProfile{ID: a.ID.String(), FullName: a.FullName, Email: a.Email},
		Tokens:      authdto.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken, TokenType: "Bearer", ExpiresIn: s.accessExp * 60},
		Permissions: perms,
	}, nil
}

func (s *service) RefreshCustomerToken(refreshToken string) (*authdto.TokenPair, error) {
	claims, err := s.jwtManager.ValidateRefresh(refreshToken)
	if err != nil || claims.Issuer != jwtpkg.IssuerCustomer {
		return nil, ErrInvalidToken
	}

	rt, err := s.findAndRevokeToken(refreshToken)
	if err != nil {
		return nil, err
	}

	c, err := s.repo.FindCustomerByID(rt.UserID)
	if err != nil || c == nil || !c.IsActive {
		return nil, ErrInvalidToken
	}

	accessToken, err := s.jwtManager.GenerateCustomerAccess(c.ID.String())
	if err != nil {
		return nil, err
	}

	newRefresh, err := s.jwtManager.GenerateRefresh(c.ID.String(), jwtpkg.IssuerCustomer)
	if err != nil {
		return nil, err
	}
	if err := s.storeRefreshToken(c.ID.String(), usermodel.UserTypeCustomer, newRefresh); err != nil {
		return nil, err
	}

	return &authdto.TokenPair{AccessToken: accessToken, RefreshToken: newRefresh, TokenType: "Bearer", ExpiresIn: s.accessExp * 60}, nil
}

func (s *service) RefreshAdminToken(refreshToken string) (*authdto.TokenPair, error) {
	claims, err := s.jwtManager.ValidateRefresh(refreshToken)
	if err != nil || claims.Issuer != jwtpkg.IssuerAdmin {
		return nil, ErrInvalidToken
	}

	rt, err := s.findAndRevokeToken(refreshToken)
	if err != nil {
		return nil, err
	}

	a, err := s.repo.FindAdminByID(rt.UserID)
	if err != nil || a == nil || !a.IsActive {
		return nil, ErrInvalidToken
	}

	perms := a.EffectivePermissions()
	accessToken, err := s.jwtManager.GenerateAdminAccess(a.ID.String(), perms)
	if err != nil {
		return nil, err
	}

	newRefresh, err := s.jwtManager.GenerateRefresh(a.ID.String(), jwtpkg.IssuerAdmin)
	if err != nil {
		return nil, err
	}
	if err := s.storeRefreshToken(a.ID.String(), usermodel.UserTypeAdmin, newRefresh); err != nil {
		return nil, err
	}

	return &authdto.TokenPair{AccessToken: accessToken, RefreshToken: newRefresh, TokenType: "Bearer", ExpiresIn: s.accessExp * 60}, nil
}

func (s *service) Logout(refreshToken string) error {
	hash := hashToken(refreshToken)
	rt, err := s.repo.FindRefreshToken(hash)
	if err != nil || rt == nil {
		return nil
	}
	return s.repo.RevokeRefreshToken(rt.ID)
}

func (s *service) buildCustomerAuthResponse(c *usermodel.Customer) (*authdto.CustomerAuthResponse, error) {
	pair, err := s.issueCustomerTokenPair(c.ID.String())
	if err != nil {
		return nil, err
	}
	return &authdto.CustomerAuthResponse{
		Customer: authdto.CustomerProfile{ID: c.ID.String(), FullName: c.FullName, Email: c.Email, Phone: c.Phone},
		Tokens:   *pair,
	}, nil
}

func (s *service) issueCustomerTokenPair(customerID string) (*authdto.TokenPair, error) {
	accessToken, err := s.jwtManager.GenerateCustomerAccess(customerID)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtManager.GenerateRefresh(customerID, jwtpkg.IssuerCustomer)
	if err != nil {
		return nil, err
	}
	if err := s.storeRefreshToken(customerID, usermodel.UserTypeCustomer, refreshToken); err != nil {
		return nil, err
	}
	return &authdto.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken, TokenType: "Bearer", ExpiresIn: s.accessExp * 60}, nil
}

func (s *service) storeRefreshToken(userID string, userType usermodel.UserType, rawToken string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	rt := &usermodel.RefreshToken{
		UserID:    uid,
		UserType:  userType,
		TokenHash: hashToken(rawToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	return s.repo.CreateRefreshToken(rt)
}

func (s *service) findAndRevokeToken(rawToken string) (*usermodel.RefreshToken, error) {
	hash := hashToken(rawToken)
	rt, err := s.repo.FindRefreshToken(hash)
	if err != nil {
		return nil, err
	}
	if rt == nil || rt.IsExpired() {
		return nil, ErrInvalidToken
	}
	if err := s.repo.RevokeRefreshToken(rt.ID); err != nil {
		return nil, err
	}
	return rt, nil
}

func hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}
