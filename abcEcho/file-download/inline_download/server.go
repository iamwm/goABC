package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(context echo.Context) error {
		return context.File("./index.html")
	})
	e.GET("/inline", func(context echo.Context) error {
		return context.File("./test.png")
	})

	e.Logger.Fatal(e.Start(":2333"))
}
