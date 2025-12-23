package jwt

import (
	jwtpkg "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"uid"`
	jwtpkg.RegisteredClaims
}
