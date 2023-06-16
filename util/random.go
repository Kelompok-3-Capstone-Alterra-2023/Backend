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

func GenerateRandomOrderNumber() string {
	currentTime := time.Now()
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(9999) + 1000
	orderNum := currentTime.Format("02012006") + "-" + fmt.Sprintf("%d", randomNum)

	return orderNum
}
