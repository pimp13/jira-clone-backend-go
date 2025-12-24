package auth

import (
	"github.com/pimp13/jira-clone-backend-go/ent"
)

// Requests
type RegisterUserDto struct {
	Name     string `json:"name" form:"name" validate:"required,min=3"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
}

type LoginUserDto struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
}

// Responses
type LoginResponse struct {
	User         *ent.User `json:"user"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	RefreshJti   string    `json:"refreshJti"`
}
