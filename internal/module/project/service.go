package project

import (
	"context"

	"github.com/pimp13/jira-clone-backend-go/ent"
	dto "github.com/pimp13/jira-clone-backend-go/internal/module/project/dto"
	"github.com/pimp13/jira-clone-backend-go/pkg/logger"
)

type ProjectService interface {
	Index(ctx context.Context)
}

type projectService struct {
	client *ent.Client
	logger logger.Logger
}

func NewProjectService(client *ent.Client, logger logger.Logger) ProjectService {
	return &projectService{
		client,
		logger,
	}
}

func (s *projectService) Index(ctx context.Context) {}

func (s *projectService) Create(
	ctx context.Context,
	bodyData dto.CreateProjectDto,
) {
}
