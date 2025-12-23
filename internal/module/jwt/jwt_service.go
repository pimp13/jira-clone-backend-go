package jwt

type JWTService interface {
	GenerateTokens(userId string)
}
