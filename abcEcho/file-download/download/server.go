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
	e.GET("/file", func(context echo.Context) error {
		return context.Inline("./test.png","test.png")
	})

	e.Logger.Fatal(e.Start(":2333"))
}
