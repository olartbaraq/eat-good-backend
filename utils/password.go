package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "Password must not be empty", errors.New("password must not be empty")

	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "could not generate hash password", err
	}

	hashedPassword := string(hash)
	return hashedPassword, nil
}

func VerifyPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err
}
