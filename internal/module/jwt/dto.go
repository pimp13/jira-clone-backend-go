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
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Email holds the value of the "email" field.
	Email string `json:"email,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Password holds the value of the "password" field.
	Password string `json:"-"`
	// IsActive holds the value of the "is_active" field.
	IsActive *bool `json:"isActive,omitempty"`
	// AvatarURL holds the value of the "avatar_url" field.
	AvatarURL *string `json:"avatarUrl,omitempty"`
	// Role holds the value of the "role" field.
	Role user.Role `json:"role,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}
