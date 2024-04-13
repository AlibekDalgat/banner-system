package service

import (
	"banner-system/internal/repository"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	signingKey = "adSj23&h#!kjWjqwnd@jnef7832N"
	tokenTTL   = 12 * time.Hour
)

type AuthService struct {
	repo repository.Authorization
}

type tokenClaims struct {
	jwt.MapClaims
	Login   string `json:"login"`
	IsAdmin bool   `json:"isAdmin"`
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo}
}

func (s *AuthService) GenerateToken(login, password string) (string, error) {
	user, isAdmin, err := s.repo.GetUser(login, password)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.MapClaims{
			"exp": time.Now().Add(tokenTTL).Unix(),
			"iat": time.Now().Unix(),
		},
		user.Login,
		isAdmin,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (bool, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Неверный метод авторизации")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return false, errors.New("клеймы токена неверного типа")
	}

	return claims.IsAdmin, nil
}
