package workspace

import (
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
	router.GET("/", ctrl.index)
	router.POST("/", ctrl.create)
}

// @Tags			[Workspace] {v1}
// @Accept		json
// @Produce		json
// @Router		/v1/workspace [GET]
// @Security	ApiKeyAuth
func (ctrl *WorkspaceController) index(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"message": "Hello World",
		"ok":      true,
	})
}

// @Tags			[Workspace] {v1}
// @Accept		json
// @Produce		json
// @Param			request	body	CreateWorkspaceDto	true	"request body"
// @Router		/v1/workspace [POST]
// @Security	ApiKeyAuth
func (ctrl *WorkspaceController) create(c echo.Context) error {
	user, err := util.GetCurrentUser(c)
	if err != nil {
		return res.JSON(c, res.ErrorMessage[struct{}]("you unauth", http.StatusUnauthorized))
	}

	var bodyData CreateWorkspaceDto
	if err := c.Bind(&bodyData); err != nil {
		return res.JSON(c, res.ErrorResponse[struct{}]("bad request body", err))
	}

	file, err := c.FormFile("image")
	if err != nil {
		return res.JSON(c, res.ErrorResponse[struct{}]("file is bad way", err))
	}

	resp := ctrl.workspaceService.Create(c.Request().Context(), bodyData, file, user.ID)

	return c.JSON(resp.StatusCode, resp)
}
