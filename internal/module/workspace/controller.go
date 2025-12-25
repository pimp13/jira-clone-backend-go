package workspace

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pimp13/jira-clone-backend-go/internal/module/auth"
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
// @Param		request	body	CreateWorkspaceDto	true	"request body"
// @Router		/v1/workspace [POST]
// @Security	ApiKeyAuth
func (ctrl *WorkspaceController) create(c echo.Context) error {
	user, err := util.GetCurrentUser(c)
	if err != nil {
		return res.JSON(c, res.ErrorMessage[struct{}]("you unauth", http.StatusUnauthorized))
	}

	var bodyData CreateWorkspaceDto
	validationErrs, err := res.ValidateRequest(c, &bodyData)
	if err != nil {
		return res.JSON(c, res.ErrorResponse[struct{}]("failed to parse body data", err))
	}
	if validationErrs != nil {
		return c.JSON(http.StatusBadRequest, validationErrs)
	}
	log.Println("BODY DATA =>", bodyData)

	file, err := c.FormFile("image")
	if err != nil {
		return res.JSON(c, res.ErrorResponse[struct{}]("file is bad way", err))
	}

	resp := ctrl.workspaceService.Create(c.Request().Context(), bodyData, file, user.ID)

	return c.JSON(resp.StatusCode, resp)
}
