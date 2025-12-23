package auth

import (
	"context"
	"net/http"

	"github.com/pimp13/jira-clone-backend-go/ent"
	"github.com/pimp13/jira-clone-backend-go/ent/user"
	"github.com/pimp13/jira-clone-backend-go/pkg/res"
	"github.com/pimp13/jira-clone-backend-go/pkg/util"
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

	hashed, err := util.HashPassword(bodyData.Password)
	if err != nil {
		return res.ErrorMessage[struct{}]("error in hash password")
	}

	if _, err = as.client.User.Create().
		SetEmail(bodyData.Email).
		SetName(bodyData.Name).
		SetPassword(hashed).
		Save(ctx); err != nil {
		return res.ErrorMessage[struct{}]("error in register user")
	}
	return res.SuccessMessage("register user is successfully!")
}

func (as *authService) userExistsByEmail(ctx context.Context, email string) (bool, error) {
	return as.client.User.Query().Where(user.EmailEQ(email)).Exist(ctx)
}
