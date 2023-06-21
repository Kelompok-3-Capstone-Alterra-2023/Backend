package util

import (
	"fmt"
	"math/rand"
	"time"
)

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GeneratePass(length int) (string, error) {
	const (
		letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digitBytes   = "0123456789"
		specialChars = "!@#$%^&*()"
	)

	charset := letterBytes + digitBytes + specialChars
	charsetLen := len(charset)

	randomBytes, err := generateRandomBytes(length)
	if err != nil {
		return "", err
	}

	password := make([]byte, length)
	for i := 0; i < length; i++ {
		password[i] = charset[int(randomBytes[i])%charsetLen]
	}

	return string(password), nil
}

func GenerateRandomString(name string) string {
	rand.Seed(time.Now().Unix())
	randomNum := rand.Intn(9999) + 1000
	imageName := name + "-" + fmt.Sprintf("%d", randomNum)

	return imageName
}

func GenerateRandomOrderNumber() string {
	currentTime := time.Now()
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(9999) + 1000
	orderNum := currentTime.Format("02012006") + "-" + fmt.Sprintf("%d", randomNum)

	return orderNum
}

func SplitBy(s string, separator rune) []string {
	var parts []string
	var currentPart string

	for _, char := range s {
		if char == separator {
			parts = append(parts, currentPart)
			currentPart = ""
		} else {
			currentPart += string(char)
		}
	}

	if len(currentPart) > 0 {
		parts = append(parts, currentPart)
	}

	return parts
}
