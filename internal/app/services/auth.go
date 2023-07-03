package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/ProvoloneStein/go-url-shortener-service/internal/app/repositories"

	"github.com/golang-jwt/jwt/v4"
)

const (
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
}

func (s *Service) GenerateToken(ctx context.Context) (string, error) {
	var userID string
	for {
		userID = uuid.New().String()
		if err := s.repo.ValidateUniqueUser(ctx, userID); err != nil {
			if !errors.Is(err, repositories.ErrUserExists) {
				return "", defaultServiceErrWrapper(err)
			}
		} else {
			break
		}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
		userID,
	})

	signedString, err := token.SignedString([]byte(signingKey))
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
		return []byte(signingKey), nil
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
