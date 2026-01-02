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
	dto "github.com/pimp13/jira-clone-backend-go/internal/module/workspace/dto"
	"github.com/pimp13/jira-clone-backend-go/pkg/res"
	"github.com/pimp13/jira-clone-backend-go/pkg/util"
)

type WorkspaceService interface {
	Index(
		ctx context.Context,
		userID uuid.UUID,
	) *res.Response[[]*dto.WorkspaceResponse]

	ShowById(
		ctx context.Context,
		workspaceId uuid.UUID,
		userID uuid.UUID,
	) *res.Response[*dto.WorkspaceResponse]

	Create(
		ctx context.Context,
		bodyData dto.CreateWorkspaceDto,
		file *multipart.FileHeader,
		userID uuid.UUID,
	) *res.Response[*dto.CreateWorkspaceResponse]

	Update(
		ctx context.Context,
		bodyData dto.UpdateWorkspaceDto,
		workspaceID uuid.UUID,
		file *multipart.FileHeader,
	) *res.Response[*dto.UpdateWorkspaceResponse]
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
	userID uuid.UUID,
) *res.Response[[]*dto.WorkspaceResponse] {
	initData, err := s.client.Workspace.Query().
		Where(entWorkspace.OwnerIDEQ(userID)).
		WithOwner().
		Order(entWorkspace.ByCreatedAt(sql.OrderDesc())).
		All(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return res.ErrorMessage[[]*dto.WorkspaceResponse](
				"workspace is not found",
				http.StatusBadRequest,
			)
		}
		return res.ErrorMessage[[]*dto.WorkspaceResponse]("failed to get workspace")
	}

	finalData := make([]*dto.WorkspaceResponse, 0, len(initData))
	for _, ws := range initData {
		finalData = append(finalData, ToWorkspaceResponse(ws))
	}

	return res.SuccessResponse(finalData, "")
}

func (s *workspaceService) ShowById(
	ctx context.Context,
	workspaceId uuid.UUID,
	userID uuid.UUID,
) *res.Response[*dto.WorkspaceResponse] {
	initData, err := s.client.Workspace.Query().
		Where(entWorkspace.IDEQ(workspaceId)).
		WithOwner().
		Order(entWorkspace.ByCreatedAt(sql.OrderDesc())).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return res.ErrorMessage[*dto.WorkspaceResponse](
				"workspace is not found",
				http.StatusBadRequest,
			)
		}
		return res.ErrorMessage[*dto.WorkspaceResponse]("failed to get workspace")
	}

	finalData := ToWorkspaceResponse(initData)
	return res.SuccessResponse(finalData, "")
}

func (s *workspaceService) Create(
	ctx context.Context,
	bodyData dto.CreateWorkspaceDto,
	file *multipart.FileHeader,
	userID uuid.UUID,
) *res.Response[*dto.CreateWorkspaceResponse] {
	var slug string
	var err error

	if bodyData.Slug != nil {
		slug, err = s.generateUniqueSlug(ctx, *bodyData.Slug, uuid.Nil)
	} else {
		slug, err = s.generateUniqueSlug(ctx, bodyData.Name, uuid.Nil)
	}
	if err != nil {
		return res.ErrorMessage[*dto.CreateWorkspaceResponse]("failed to generate slug")
	}

	// TODO: imageURL or nil or default placeholder image
	var imageURL *string = nil
	var filePath *string = nil

	if file != nil {
		uploadResult, err := s.fileUploadService.UploadImage(ctx, file)
		if err != nil {
			return res.ErrorResponse[*dto.CreateWorkspaceResponse]("failed to upload file", err)
		}
		imageURL = &uploadResult.URL
		filePath = &uploadResult.FilePath
	}

	builder := s.client.Workspace.Create().
		SetName(bodyData.Name).
		SetSlug(slug).
		SetOwnerID(userID).
		SetInviteCode(util.GenerateInviteCode(0)).
		SetNillableImageURL(imageURL).
		SetNillableImagePath(filePath)

	newWorkspace, err := builder.Save(ctx)
	if err != nil {
		if file != nil && filePath != nil {
			_ = s.fileUploadService.DeleteImage(ctx, *filePath)
		}
		return res.ErrorResponse[*dto.CreateWorkspaceResponse]("failed to save workspace", err)
	}

	return res.SuccessResponse(
		&dto.CreateWorkspaceResponse{
			ID: newWorkspace.ID,
		},
		"workspace is saved by successfully!",
		http.StatusCreated,
	)
}

func (s *workspaceService) Update(
	ctx context.Context,
	bodyData dto.UpdateWorkspaceDto,
	workspaceID uuid.UUID,
	file *multipart.FileHeader,
) *res.Response[*dto.UpdateWorkspaceResponse] {
	var err error

	workspace, err := s.client.Workspace.
		Query().
		Where(entWorkspace.IDEQ(workspaceID)).
		Select(
			entWorkspace.FieldID,
			entWorkspace.FieldImageURL,
			entWorkspace.FieldImagePath,
		).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return res.ErrorMessage[*dto.UpdateWorkspaceResponse]("workspace is not found")
		}
		return res.ErrorResponse[*dto.UpdateWorkspaceResponse]("failed to get workspace", err)
	}

	var slug string
	if bodyData.Slug != nil {
		slug, err = s.generateUniqueSlug(ctx, *bodyData.Slug, workspaceID)
	} else if bodyData.Name != nil {
		slug, err = s.generateUniqueSlug(ctx, *bodyData.Name, workspaceID)
	}
	if err != nil {
		return res.ErrorResponse[*dto.UpdateWorkspaceResponse]("failed to generate slug", err)
	}

	var imageURL *string = nil
	var filePath *string = nil
	var shouldDeleteOldImage bool
	if file != nil {
		uploadResult, err := s.fileUploadService.UploadImage(ctx, file)
		if err != nil {
			return res.ErrorResponse[*dto.UpdateWorkspaceResponse]("failed to upload file", err)
		}
		imageURL = &uploadResult.URL
		filePath = &uploadResult.FilePath
		shouldDeleteOldImage = true
	}

	builder := s.client.Workspace.
		UpdateOneID(workspaceID).
		SetNillableName(bodyData.Name).
		SetNillableImagePath(filePath).
		SetNillableImageURL(imageURL)

	if slug != "" {
		builder.SetSlug(slug)
	}

	updated, err := builder.Save(ctx)
	if err != nil {
		if file != nil && filePath != nil {
			_ = s.fileUploadService.DeleteImage(ctx, *filePath)
		}
		return res.ErrorResponse[*dto.UpdateWorkspaceResponse]("failed to update workspace", err)
	}

	if shouldDeleteOldImage && workspace.ImageURL != nil && workspace.ImagePath != nil {
		_ = s.fileUploadService.DeleteImage(ctx, *workspace.ImagePath)
	}

	return res.SuccessResponse(
		&dto.UpdateWorkspaceResponse{
			ID: updated.ID,
		},
		"workspace is updated successfully!",
	)
}

func (s *workspaceService) generateUniqueSlug(
	ctx context.Context,
	title string,
	excludeID uuid.UUID,
) (string, error) {
	baseSlug := util.Slugify(title)
	slug := baseSlug
	count := 1

	for {
		q := s.client.
			Workspace.
			Query().
			Where(entWorkspace.SlugEQ(slug))

		if excludeID != uuid.Nil {
			q = q.Where(entWorkspace.IDNEQ(excludeID))
		}

		exists, err := q.Exist(ctx)
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
