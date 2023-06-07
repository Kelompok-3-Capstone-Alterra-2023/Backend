package main

import (
	"capstone/config"

	"capstone/controller"
	"capstone/middleware"
	"capstone/model"
	"net/http"
  "capstone/route"

	"github.com/labstack/echo/v4"
)


	


func main() {

	config.Open()

	e:=route.New()
	


	e.Logger.Fatal(e.Start(":8080"))
}
