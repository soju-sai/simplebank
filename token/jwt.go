package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secreteKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size, must not less than %d", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	claims, err := NewClaims(username, "local", duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(maker.secreteKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*MyClaims, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrInvalidKey
		}
		return []byte(maker.secreteKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &MyClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}
	myClaim, ok := jwtToken.Claims.(*MyClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return myClaim, nil
}
