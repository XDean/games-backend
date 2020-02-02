package test

import (
	"games-backend/app"
	"github.com/labstack/echo/v4"
	"net/http"
)

func init() {
	app.RegisterWeb(func(c *echo.Echo) {
		c.GET("/ping", func(c echo.Context) error {
			return c.JSON(http.StatusOK, "pong")
		})
	})
}
