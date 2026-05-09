package auth

import (
	"errors"

	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// CustomerRegister godoc
// @Summary      Register customer
// @Description  Create a new customer account
// @Tags         Customer Auth
// @Accept       json
// @Produce      json
// @Param        request  body      CustomerRegisterRequest                          true  "Register payload"
// @Success      201      {object}  response.Response{data=CustomerAuthResponse}     "Registration successful"
// @Failure      400      {object}  response.ErrorResponse                           "Validation failed"
// @Failure      409      {object}  response.ErrorResponse                           "Email already in use"
// @Failure      500      {object}  response.ErrorResponse
// @Router       /customer/auth/register [post]
func (h *Handler) CustomerRegister(c *gin.Context) {
	var req CustomerRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	resp, err := h.svc.RegisterCustomer(req)
	if err != nil {
		if errors.Is(err, ErrEmailTaken) {
			response.Conflict(c, "Email already in use")
			return
		}
		response.InternalError(c, "Registration failed")
		return
	}

	response.Created(c, "Registration successful", resp)
}

// CustomerLogin godoc
// @Summary      Login customer
// @Description  Authenticate a customer and return access + refresh tokens
// @Tags         Customer Auth
// @Accept       json
// @Produce      json
// @Param        request  body      CustomerLoginRequest                             true  "Login payload"
// @Success      200      {object}  response.Response{data=CustomerAuthResponse}    "Login successful"
// @Failure      400      {object}  response.ErrorResponse                          "Validation failed"
// @Failure      401      {object}  response.ErrorResponse                          "Invalid credentials"
// @Failure      403      {object}  response.ErrorResponse                          "Account inactive"
// @Failure      500      {object}  response.ErrorResponse
// @Router       /customer/auth/login [post]
func (h *Handler) CustomerLogin(c *gin.Context) {
	var req CustomerLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	resp, err := h.svc.LoginCustomer(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCreds) {
			response.Unauthorized(c, "Invalid email or password")
			return
		}
		if errors.Is(err, ErrAccountInactive) {
			response.Forbidden(c, "Account is inactive")
			return
		}
		response.InternalError(c, "Login failed")
		return
	}

	response.OK(c, "Login successful", resp)
}

// AdminLogin godoc
// @Summary      Login admin
// @Description  Authenticate an admin and return tokens with RBAC permissions
// @Tags         Admin Auth
// @Accept       json
// @Produce      json
// @Param        request  body      AdminLoginRequest                               true  "Login payload"
// @Success      200      {object}  response.Response{data=AdminAuthResponse}       "Login successful"
// @Failure      400      {object}  response.ErrorResponse                          "Validation failed"
// @Failure      401      {object}  response.ErrorResponse                          "Invalid credentials"
// @Failure      403      {object}  response.ErrorResponse                          "Account inactive"
// @Failure      500      {object}  response.ErrorResponse
// @Router       /admin/auth/login [post]
func (h *Handler) AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	resp, err := h.svc.LoginAdmin(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCreds) {
			response.Unauthorized(c, "Invalid email or password")
			return
		}
		if errors.Is(err, ErrAccountInactive) {
			response.Forbidden(c, "Account is inactive")
			return
		}
		response.InternalError(c, "Login failed")
		return
	}

	response.OK(c, "Login successful", resp)
}

// CustomerRefresh godoc
// @Summary      Refresh customer token
// @Description  Exchange a valid refresh token for a new access token (rotation)
// @Tags         Customer Auth
// @Accept       json
// @Produce      json
// @Param        request  body      RefreshRequest                                  true  "Refresh payload"
// @Success      200      {object}  response.Response{data=TokenPair}               "Token refreshed"
// @Failure      400      {object}  response.ErrorResponse                          "Validation failed"
// @Failure      401      {object}  response.ErrorResponse                          "Invalid or expired token"
// @Router       /customer/auth/refresh [post]
func (h *Handler) CustomerRefresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	pair, err := h.svc.RefreshCustomerToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "Invalid or expired token")
		return
	}

	response.OK(c, "Token refreshed", pair)
}

// AdminRefresh godoc
// @Summary      Refresh admin token
// @Description  Exchange a valid admin refresh token for a new access token (rotation)
// @Tags         Admin Auth
// @Accept       json
// @Produce      json
// @Param        request  body      RefreshRequest                                  true  "Refresh payload"
// @Success      200      {object}  response.Response{data=TokenPair}               "Token refreshed"
// @Failure      400      {object}  response.ErrorResponse                          "Validation failed"
// @Failure      401      {object}  response.ErrorResponse                          "Invalid or expired token"
// @Router       /admin/auth/refresh [post]
func (h *Handler) AdminRefresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	pair, err := h.svc.RefreshAdminToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "Invalid or expired token")
		return
	}

	response.OK(c, "Token refreshed", pair)
}

// Logout godoc
// @Summary      Logout
// @Description  Revoke the provided refresh token. Works for both customer and admin.
// @Tags         Customer Auth
// @Accept       json
// @Produce      json
// @Param        request  body      RefreshRequest           true  "Logout payload"
// @Success      200      {object}  response.Response        "Logged out"
// @Failure      400      {object}  response.ErrorResponse
// @Router       /customer/auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	_ = h.svc.Logout(req.RefreshToken)
	response.OK(c, "Logged out", nil)
}
