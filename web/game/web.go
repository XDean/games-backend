package game

import (
	"games-backend/app"
	"games-backend/games/host"
	"games-backend/util"
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
	id := util.IntParam(c, "id")

	server, ok := getServer(gameName, id)
	if !ok {

		_ = ws.WriteJSON(host.ErrorEvent("房间不存在"))
		return nil
	}

	client := server.newClient(user, ws)
	client.run()
	return nil
}

func createGame(c echo.Context) error {
	gameName := c.Param("game")

	server, err := createServer(gameName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	server.run()

	return c.JSON(http.StatusOK, xecho.J{
		"id": server.id,
	})
}
