package main

import (
	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	e.GET("/", func(context echo.Context) error {
		return context.String(http.StatusOK, "23333")
	})

	e.Server.Addr = ":2333"
	e.Logger.Fatal(gracehttp.Serve(e.Server))
}
