package auth

import (
	"context"

	"github.com/pimp13/jira-clone-backend-go/ent"
)

type AuthService interface {
	Register(ctx context.Context)
}

type authService struct {
	client *ent.Client
}

func NewAuthService(client *ent.Client) *AuthService {
	return &authService{
		client,
	}
}
