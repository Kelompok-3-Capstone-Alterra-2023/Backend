package main

import (
	"capstone/config"
	"log"
	"net/http"

	"capstone/route"
)

func main() {

	config.Open()

	e := route.New()

	if err := e.StartTLS(":8080", "/app/host.cert", "/app/host.key"); err != http.ErrServerClosed {
		log.Fatal(err)
	  }
}
