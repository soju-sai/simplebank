package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token has invalid claims: token is expired")
	ErrInvalidKey   = errors.New("token is unverifiable: error while executing keyfunc: key is invalid")
)

type MyClaims struct {
	Issuer    string           `json:"iss,omitempty"`
	Subject   string           `json:"sub,omitempty"`
	Audience  jwt.ClaimStrings `json:"aud,omitempty"`
	ExpiresAt *jwt.NumericDate `json:"exp,omitempty"`
	NotBefore *jwt.NumericDate `json:"nbf,omitempty"`
	IssuedAt  *jwt.NumericDate `json:"iat,omitempty"`
	ID        uuid.UUID        `json:"jti,omitempty"`
}

func NewClaims(issuer string, subject string, duration time.Duration) (jwt.Claims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	claims := &MyClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    issuer,
		Subject:   subject,
		ID:        tokenID,
		Audience:  []string{"aud"},
	}

	return claims, nil
}

func (claims *MyClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return claims.ExpiresAt, nil
}
func (claims *MyClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return claims.IssuedAt, nil
}
func (claims *MyClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return claims.NotBefore, nil
}
func (claims *MyClaims) GetIssuer() (string, error) {
	return claims.Issuer, nil
}
func (claims *MyClaims) GetSubject() (string, error) {
	return claims.Subject, nil
}
func (claims *MyClaims) GetAudience() (jwt.ClaimStrings, error) {
	return claims.Audience, nil
}
