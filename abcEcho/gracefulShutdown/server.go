package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.GET("/", func(context echo.Context) error {
		time.Sleep(5 * time.Second)
		return context.JSON(http.StatusOK, "OK")
	})

	go func() {
		if err := e.Start(":2333"); err != nil {
			e.Logger.Info("Shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
