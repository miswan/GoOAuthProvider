package main

import (
	"log"
	"oauth2-provider/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	h := handlers.NewHTMLHandler()

	e.GET("/", h.Index)
	e.GET("/login", h.Login)

    e.POST("/login", func(c echo.Context) error {
        return c.String(200, "Login POST mock")
    })

	log.Println("Starting verification server on :8081")
	e.Start(":8081")
}
