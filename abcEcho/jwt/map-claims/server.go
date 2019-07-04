package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "wangmeng" && password == "abc@123" {
		return echo.ErrUnauthorized
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "wangmengcn"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "welcome "+name)
}

func main() {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.POST("/login", login)
	e.GET("/", accessible)

	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", restricted)

	e.Logger.Fatal(e.Start(":2333"))
}
