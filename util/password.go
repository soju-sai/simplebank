package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hashing password failed: %w", err)
	}

	return string(hp), nil
}

func CheckHashedPassword(password string, hp string) error {
	return bcrypt.CompareHashAndPassword([]byte(hp), []byte(password))
}
