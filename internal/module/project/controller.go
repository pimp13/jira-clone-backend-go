package project

import (
	"github.com/labstack/echo/v4"
	"github.com/pimp13/jira-clone-backend-go/internal/module/auth"
	dto "github.com/pimp13/jira-clone-backend-go/internal/module/project/dto"
	"github.com/pimp13/jira-clone-backend-go/pkg/res"
)

type ProjectController struct {
	projectService ProjectService
	authMiddleware auth.AuthMiddleware
}

func NewProjectController(
	projectService ProjectService,
	authMiddleware auth.AuthMiddleware,
) *ProjectController {
	return &ProjectController{
		projectService,
		authMiddleware,
	}
}

func (ctrl *ProjectController) Routes(r *echo.Group) {
	router := r.Group("/project", ctrl.authMiddleware.SetAuthMiddleware)
	router.GET("", ctrl.index)
}

// @Tags		[Project] {v1}
// @Accept		json
// @Produce	json
// @Router		/v1/project [GET]
// @Security	ApiKeyAuth
func (ctrl *ProjectController) index(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"message": "Hello World",
		"ok":      true,
	})
}

//	@Tags		[Project] {v1}
//	@Accept		json
//	@Produce	json
//	@Param		request	body	CreateProjectDto	true	"request body"
//	@Router		/v1/project [POST]
//	@Security	ApiKeyAuth
func (ctrl *ProjectController) create(c echo.Context) error {
	var bodyData dto.CreateProjectDto
	validationErrs, err := res.ValidateRequest(c, &bodyData)
	if err != nil {
		return res.
	}
}