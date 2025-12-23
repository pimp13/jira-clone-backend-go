package auth

import (
	"context"
	"net/http"

	"github.com/pimp13/jira-clone-backend-go/ent"
	"github.com/pimp13/jira-clone-backend-go/ent/user"
	"github.com/pimp13/jira-clone-backend-go/pkg/res"
)

type AuthService interface {
	Register(ctx context.Context, bodyData *RegisterUserDto) *res.Response[struct{}]
}

type authService struct {
	client *ent.Client
}

func NewAuthService(client *ent.Client) AuthService {
	return &authService{
		client,
	}
}

func (as *authService) Register(ctx context.Context, bodyData *RegisterUserDto) *res.Response[struct{}] {
	existsByEmail, err := as.userExistsByEmail(ctx, bodyData.Email)
	if err != nil {
		return res.ErrorMessage[struct{}]("error in find user by email", http.StatusInternalServerError)
	}
	if existsByEmail {
		return res.ErrorMessage[struct{}]("email is exists", http.StatusBadRequest)
	}
	return res.SuccessMessage("ok")
}

func (as *authService) userExistsByEmail(ctx context.Context, email string) (bool, error) {
	return as.client.User.Query().Where(user.EmailEQ(email)).Exist(ctx)
}
