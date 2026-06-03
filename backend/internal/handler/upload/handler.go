package upload

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"petshop/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var validExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
}

type Handler struct {
	uploadDir string
	maxSize   int64
	baseURL   string
}

func NewHandler(uploadDir, baseURL string, maxSize int64) *Handler {
	os.MkdirAll(uploadDir, 0755)
	return &Handler{uploadDir: uploadDir, baseURL: baseURL, maxSize: maxSize}
}

// Upload godoc
// @Summary      Upload image
// @Description  Upload a product image (max 5MB, jpg/png/webp)
// @Tags         Admin Upload
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Image file"
// @Success      200  {object}  response.Response{data=UploadResponse}
// @Router       /admin/upload [post]
func (h *Handler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "No file provided")
		return
	}
	defer file.Close()

	if header.Size > h.maxSize {
		response.BadRequest(c, fmt.Sprintf("File too large. Max %d MB", h.maxSize/(1024*1024)))
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !validExtensions[ext] {
		response.BadRequest(c, "Invalid file type. Allowed: jpg, png, webp")
		return
	}

	filename := uuid.NewString() + ext
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

	url := fmt.Sprintf("%s/uploads/%s", h.baseURL, filename)
	response.OK(c, "Upload successful", gin.H{"url": url, "filename": filename})
}
