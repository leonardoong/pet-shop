package user

import (
	"errors"

	userdto "petshop/internal/dto/user"
	"petshop/internal/middleware"
	usersvc "petshop/internal/service/user"
	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct{ svc usersvc.ProfileService }

func NewProfileHandler(svc usersvc.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

// GetProfile godoc
// @Summary      Get customer profile
// @Description  Returns the authenticated customer's profile
// @Tags         Customer Profile
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.ErrorResponse
// @Router       /customer/me [get]
func (h *ProfileHandler) Get(c *gin.Context) {
	customerID := middleware.MustGetCustomerID(c)
	profile, err := h.svc.GetProfile(customerID)
	if err != nil {
		if errors.Is(err, usersvc.ErrCustomerNotFound) {
			response.NotFound(c, "Customer not found")
			return
		}
		response.InternalError(c, "Failed to fetch profile")
		return
	}
	response.OK(c, "Success", profile)
}

// UpdateProfile godoc
// @Summary      Update profile
// @Description  Update name and phone number
// @Tags         Customer Profile
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      userdto.UpdateProfileRequest  true  "Profile update"
// @Success      200      {object}  response.Response
// @Router       /customer/me [put]
func (h *ProfileHandler) Update(c *gin.Context) {
	var req userdto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	customerID := middleware.MustGetCustomerID(c)
	profile, err := h.svc.UpdateProfile(customerID, req)
	if err != nil {
		if errors.Is(err, usersvc.ErrCustomerNotFound) {
			response.NotFound(c, "Customer not found")
			return
		}
		response.InternalError(c, "Failed to update profile")
		return
	}
	response.OK(c, "Profile updated", profile)
}

// ChangeEmail godoc
// @Summary      Change email
// @Description  Change email (requires current password)
// @Tags         Customer Profile
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      userdto.ChangeEmailRequest  true  "Email change"
// @Success      200      {object}  response.Response
// @Router       /customer/me/email [put]
func (h *ProfileHandler) ChangeEmail(c *gin.Context) {
	var req userdto.ChangeEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	customerID := middleware.MustGetCustomerID(c)
	if err := h.svc.ChangeEmail(customerID, req); err != nil {
		if errors.Is(err, usersvc.ErrWrongPassword) {
			response.BadRequest(c, "Password is incorrect")
			return
		}
		if errors.Is(err, usersvc.ErrEmailTaken) {
			response.Conflict(c, "Email already in use")
			return
		}
		response.InternalError(c, "Failed to change email")
		return
	}
	response.OK(c, "Email changed", nil)
}

// ChangePassword godoc
// @Summary      Change password
// @Description  Change password (requires current password)
// @Tags         Customer Profile
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body      userdto.ChangePasswordRequest  true  "Password change"
// @Success      200      {object}  response.Response
// @Router       /customer/me/password [put]
func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	var req userdto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	customerID := middleware.MustGetCustomerID(c)
	if err := h.svc.ChangePassword(customerID, req); err != nil {
		if errors.Is(err, usersvc.ErrWrongPassword) {
			response.BadRequest(c, "Password is incorrect")
			return
		}
		response.InternalError(c, "Failed to change password")
		return
	}
	response.OK(c, "Password changed", nil)
}
