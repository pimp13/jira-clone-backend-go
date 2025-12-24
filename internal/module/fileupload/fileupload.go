package fileupload

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pimp13/jira-clone-backend-go/internal/infrastructure/config"
)

type FileUploadService interface {
	UploadImage(ctx context.Context, image *multipart.FileHeader) (*FileUploadResultDto, error)

	DeleteImage(ctx context.Context, imageURL string) error
}

type fileUploadService struct {
	uploadDir string
	baseURL   string
}

func NewFileUploadService(uploadDir string, baseURL ...string) FileUploadService {
	var baseUrl string
	if baseURL[0] != "" {
		baseUrl = baseURL[0]
	} else {
		baseUrl = config.Envs.App.Url
	}

	return &fileUploadService{
		uploadDir: uploadDir,
		baseURL:   baseUrl,
	}
}

func (s *fileUploadService) UploadImage(ctx context.Context, image *multipart.FileHeader) (*FileUploadResultDto, error) {
	src, err := image.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return nil, err
	}

	ext := filepath.Ext(image.Filename)
	fileName := fmt.Sprintf("apophis_%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(s.uploadDir, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	urlPath := "/" + filepath.ToSlash(filepath.Join(s.uploadDir, fileName))

	return &FileUploadResultDto{
		URL:      s.baseURL + urlPath,
		FilePath: filePath,
	}, nil
}

func (s *fileUploadService) DeleteImage(ctx context.Context, imageURL string) error {
	if imageURL == "" {
		return nil
	}

	relPath := strings.TrimPrefix(imageURL, s.baseURL)
	if relPath == "" {
		return nil
	}

	imagePath := "." + relPath
	return os.Remove(imagePath)
}
