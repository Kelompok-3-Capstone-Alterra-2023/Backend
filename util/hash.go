package util

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}

func CompareHashAndPassword(hashedPass, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(pass))
	if err != nil {
		return err
	}

	return nil
}
