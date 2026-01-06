package workspace

import (
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/pimp13/jira-clone-backend-go/internal/module/auth"
	workspace "github.com/pimp13/jira-clone-backend-go/internal/module/workspace/dto"
	"github.com/pimp13/jira-clone-backend-go/pkg/res"
	"github.com/pimp13/jira-clone-backend-go/pkg/util"
)

type WorkspaceController struct {
	workspaceService WorkspaceService
	authMiddleware   auth.AuthMiddleware
}

func NewWorkspaceController(workspaceService WorkspaceService, authMiddleware auth.AuthMiddleware) *WorkspaceController {
	return &WorkspaceController{
		workspaceService,
		authMiddleware,
	}
}

func (ctrl *WorkspaceController) Routes(r *echo.Group) {
	router := r.Group("/workspace", ctrl.authMiddleware.SetAuthMiddleware)
	router.GET("", ctrl.index)
	router.POST("", ctrl.create)
	router.GET("/:id", ctrl.showById)
	router.PATCH("/:id", ctrl.update)
}

// @Tags		[Workspace] {v1}
// @Accept		json
// @Produce	json
// @Router		/v1/workspace [GET]
// @Security	ApiKeyAuth
func (ctrl *WorkspaceController) index(c echo.Context) error {
	user, err := util.GetCurrentUser(c)
	if err != nil {
		return res.JSON(c, res.ErrorMessage[struct{}]("you unauth", http.StatusUnauthorized))
	}

	resp := ctrl.workspaceService.Index(c.Request().Context(), user.ID)

	return c.JSON(resp.StatusCode, resp)
}

// @Tags		[Workspace] {v1}
// @Accept		json
// @Produce	json
// @Router		/v1/workspace/{id} [GET]
// @Param		id	path	string	true	"workspace id"
// @Security	ApiKeyAuth
func (ctrl *WorkspaceController) showById(c echo.Context) error {
	workspaceId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return res.JSON(c, res.ErrorMessage[struct{}](
			"bad request param",
			http.StatusUnauthorized,
		))
	}

	user, err := util.GetCurrentUser(c)
	if err != nil {
		return res.JSON(c, res.ErrorMessage[struct{}](
			"you unauth",
			http.StatusUnauthorized,
		))
	}

	resp := ctrl.workspaceService.ShowById(c.Request().Context(), workspaceId, user.ID)

	return c.JSON(resp.StatusCode, resp)
}

// @Tags		[Workspace] {v1}
// @Accept		json
// @Produce	json
// @Param		request	body	CreateWorkspaceDto	true	"request body"
// @Router		/v1/workspace [POST]
// @Security	ApiKeyAuth
func (ctrl *WorkspaceController) create(c echo.Context) error {
	user, err := util.GetCurrentUser(c)
	if err != nil {
		return res.JSON(c, res.ErrorMessage[struct{}]("you unauth", http.StatusUnauthorized))
	}

	var bodyData workspace.CreateWorkspaceDto
	validationErrs, err := res.ValidateRequest(c, &bodyData)
	if err != nil {
		return res.JSON(c, res.ErrorResponse[struct{}]("failed to parse body data", err))
	}
	if validationErrs != nil {
		return c.JSON(http.StatusBadRequest, validationErrs)
	}

	file, err := c.FormFile("image")
	var uploadedFile *multipart.FileHeader = nil
	if err == nil {
		uploadedFile = file
	} else if !errors.Is(err, http.ErrMissingFile) {
		return res.JSON(c, res.ErrorResponse[struct{}]("file is bad way", err))
	}

	resp := ctrl.workspaceService.Create(c.Request().Context(), bodyData, uploadedFile, user.ID)

	return c.JSON(resp.StatusCode, resp)
}

// @Tags		[Workspace] {v1}
// @Accept		json
// @Produce	json
// @Param		request	body	UpdateWorkspaceDto	true	"request body"
// @Param		id		path	string				true	"workspace id"
// @Router		/v1/workspace/{id} [PATCH]
// @Security	ApiKeyAuth
func (ctrl *WorkspaceController) update(c echo.Context) error {
	workspaceId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return res.JSON(c, res.ErrorMessage[struct{}](
			"bad request param",
			http.StatusUnauthorized,
		))
	}

	var bodyData workspace.UpdateWorkspaceDto
	validationErrs, err := res.ValidateRequest(c, &bodyData)
	if err != nil {
		return res.JSON(c, res.ErrorResponse[struct{}]("failed to parse body data", err))
	}
	if validationErrs != nil {
		return c.JSON(http.StatusBadRequest, validationErrs)
	}

	file, err := c.FormFile("image")
	var uploadedFile *multipart.FileHeader = nil
	if err == nil {
		uploadedFile = file
	} else if !errors.Is(err, http.ErrMissingFile) {
		return res.JSON(c, res.ErrorResponse[struct{}]("file is bad way", err))
	}

	resp := ctrl.workspaceService.Update(c.Request().Context(), bodyData, workspaceId, uploadedFile)

	return c.JSON(resp.StatusCode, resp)
}
