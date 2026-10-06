package upload_image

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/response"

	"github.com/MamangRust/monolith-pointofsale-pkg/logger"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ImageUploads interface {
	EnsureUploadDirectory(uploadDir string) error
	ProcessImageUpload(c echo.Context, file *multipart.FileHeader) (string, error)
	CleanupImageOnFailure(imagePath string)
	SaveUploadedFile(file *multipart.FileHeader, dst string) error
}

type ImageUpload struct {
	logger logger.LoggerInterface
}

func NewImageUpload(logger logger.LoggerInterface) ImageUploads {
	return &ImageUpload{logger: logger}
}

func (h *ImageUpload) EnsureUploadDirectory(uploadDir string) error {
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			h.logger.Error("Failed to create upload directory",
				zap.String("directory", uploadDir),
				zap.Error(err),
			)
			return err
		}
	}
	return nil
}

func (h *ImageUpload) ProcessImageUpload(c echo.Context, file *multipart.FileHeader) (string, error) {
	allowedTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedTypes[ext] {
		return "", c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "invalid_image_type",
			Message: "Only JPG, JPEG, and PNG",
			Code:    http.StatusBadRequest,
		})
	}

	if file.Size > 5<<20 {
		return "", c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Status:  "invalid_image_size",
			Message: "Image size must be less than 5MB",
			Code:    http.StatusBadRequest,
		})
	}

	uploadDir := "uploads/products"
	if err := h.EnsureUploadDirectory(uploadDir); err != nil {
		return "", c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "server_error",
			Message: "Failed to prepare storage for upload",
			Code:    http.StatusInternalServerError,
		})
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	imagePath := filepath.Join(uploadDir, filename)

	if err := h.SaveUploadedFile(file, imagePath); err != nil {
		h.logger.Error("Failed to save image",
			zap.String("path", imagePath),
			zap.Error(err),
		)
		return "", c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Status:  "upload_failed",
			Message: "Failed to save uploaded image",
			Code:    http.StatusInternalServerError,
		})
	}

	h.logger.Debug("Successfully saved uploaded file",
		zap.String("path", imagePath),
		zap.Int64("size", file.Size),
	)

	return imagePath, nil
}

func (h *ImageUpload) CleanupImageOnFailure(imagePath string) {
	if removeErr := os.Remove(imagePath); removeErr != nil {
		h.logger.Debug("Failed to clean up uploaded file after failure",
			zap.String("path", imagePath),
			zap.Error(removeErr),
		)
	}
}

func (h *ImageUpload) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	if stat, err := os.Stat(dst); err != nil || stat.Size() == 0 {
		return fmt.Errorf("failed to verify file write: %w", err)
	}

	return nil
}
