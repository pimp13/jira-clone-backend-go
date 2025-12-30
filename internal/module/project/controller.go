package project

import "github.com/labstack/echo/v4"

type ProjectController struct {
	projectService ProjectService
}

func NewProjectController(projectService ProjectService) *ProjectController {
	return &ProjectController{
		projectService,
	}
}

func (ctrl *ProjectController) Routes(r *echo.Group) {
	r.GET("/project", ctrl.index)
}

func (ctrl *ProjectController) index(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"message": "Hello World",
		"ok":      true,
	})
}
