package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pimp13/jira-clone-backend-go/internal/module/jwt"
	"github.com/pimp13/jira-clone-backend-go/pkg/res"
	"github.com/pimp13/jira-clone-backend-go/pkg/util"
)

type AuthMiddleware interface {
	SetAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type authMiddleware struct {
	jwtSvc jwt.JWTService
}

func NewAuthMiddleware(jwtSvc jwt.JWTService) AuthMiddleware {
	return &authMiddleware{
		jwtSvc: jwtSvc,
	}
}

func (m *authMiddleware) SetAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenStr, err := util.ExtractTokenFromRequest(c)
		if err != nil {
			return res.JSON(c, res.ErrorResponse[struct{}]("can't found token", err, http.StatusUnauthorized))
		}

		claims, err := m.jwtSvc.ParseAccessToken(tokenStr)
		if err != nil {
			return res.JSON(c, res.ErrorResponse[struct{}]("can't found claims", err, http.StatusUnauthorized))
		}

		user, err := m.jwtSvc.GetUserByID(c.Request().Context(), claims.UserID)
		if err != nil {
			return res.JSON(c, res.ErrorResponse[struct{}]("can't found user", err, http.StatusUnauthorized))
		}

		c.Set(util.UserAuthKey.String(), user)

		return next(c)
	}
}
