package main

import (
	"github.com/labstack/echo/v4"
	"oauth2-provider/handlers"
)

func main() {
	e := echo.New()
	htmlHandler := handlers.NewHTMLHandler()
	e.GET("/", htmlHandler.Index)
	e.Start(":8081")
}
