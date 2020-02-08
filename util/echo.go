package util

import (
	"github.com/labstack/echo/v4"
	"github.com/xdean/goex/xecho"
	"net/http"
	"strconv"
)

func IntParam(c echo.Context, name string) int {
	param := c.Param(name)
	if i, err := strconv.Atoi(param); err == nil {
		return i
	} else {
		xecho.MustNoError(echo.NewHTTPError(http.StatusBadRequest, "Unrecognized param '"+name+"': "+param))
		return 0
	}
}
