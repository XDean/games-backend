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
		Port int
	}{
		Port: 11071,
	}

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

func RegisterWeb(f func(c *echo.Echo)) {
	App.RegisterInit(xapp.InitTask{
		Ready: func() bool {
			return Context.Echo != nil
		},
		Init: func() {
			f(Context.Echo)
		},
	})
}
