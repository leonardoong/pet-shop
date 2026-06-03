package auth

import (
	"errors"

	authdto "petshop/internal/dto/auth"
	authsvc "petshop/internal/service/auth"
	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
)

type ResetHandler struct{ svc authsvc.ResetService }

func NewResetHandler(svc authsvc.ResetService) *ResetHandler {
	return &ResetHandler{svc: svc}
}

// ForgotPassword godoc
// @Summary      Request password reset
// @Description  Sends a password reset token (logs to console in dev)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      authdto.ForgotPasswordRequest  true  "Email"
// @Success      200      {object}  response.Response
// @Router       /auth/forgot-password [post]
func (h *ResetHandler) ForgotPassword(c *gin.Context) {
	var req authdto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	token, name, err := h.svc.ForgotPassword(req.Email)
	if err != nil {
		response.InternalError(c, "Failed to process request")
		return
	}

	if token == "" {
		response.OK(c, "If the email exists, a reset link has been sent", nil)
		return
	}

	resetURL := "http://localhost:3000/reset-password?token=" + token
	logResetEmail(req.Email, name, resetURL)

	response.OK(c, "If the email exists, a reset link has been sent", nil)
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Resets password using a valid reset token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      authdto.ResetPasswordRequest  true  "Reset payload"
// @Success      200      {object}  response.Response
// @Router       /auth/reset-password [post]
func (h *ResetHandler) ResetPassword(c *gin.Context) {
	var req authdto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Validation failed", err.Error())
		return
	}

	if err := h.svc.ResetPassword(req.Token, req.NewPassword); err != nil {
		if errors.Is(err, authsvc.ErrTokenExpired) {
			response.BadRequest(c, "Token expired or already used")
			return
		}
		response.InternalError(c, "Failed to reset password")
		return
	}

	response.OK(c, "Password has been reset", nil)
}

func logResetEmail(email, name, url string) {
	// In production this would send an actual email
	println("========================================")
	println("PASSWORD RESET TOKEN (dev mode)")
	println("To:", email, "(", name, ")")
	println("Reset URL:", url)
	println("========================================")
}
