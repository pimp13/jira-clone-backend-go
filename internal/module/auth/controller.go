package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pimp13/jira-clone-backend-go/internal/infrastructure/config"
	"github.com/pimp13/jira-clone-backend-go/pkg/res"
	"github.com/pimp13/jira-clone-backend-go/pkg/util"
)

type AuthController struct {
	authService    AuthService
	authMiddleware AuthMiddleware
}

func NewAuthController(authService AuthService, authMiddleware AuthMiddleware) *AuthController {
	return &AuthController{
		authService,
		authMiddleware,
	}
}

func (ac *AuthController) Routes(r *echo.Group) {
	r.POST("/auth/register", ac.handleRegister)
	r.POST("/auth/login", ac.handleLogin)
	r.GET("/auth/info", ac.handleUserInfo, ac.authMiddleware.SetAuthMiddleware)
	r.POST("/auth/logout", ac.handleLogout)
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

// @Tags		[Auth] {v1}
// @Accept		json
// @Produce	json
// @Param		request	body	LoginUserDto	true	"request body"
// @Router		/v1/auth/login [POST]
func (ac *AuthController) handleLogin(c echo.Context) error {
	var bodyData LoginUserDto
	validateResp, err := res.ValidateRequest(c, &bodyData)
	if err != nil {
		return res.JSON(c, res.ErrorMessage[struct{}]("failed to validate", http.StatusBadRequest))
	}
	if validateResp != nil {
		return c.JSON(http.StatusBadRequest, validateResp)
	}

	resp := ac.authService.Login(c.Request().Context(), &bodyData)

	if !resp.OK {
		return c.JSON(resp.StatusCode, resp)
	}

	c.SetCookie(&http.Cookie{
		Name:     config.Envs.App.AuthCookieName,
		Value:    resp.Data.AccessToken,
		Path:     "/",
		Domain:   config.Envs.App.Url,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})

	return c.JSON(resp.StatusCode, resp)
}

// @Tags		[Auth] {v1}
// @Accept		json
// @Produce	json
// @Router		/v1/auth/logout [POST]
func (ac *AuthController) handleLogout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     config.Envs.App.AuthCookieName,
		Value:    "",
		Path:     "/",
		Domain:   config.Envs.App.Url,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	return res.JSON(c, res.SuccessMessage("you are logouted!"))
}

// @Tags		[Auth] {v1}
// @Accept		json
// @Produce	json
// @Router		/v1/auth/info [GET]
// @Security	ApiKeyAuth
func (ac *AuthController) handleUserInfo(c echo.Context) error {
	user, err := util.GetCurrentUser(c)
	if err != nil {
		return res.JSON(c, res.ErrorMessage[struct{}]("you unauth", http.StatusUnauthorized))
	}

	return res.JSON(c, res.SuccessResponse(user, "you is logged!"))
}
