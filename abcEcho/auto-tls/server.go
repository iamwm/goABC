package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
)

func main() {
	e := echo.New()

	e.AutoTLSManager.Cache = autocert.DirCache("/tmp/www/.cache")
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, `
				<h1>Echo TLS</h1>`)
	})
	e.Logger.Fatal(e.StartAutoTLS(":443"))
}
