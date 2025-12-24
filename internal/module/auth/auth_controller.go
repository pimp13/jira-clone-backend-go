package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pimp13/jira-clone-backend-go/pkg/res"
)

type AuthController struct {
	authService AuthService
}

func NewAuthController(authService AuthService) *AuthController {
	return &AuthController{authService}
}

func (ac *AuthController) Routes(r *echo.Group) {
	r.POST("/auth/register", ac.handleRegister)
}

// @Tags		[Auth] {v1}
// @Accept		json
// @Produce	json
// @Param		request	body	RegisterUserDto	true	"request body"
// @Router		/v1/auth/register [POST]
func (ac *AuthController) handleRegister(c echo.Context) error {
	var bodyData RegisterUserDto
	if err := c.Bind(&bodyData); err != nil {
		return c.JSON(http.StatusBadRequest, res.Error[struct{}](err, http.StatusBadRequest))
	}

	file, err := c.FormFile("image")
	if err != nil {
		return res.JSON(c, res.ErrorResponse[struct{}]("failed to get image", err))
	}

	resp := ac.authService.Register(c.Request().Context(), &bodyData, file)

	return res.JSON(c, resp)
}
