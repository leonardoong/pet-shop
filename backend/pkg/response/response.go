package response

import "github.com/gin-gonic/gin"

// Response is the standard success envelope for all API responses.
// @Description Standard API success response
type Response struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
}

// ErrorResponse is the standard error envelope.
// @Description Standard API error response
type ErrorResponse struct {
	Success bool     `json:"success"          example:"false"`
	Message string   `json:"message"          example:"Validation failed"`
	Errors  []string `json:"errors,omitempty" example:"field is required"`
}

func OK(c *gin.Context, message string, data any) {
	c.JSON(200, Response{Success: true, Message: message, Data: data})
}

func Created(c *gin.Context, message string, data any) {
	c.JSON(201, Response{Success: true, Message: message, Data: data})
}

func BadRequest(c *gin.Context, message string, errors ...string) {
	c.JSON(400, ErrorResponse{Success: false, Message: message, Errors: errors})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(401, ErrorResponse{Success: false, Message: message})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(403, ErrorResponse{Success: false, Message: message})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(404, ErrorResponse{Success: false, Message: message})
}

func Conflict(c *gin.Context, message string) {
	c.JSON(409, ErrorResponse{Success: false, Message: message})
}

func InternalError(c *gin.Context, message string) {
	c.JSON(500, ErrorResponse{Success: false, Message: message})
}
