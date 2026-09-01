package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"new_project/internal/models"

	"github.com/google/uuid"
)

type MediaService interface {
	Upload(ctx context.Context, file *multipart.FileHeader, folder string) (string, error)
}

type localMediaService struct {
	baseUploadPath string
}

// @inject
func NewMediaService() MediaService {
	return &localMediaService{
		baseUploadPath: "uploads",
	}
}

func (s *localMediaService) Upload(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (string, error) {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true,
	}
	if !allowedExts[ext] {
		return "", models.NewAppError(models.ErrorTypeValidation, "invalid file extension. Only images allowed.", string(models.ErrCodeInvalidInput))
	}

	cleanFolder := filepath.Clean(folder)

	if strings.HasPrefix(cleanFolder, "..") || filepath.IsAbs(cleanFolder) {
		return "", models.NewAppError(models.ErrorTypeValidation, "invalid folder path", string(models.ErrCodeInvalidInput))
	}

	parts := strings.Split(cleanFolder, string(filepath.Separator))
	rootFolder := parts[0]

	allowedRootFolders := map[string]bool{
		"avatars": true, "clubs": true, "proofs": true, "misc": true,
	}

	if !allowedRootFolders[rootFolder] {
		cleanFolder = filepath.Join("misc", cleanFolder)
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	newFileName := uuid.New().String() + ext

	targetDir := filepath.Join(s.baseUploadPath, cleanFolder)
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	targetFilePath := filepath.Join(targetDir, newFileName)
	dst, err := os.Create(targetFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create target file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file content: %w", err)
	}

	webPath := fmt.Sprintf("/%s/%s/%s", s.baseUploadPath, filepath.ToSlash(cleanFolder), newFileName)
	return webPath, nil
}
