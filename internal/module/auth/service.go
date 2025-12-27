package auth

import (
	"context"
	"net/http"

	"github.com/pimp13/jira-clone-backend-go/ent"
	"github.com/pimp13/jira-clone-backend-go/ent/user"
	"github.com/pimp13/jira-clone-backend-go/internal/module/jwt"
	"github.com/pimp13/jira-clone-backend-go/pkg/logger"
	"github.com/pimp13/jira-clone-backend-go/pkg/res"
	"github.com/pimp13/jira-clone-backend-go/pkg/util"
)

type AuthService interface {
	Register(ctx context.Context, bodyData *RegisterUserDto) *res.Response[struct{}]

	Login(ctx context.Context, bodyData *LoginUserDto) *res.Response[LoginResponse]
}

type authService struct {
	client     *ent.Client
	jwtService jwt.JWTService
	logger     logger.Logger
}

func NewAuthService(client *ent.Client, jwtService jwt.JWTService, logger logger.Logger) AuthService {
	return &authService{
		client:     client,
		jwtService: jwtService,
		logger:     logger,
	}
}

func (as *authService) Register(
	ctx context.Context,
	bodyData *RegisterUserDto,
) *res.Response[struct{}] {
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

	as.logger.Info().Str("email", bodyData.Email).Msg("A new user registered!")

	return res.SuccessMessage("register user is successfully!")
}

func (as *authService) Login(
	ctx context.Context,
	bodyData *LoginUserDto,
) *res.Response[LoginResponse] {
	user, err := as.findUserByEmail(ctx, bodyData.Email)
	if err != nil {
		return res.ErrorMessage[LoginResponse]("email or password is invalid", http.StatusBadRequest)
	}

	if !util.CheckPasswordHash(bodyData.Password, user.Password) {
		return res.ErrorMessage[LoginResponse]("password or email is invalid", http.StatusBadRequest)
	}

	accessToken, refreshToken, refreshJti, err := as.jwtService.GenerateTokens(user.ID)
	if err != nil {
		return res.ErrorMessage[LoginResponse]("failed to generate auth token")
	}

	as.logger.Info().Msgf("User by email %s is login...", bodyData.Email)

	return res.SuccessResponse(LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		RefreshJti:   refreshJti,
	}, "user logged in successfully")
}

func (as *authService) userExistsByEmail(ctx context.Context, email string) (bool, error) {
	return as.client.User.Query().Where(user.EmailEQ(email)).Exist(ctx)
}

func (as *authService) findUserByEmail(ctx context.Context, email string) (*ent.User, error) {
	return as.client.User.Query().Where(user.EmailEQ(email)).First(ctx)
}
