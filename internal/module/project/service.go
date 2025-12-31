package project

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/pimp13/jira-clone-backend-go/ent"
	"github.com/pimp13/jira-clone-backend-go/internal/module/fileupload"
	dto "github.com/pimp13/jira-clone-backend-go/internal/module/project/dto"
	"github.com/pimp13/jira-clone-backend-go/pkg/logger"
	"github.com/pimp13/jira-clone-backend-go/pkg/res"
)

type ProjectService interface {
	Index(ctx context.Context)

	Create(
		ctx context.Context,
		bodyData dto.CreateProjectDto,
		file *multipart.FileHeader,
	) *res.Response[struct{}]
}

type projectService struct {
	client            *ent.Client
	logger            logger.Logger
	fileUploadService fileupload.FileUploadService
}

func NewProjectService(client *ent.Client, logger logger.Logger) ProjectService {
	return &projectService{
		client:            client,
		logger:            logger,
		fileUploadService: fileupload.NewFileUploadService("public/uploads/project", ""),
	}
}

func (s *projectService) Index(ctx context.Context) {}

func (s *projectService) Create(
	ctx context.Context,
	bodyData dto.CreateProjectDto,
	file *multipart.FileHeader,
) *res.Response[struct{}] {
	var imageURL *string = nil
	var filePath *string = nil

	if file != nil {
		uploadResult, err := s.fileUploadService.UploadImage(ctx, file)
		if err != nil {
			return res.ErrorResponse[struct{}]("failed to upload file", err)
		}
		imageURL = &uploadResult.URL
		filePath = &uploadResult.FilePath
	}

	builder := s.client.Project.Create().
		SetName(bodyData.Name).
		SetWorkspaceID(bodyData.WorkspaceId).
		SetNillableImageURL(imageURL).
		SetNillableDescription(bodyData.Description).
		SetNillableIsActive(bodyData.IsActive)

	if _, err := builder.Save(ctx); err != nil {
		if file != nil && filePath != nil {
			_ = s.fileUploadService.DeleteImage(ctx, *filePath)
		}
		return res.ErrorResponse[struct{}]("failed to save project", err)
	}
	return res.SuccessMessage("workspace is saved by successfully!", http.StatusCreated)
}
