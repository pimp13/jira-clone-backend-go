package project

import "context"

type ProjectService interface {
	Index(ctx context.Context)
}

type ProjectServiceImpl struct {}

func NewProjectService() ProjectService {
	return &ProjectServiceImpl{}
}

func (s *ProjectServiceImpl) Index(ctx context.Context) {}
