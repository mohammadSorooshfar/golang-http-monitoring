package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	routes "github.com/mohammadSorooshfar/golang-http-monitoring/internal/Routes"
)

func main() {

	e := echo.New()
	userGroup := e.Group("/User")
	urlGroup := e.Group("/Url")
	routes.HandleUserRoutes(userGroup)
	routes.HandleUrlRoutes(urlGroup)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
