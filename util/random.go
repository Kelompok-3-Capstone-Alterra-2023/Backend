package util

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateRandomString(name string) string {
	rand.Seed(time.Now().Unix())
	randomNum := rand.Intn(9999) + 1000
	imageName := name + "-" + fmt.Sprintf("%d", randomNum)

	return imageName
}
