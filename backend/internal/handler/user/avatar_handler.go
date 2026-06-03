package user

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"petshop/internal/middleware"
	userrepo "petshop/internal/repository/user"
	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AvatarHandler struct {
	repo     userrepo.AdminRepository
	uploadDir string
	baseURL  string
}

func NewAvatarHandler(repo userrepo.AdminRepository, uploadDir, baseURL string) *AvatarHandler {
	os.MkdirAll(uploadDir, 0755)
	return &AvatarHandler{repo: repo, uploadDir: uploadDir, baseURL: baseURL}
}

// UploadAvatar godoc
// @Summary      Upload customer avatar
// @Description  Upload and crop avatar image (max 2MB, jpg/png/webp)
// @Tags         Customer Profile
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Avatar image"
// @Success      200   {object}  response.Response
// @Router       /customer/me/avatar [post]
func (h *AvatarHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "No file provided")
		return
	}
	defer file.Close()

	if header.Size > 2*1024*1024 {
		response.BadRequest(c, "File too large. Max 2 MB")
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		response.BadRequest(c, "Invalid file type. Allowed: jpg, png, webp")
		return
	}

	customerID := middleware.MustGetCustomerID(c)
	cust, err := h.repo.FindCustomerByID(customerID.String())
	if err != nil || cust == nil {
		response.InternalError(c, "Customer not found")
		return
	}

	filename := "avatar_" + uuid.NewString() + ext
	outPath := filepath.Join(h.uploadDir, filename)

	out, err := os.Create(outPath)
	if err != nil {
		response.InternalError(c, "Failed to save file")
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		response.InternalError(c, "Failed to save file")
		return
	}

	url := fmt.Sprintf("%s/uploads/avatars/%s", h.baseURL, filename)
	cust.AvatarURL = &url
	h.repo.UpdateCustomer(cust)

	response.OK(c, "Avatar updated", gin.H{"avatar_url": url})
}
