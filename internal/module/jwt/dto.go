package jwt

import (
	"time"

	jwtpkg "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pimp13/jira-clone-backend-go/ent/user"
)

type Claims struct {
	UserID uuid.UUID `json:"uid"`
	jwtpkg.RegisteredClaims
}

type UserInfo struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	Name      string    `json:"name,omitempty"`
	Password  string    `json:"-"`
	IsActive  *bool     `json:"isActive,omitempty"`
	AvatarURL *string   `json:"avatarUrl,omitempty"`
	Role      user.Role `json:"role,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}
