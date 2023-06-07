package main

import (
	"capstone/config"
	"capstone/route"
)

func main() {
	config.Open()

	e:=route.New()
	
	e.Logger.Fatal(e.Start(":8080"))
}
