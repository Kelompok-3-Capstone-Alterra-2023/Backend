package email

import (
	"math/rand"
	"time"

	sendinblue "github.com/CyCoreSystems/sendinblue"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	otpBytes := make([]byte, 8)
	for i := range otpBytes {
		otpBytes[i] = chars[rand.Intn(len(chars))]
	}
	return string(otpBytes)
}

func SendEmail(toName, toEmail, subject, content string) error {
	sender := sendinblue.Address{
		Name:  "capstone project",
		Email: "test@example.com",
	}
	recipient := sendinblue.Address{
		Name:  toName,
		Email: toEmail,
	}
	message := sendinblue.Message{
		Sender:      &sender,
		To:          []*sendinblue.Address{&recipient},
		Subject:     subject,
		TextContent: content,
	}
	return message.Send("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")
}
