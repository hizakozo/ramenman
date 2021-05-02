package main

import (
	"github.com/labstack/echo"
	"net/http"
	"ramenman/handler"
)

func HandleRequest(c echo.Context) error {
	resp, err := handler.GetShops(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp.Results)
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return HandleRequest(c)
	})
	e.Logger.Fatal(e.Start(":1313"))
}