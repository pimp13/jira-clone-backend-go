package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	authService *AuthService
}

func NewAuthController(authService *AuthService) *AuthController {
	return &AuthController{authService}
}

func (ac *AuthController) Routes(r *echo.Group) {
	r.GET("/auth", ac.GetForTest)
}

// @Tags		users
// @Accept		json
// @Produce	json
// @Router		/v1/auth [GET]
func (ac *AuthController) GetForTest(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"ok":   true,
		"data": ac.authService.GetForTestInService("Pouya"),
	})
}
