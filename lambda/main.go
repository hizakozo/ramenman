package main

import (
	"github.com/aws/aws-lambda-go/lambda"
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
	lambda.Start(HandleRequest)
}