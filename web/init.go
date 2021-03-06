package web

import (
	"fmt"
	"games-backend/app"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/xdean/goex/xecho"
	"net/http"

	_ "games-backend/web/game"
)

func init() {
	app.App.RegisterInitFunc(func() {
		e := echo.New()

		e.Validator = xecho.NewValidator()

		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.Use(xecho.BreakErrorRecover())

		e.GET("/ping", func(c echo.Context) error {
			return c.JSON(http.StatusOK, "pong")
		})

		app.Context.Echo = e
	})
	app.App.RegisterRun(func() {
		logrus.Fatal(app.Context.Echo.Start(fmt.Sprintf(":%d", app.Config.Port)))
	})
}
