package main

import (
	"github.com/labstack/echo/v4"
	"github.com/tylerb/graceful"
	"net/http"
	"time"
)

func main() {
	e := echo.New()
	e.GET("/", func(context echo.Context) error {
		return context.String(http.StatusOK, "2333")
	})
	e.Server.Addr = ":2333"

	graceful.ListenAndServe(e.Server, 5*time.Second)
}
