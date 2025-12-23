package auth

import (
	"fmt"

	"github.com/pimp13/jira-clone-backend-go/ent"
)

type AuthService struct {
	client *ent.Client
}

func NewAuthService(client *ent.Client) *AuthService {
	return &AuthService{
		client,
	}
}

func (as *AuthService) GetForTestInService(name string) string {
	return fmt.Sprintf("Hello, %s", name)
}
