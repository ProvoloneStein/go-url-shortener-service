package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v4"
)

const (
	tokenTTL = 12 * time.Hour
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

func (s *Service) GenerateToken(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", defaultServiceErrWrapper(ctx.Err())
	default:
	}
	userID := uuid.New().String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
		userID,
	})

	signedString, err := token.SignedString([]byte(s.cfg.SigningKey))
	if err != nil {
		return "", defaultServiceErrWrapper(err)
	}

	return signedString, nil
}

func (s *Service) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("service: invalid signing method")
		}
		return []byte(s.cfg.SigningKey), nil
	})
	if err != nil {
		return "", defaultServiceErrWrapper(err)
	}
	if !token.Valid {
		return "", errors.New("service: token is not valid")
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("service: token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}
