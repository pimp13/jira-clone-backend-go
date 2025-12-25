package workspace

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/pimp13/jira-clone-backend-go/ent"
	entWorkspace "github.com/pimp13/jira-clone-backend-go/ent/workspace"
	"github.com/pimp13/jira-clone-backend-go/internal/module/fileupload"
	"github.com/pimp13/jira-clone-backend-go/pkg/res"
	"github.com/pimp13/jira-clone-backend-go/pkg/util"
)

type WorkspaceService interface {
	Index(ctx context.Context, userId uuid.UUID) *res.Response[[]*WorkspaceResponse]

	Create(
		ctx context.Context,
		bodyData CreateWorkspaceDto,
		file *multipart.FileHeader,
		userId uuid.UUID,
	) *res.Response[struct{}]
}

type workspaceService struct {
	client            *ent.Client
	fileUploadService fileupload.FileUploadService
}

func NewWorkspaceService(client *ent.Client) WorkspaceService {
	return &workspaceService{
		client:            client,
		fileUploadService: fileupload.NewFileUploadService("public/uploads/workspace", ""),
	}
}

func (s *workspaceService) Index(ctx context.Context, userId uuid.UUID) *res.Response[[]*WorkspaceResponse] {
	data, err := s.client.Workspace.Query().
		Where(entWorkspace.OwnerIDEQ(userId)).
		WithOwner().
		Order(entWorkspace.ByCreatedAt(sql.OrderDesc())).
		All(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return res.ErrorMessage[[]*WorkspaceResponse]("workspace is not found", http.StatusBadRequest)
		}
		return res.ErrorMessage[[]*WorkspaceResponse]("failed to get workspace")
	}

	finalData := make([]*WorkspaceResponse, 0, len(data))
	for _, ws := range data {
		finalData = append(finalData, ToWorkspaceResponse(ws))
	}

	return res.SuccessResponse(finalData, "")
}

func (s *workspaceService) Create(
	ctx context.Context,
	bodyData CreateWorkspaceDto,
	file *multipart.FileHeader,
	userId uuid.UUID,
) *res.Response[struct{}] {
	var slug string
	var err error

	if bodyData.Slug != nil {
		slug, err = s.generateUniqueSlug(ctx, *bodyData.Slug)
	} else {
		slug, err = s.generateUniqueSlug(ctx, bodyData.Name)
	}
	if err != nil {
		return res.ErrorMessage[struct{}]("failed to generate slug")
	}

	uploadResult, err := s.fileUploadService.UploadImage(ctx, file)
	if err != nil {
		return res.ErrorResponse[struct{}]("failed to upload file", err)
	}

	if _, err := s.client.Workspace.Create().
		SetName(bodyData.Name).
		SetImageURL(uploadResult.URL).
		SetSlug(slug).
		SetOwnerID(userId).
		Save(ctx); err != nil {
		_ = s.fileUploadService.DeleteImage(ctx, uploadResult.FilePath)
		return res.ErrorResponse[struct{}]("failed to save workspace", err)
	}

	return res.SuccessMessage("workspace is saved by successfully!", http.StatusCreated)
}

func (s *workspaceService) generateUniqueSlug(ctx context.Context, title string) (string, error) {
	baseSlug := util.Slugify(title)
	slug := baseSlug
	count := 1

	for true {
		exists, err := s.client.Workspace.Query().Where(entWorkspace.SlugEQ(slug)).Exist(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				break
			}
			return "", err
		}
		if !exists {
			break
		}
		slug = fmt.Sprintf("%s-%d", baseSlug, count)
		count++
	}

	return slug, nil
}
