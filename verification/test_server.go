package main

import (
	"log"
	"oauth2-provider/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	htmlHandler := handlers.NewHTMLHandler()
	e.GET("/", htmlHandler.Welcome)

	log.Println("Starting verification server on :8001")
	if err := e.Start(":8001"); err != nil {
		log.Fatal(err)
	}
}
