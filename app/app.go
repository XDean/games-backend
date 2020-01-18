package app

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/xdean/goex/xapp"
)

var (
	App = xapp.App{
		Config: xapp.ConfigRegistry{},
	}

	Config = struct {
	}{}

	Context = struct {
		Echo *echo.Echo
	}{}
)

func init() {
	App.Config.Register(&Config, "app")
}

func Debug() {
	logrus.SetLevel(logrus.DebugLevel)
}
