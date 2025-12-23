package fileupload

import (
	"context"
	"mime/multipart"

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

// DeleteImage implements FileUploadService.
func (f *fileUploadService) DeleteImage(ctx context.Context, imageURL string) error {
	panic("unimplemented")
}

// UploadImage implements FileUploadService.
func (f *fileUploadService) UploadImage(ctx context.Context, image *multipart.FileHeader) (*FileUploadResultDto, error) {
	panic("unimplemented")
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
