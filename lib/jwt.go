package lib

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type MyClaims struct {
	jwt.RegisteredClaims
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewClaims(id string, username, email, issuer string, duration time.Time) MyClaims {
	return MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: issuer,
			ExpiresAt: jwt.NewNumericDate(
				duration,
			),
		},
		ID:       id,
		Username: username,
		Email:    email,
	}
}

func GenerateToken(claims MyClaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString, secret string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New(ErrTokenExpired)
		}
		return nil, err
	}

	claims, ok := token.Claims.(*MyClaims)
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New(ErrTokenExpired)
	}
	if !ok {
		return nil, err
	}

	return claims, nil
}
