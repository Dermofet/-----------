package entity

import (
	"fmt"
	"music-backend-test/cmd/music-backend-test/config"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Token struct {
	Token *jwt.Token
}

type TokenShow struct {
	Token string `json:"token"`
}

func (t *Token) String() (string, error) {
	cfg, err := config.GetAppConfig()
	if err != nil {
		return "", fmt.Errorf("can't generate token: %w", err)
	}

	return t.Token.SignedString([]byte(cfg.ApiKey))
}

func GenerateToken(id uuid.UUID) *Token {
	return &Token{
		Token: jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Id:        id.String(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Subject:   "auth",
		}),
	}
}

func ParseToken(tokenString string) (uuid.UUID, error) {
	cfg, err := config.GetAppConfig()
	if err != nil {
		return uuid.Nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.ApiKey), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	id, err := uuid.Parse(claims.Id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
