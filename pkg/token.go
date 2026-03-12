package util

import (
	"arena-ban/internal/domain"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)
var JwtUtil *TokenUtil

type TokenUtil struct {
	SecretKey      string
	ExpiryDuration time.Duration
}

func InitJwt() {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		panic("JWT_SECRET not found")
	}

	JwtUtil = &TokenUtil{
		SecretKey:      secretKey,
		ExpiryDuration: time.Hour * 24,
	}
}

func (t *TokenUtil) CreateToken(user *domain.User) (string, error) {
	claims := &domain.TokenClaims{
		Email: user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.ExpiryDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.SecretKey))
}

func (t *TokenUtil) ParseToken(tokenString string) (*domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid Signing Method: %v", token.Header["alg"])
		}
		return []byte(t.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("Expired Token")
		}
		return nil, fmt.Errorf("gagal memproses token: %v", err)
	}

	if claims, ok := token.Claims.(*domain.TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Invalid Token")
}