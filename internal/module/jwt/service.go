package jwt

import (
	"context"
	"fmt"
	"time"

	jwtpkg "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/pimp13/jira-clone-backend-go/ent"
	"github.com/pimp13/jira-clone-backend-go/ent/user"
	"github.com/pimp13/jira-clone-backend-go/internal/infrastructure/config"
)

type JWTService interface {
	GenerateTokens(userId uuid.UUID) (accessToken string, refreshToken string, refreshJti string, err error)

	ParseAccessToken(tokenStr string) (*Claims, error)

	GetUserByID(ctx context.Context, userID uuid.UUID) (*UserInfo, error)
}

type jwtService struct {
	config     *config.Config
	client     *ent.Client
	secretKey  []byte
	accessTTL  uint64
	refreshTTL uint64
}

func NewJWTService(client *ent.Client, config *config.Config) JWTService {
	js := &jwtService{
		config: config,
		client: client,
	}
	js.secretKey = []byte(config.App.SecretKey)
	js.accessTTL = 258000000
	js.refreshTTL = 2540000000
	return js
}

func (s *jwtService) GenerateTokens(userID uuid.UUID) (
	accessToken string, refreshToken string, refreshJti string, err error,
) {
	now := time.Now()

	atClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwtpkg.RegisteredClaims{
			IssuedAt: jwtpkg.NewNumericDate(now),
			// TODO: get from accessTTL
			ExpiresAt: jwtpkg.NewNumericDate(now.Add(time.Duration(time.Hour * 24))),
			ID:        userID.String(),
			Subject:   userID.String(),
		},
	}

	at := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, atClaims)
	accessToken, err = at.SignedString(s.secretKey)
	if err != nil {
		return
	}

	// refresh token (we put jti)
	rtClaims := &Claims{
		UserID: userID,
		RegisteredClaims: jwtpkg.RegisteredClaims{
			IssuedAt: jwtpkg.NewNumericDate(now),
			// TODO: get from refreshTTL
			ExpiresAt: jwtpkg.NewNumericDate(now.Add(time.Duration(time.Hour * 24 * 7))),
			ID:        userID.String(),
			Subject:   userID.String(),
		},
	}
	rt := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, rtClaims)
	refreshToken, err = rt.SignedString(s.secretKey)
	if err != nil {
		return
	}

	return
}

func (s *jwtService) ParseAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwtpkg.ParseWithClaims(tokenStr, &Claims{}, func(t *jwtpkg.Token) (any, error) {
		if _, ok := t.Method.(*jwtpkg.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secretKey, nil
	}, jwtpkg.WithLeeway(time.Second), jwtpkg.WithValidMethods([]string{jwtpkg.SigningMethodHS256.Name}))

	if err != nil {
		return nil, fmt.Errorf("can't parse jwt token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwtpkg.ErrTokenInvalidClaims
}

func (s *jwtService) GetUserByID(ctx context.Context, userID uuid.UUID) (*UserInfo, error) {
	user, err := s.client.User.Query().
		Where(user.IDEQ(userID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Password:  user.Password,
		IsActive:  user.IsActive,
		AvatarURL: user.AvatarURL,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
