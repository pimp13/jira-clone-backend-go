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
	workspace "github.com/pimp13/jira-clone-backend-go/internal/module/workspace/dto"
	"github.com/pimp13/jira-clone-backend-go/pkg/res"
	"github.com/pimp13/jira-clone-backend-go/pkg/util"
)

type WorkspaceService interface {
	Index(
		ctx context.Context,
		userId uuid.UUID,
	) *res.Response[[]*workspace.WorkspaceResponse]

	ShowById(
		ctx context.Context,
		workspaceId uuid.UUID,
		userId uuid.UUID,
	) *res.Response[*workspace.WorkspaceResponse]

	Create(
		ctx context.Context,
		bodyData workspace.CreateWorkspaceDto,
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

func (s *workspaceService) Index(
	ctx context.Context,
	userId uuid.UUID,
) *res.Response[[]*workspace.WorkspaceResponse] {
	initData, err := s.client.Workspace.Query().
		Where(entWorkspace.OwnerIDEQ(userId)).
		WithOwner().
		Order(entWorkspace.ByCreatedAt(sql.OrderDesc())).
		All(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return res.ErrorMessage[[]*workspace.WorkspaceResponse](
				"workspace is not found",
				http.StatusBadRequest,
			)
		}
		return res.ErrorMessage[[]*workspace.WorkspaceResponse]("failed to get workspace")
	}

	finalData := make([]*workspace.WorkspaceResponse, 0, len(initData))
	for _, ws := range initData {
		finalData = append(finalData, ToWorkspaceResponse(ws))
	}

	return res.SuccessResponse(finalData, "")
}

func (s *workspaceService) ShowById(
	ctx context.Context,
	workspaceId uuid.UUID,
	userId uuid.UUID,
) *res.Response[*workspace.WorkspaceResponse] {
	initData, err := s.client.Workspace.Query().
		Where(entWorkspace.IDEQ(workspaceId)).
		WithOwner().
		Order(entWorkspace.ByCreatedAt(sql.OrderDesc())).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return res.ErrorMessage[*workspace.WorkspaceResponse](
				"workspace is not found",
				http.StatusBadRequest,
			)
		}
		return res.ErrorMessage[*workspace.WorkspaceResponse]("failed to get workspace")
	}

	finalData := ToWorkspaceResponse(initData)
	return res.SuccessResponse(finalData, "")
}

func (s *workspaceService) Create(
	ctx context.Context,
	bodyData workspace.CreateWorkspaceDto,
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

	// TODO: imageURL or nil or default placeholder image
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

	builder := s.client.Workspace.Create().
		SetName(bodyData.Name).
		SetSlug(slug).
		SetOwnerID(userId).
		SetInviteCode(util.GenerateInviteCode(0)).
		SetNillableImageURL(imageURL)

	// TODO: return new workspaceId in response
	if _, err := builder.Save(ctx); err != nil {
		if file != nil && filePath != nil {
			_ = s.fileUploadService.DeleteImage(ctx, *filePath)
		}
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
