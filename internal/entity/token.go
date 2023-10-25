package entity

import (
	"fmt"
	"go-test-grpc-http/cmd/go-test-grpc-http/config"
	"time"

	"github.com/golang-jwt/jwt"
)

type Token struct {
	Token *jwt.Token
}

func (t *Token) String() (string, error) {
	cfg, err := config.GetAppConfig()
	if err != nil {
		return "", fmt.Errorf("can't generate token: %w", err)
	}

	return t.Token.SignedString([]byte(cfg.ApiKey))
}

func GenerateToken(id *UserID) *Token {
	return &Token{
		Token: jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Id:        id.String(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Subject:   "auth",
		}),
	}
}

func ParseToken(tokenString string) (*UserID, error) {
	cfg, err := config.GetAppConfig()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.ApiKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	var id UserID
	err = id.FromString(claims.Id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}
