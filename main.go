package main

import (
	"github.com/labstack/echo"
	"ramenman/handler"
)

func main() {
	e := echo.New()
	e.GET("/", func(e echo.Context) error {
		return handler.HandleRequest(e)
	})
	e.Logger.Fatal(e.Start(":1313"))
}