package main

import (
	"capstone/config"
	"capstone/route"
)

func main() {
	config.InitDB()
	e := route.New()
	e.Logger.Fatal(e.Start(":8080"))
}
