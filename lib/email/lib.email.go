package email

import (
	"fmt"
	"math/rand"
	"time"

	sendinblue "github.com/CyCoreSystems/sendinblue"
)

func GenerateOTP() string{
	rand.Seed(time.Now().UnixNano())
	chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	otpBytes := make([]byte, 8)
	for i := range otpBytes {
		otpBytes[i] = chars[rand.Intn(len(chars))]
	}
	return string(otpBytes)
}

func SendEmail(toName, toEmail , otp string) error {
	sender := sendinblue.Address{
		Name:  "Miniproject",
		Email: "test@example.com",
	}
	recipient := sendinblue.Address{
		Name:  toName,
		Email: toEmail,
	}
    message := sendinblue.Message{
		Sender: &sender,
		To:     []*sendinblue.Address{&recipient},
		Subject: "Account Creation",
		TextContent:    fmt.Sprintf("OTP: %s", otp),
	}
	return message.Send("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")
}
