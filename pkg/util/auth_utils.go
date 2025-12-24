package util

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/pimp13/jira-clone-backend-go/internal/infrastructure/config"
	"github.com/pimp13/jira-clone-backend-go/internal/module/jwt"
)

type UserAuth string

func (s UserAuth) String() string {
	return string(s)
}

const UserAuthKey UserAuth = "user_auth"

func ExtractTokenFromRequest(c echo.Context) (string, error) {

	cookie, err := c.Cookie(config.Envs.App.AuthCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil

	// authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
	// if strings.HasPrefix(authHeader, "Bearer ") {
	// 	return strings.TrimPrefix(authHeader, "Bearer "), nil
	// } else if authHeader != "" {
	// 	return authHeader, nil
	// }
	// return "", errors.New(echo.ErrUnauthorized.Error())
}

func GetCurrentUser(c echo.Context) (*jwt.UserInfo, error) {
	user, exists := c.Get(UserAuthKey.String()).(*jwt.UserInfo)
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}
