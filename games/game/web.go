package game

import (
	"fmt"
	"games-backend/app"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/xdean/goex/xecho"
	"net/http"
)

func init() {
	app.RegisterWeb(func(c *echo.Echo) {
		c.POST("/api/game/:game", createGame)
		c.GET("/socket/game/:game/:id", gameSocket)
	})
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func gameSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	xecho.MustNoError(err)

	user := c.QueryParam("user")
	gameName := c.Param("game")
	id := IntParam(c, "id")

	host := GetHost(gameName, id)
	if host == nil {
		_ = ws.WriteJSON(TopicEvent{
			Topic:   "error",
			Payload: "房间不存在",
		})
		_ = ws.Close()
		return nil
	}

	client := NewClient(user, host, ws)
	client.Start()
	return nil
}

func createGame(c echo.Context) error {
	gameName := c.Param("game")

	factory := GetFactory(gameName)
	if factory == nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("No such game: %s", gameName))
	}

	config := factory.NewConfig()
	if config != nil {
		xecho.MustBindAndValidate(c, &config)
	}
	game := factory.NewGame(config)
	host := CreateHost(gameName, game)
	return c.JSON(http.StatusOK, xecho.J{
		"id": host.Id,
	})
}