package main

import (
	"github.com/GeertJohan/go.rice"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	assetsHandler := http.FileServer(rice.MustFindBox("app").HTTPBox())

	e.GET("/", echo.WrapHandler(assetsHandler))
	e.GET("/qrcode/*", echo.WrapHandler(http.StripPrefix("/qrcode/", assetsHandler)))

	e.Logger.Fatal(e.Start(":5678"))
}
