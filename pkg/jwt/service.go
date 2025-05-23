package jwt

import (
	"time"
)

type JwtService interface {
	GenerateToken(userID string, expireDuration time.Duration) (string, error)
	ValidateToken(tokenStr string) (*Claims, error)
}

type jwtService struct{}

func NewJwtService() JwtService {
	return &jwtService{}
}

func (j *jwtService) GenerateToken(userID string, expireDuration time.Duration) (string, error) {
	return GenerateToken(userID, expireDuration)
}

func (j *jwtService) ValidateToken(tokenStr string) (*Claims, error) {
	return ValidateToken(tokenStr)
}
