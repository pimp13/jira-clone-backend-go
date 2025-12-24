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
	validateErrs, err := res.ValidateRequest(c, &bodyData)
	if err != nil {
		return res.JSON(c, res.ErrorMessage[struct{}]("failed to validate", http.StatusBadRequest))
	}
	if validateErrs != nil {
		return c.JSON(http.StatusBadRequest, validateErrs)
	}

	resp := ac.authService.Register(c.Request().Context(), &bodyData)

	return res.JSON(c, resp)
}
