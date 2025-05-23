package usecase

import (
	"context"
	"errors"

	"monitoring/internal/domain/entity"
	iface "monitoring/internal/domain/interface"
	"monitoring/pkg/db"
	"monitoring/pkg/jwt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	userRepo iface.UserRepository
	cache0   *db.Client
}

func NewAuthUsecase(userRepo iface.UserRepository, cache *db.Client) iface.AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
		cache0:   cache,
	}
}

func (u *authUsecase) Register(ctx context.Context, req *entity.RegisterRequest) error {
	// cek username
	_, err := u.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return errors.New("username already exists")
	}

	// cek email
	_, err = u.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return errors.New("email already exists")
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	return u.userRepo.Create(ctx, &user)
}

func (u *authUsecase) Login(ctx context.Context, username, password string) (string, string, error) {
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", "", errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid username or password")
	}

	// generate access & refresh token
	accessToken, err := jwt.GenerateToken(user.ID.String(), time.Minute*15)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.GenerateToken(user.ID.String(), time.Hour*24*7)
	if err != nil {
		return "", "", err
	}

	key := "refresh_token:" + user.ID.String()
	err = u.cache0.Set(ctx, key, refreshToken, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (u *authUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := jwt.ValidateToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	userID := claims.UserID
	key := "refresh_token:" + userID

	// Ambil refresh token dari Redis db 1
	storedToken, err := u.cache0.Get(ctx, key)
	if err != nil {
		return "", "", errors.New("refresh token not found or expired")
	}

	if storedToken != refreshToken {
		return "", "", errors.New("refresh token does not match")
	}

	newAccessToken, err := jwt.GenerateToken(userID, time.Minute*15)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := jwt.GenerateToken(userID, time.Hour*24*7)
	if err != nil {
		return "", "", err
	}

	// Update refresh token di Redis
	err = u.cache0.Set(ctx, key, newRefreshToken, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}
