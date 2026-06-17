package provider

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JsonWebToken struct {
	secret string
}

func NewJsonWebTokenProvider(secret string) *JsonWebToken {
	return &JsonWebToken{
		secret: secret,
	}
}

func (j *JsonWebToken) ValidateToken(ctx context.Context, token string) (int64, error) {
	var (
		userID int64
	)

	result, err := jwt.Parse(
		token,
		func(token *jwt.Token) (any, error) {
			return []byte(j.secret), nil
		},
	)

	if err != nil || !result.Valid {
		return userID, errors.New("invalid jwt token")
	}

	claims, ok := result.Claims.(jwt.MapClaims)
	if !ok {
		return userID, errors.New("failed to parse token claims")
	}

	sub, ok := claims["sub"]
	if !ok {
		return userID, errors.New("user_id not found in token")
	}

	value, ok := sub.(float64)
	if !ok {
		return userID, errors.New("invalid user_id data type")
	}

	userID = int64(value)

	return userID, nil
}

func (j *JsonWebToken) GenerateToken(userID int64) (string, error) {

	now := time.Now()

	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(j.secret),
	)
}
