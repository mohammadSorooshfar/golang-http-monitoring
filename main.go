package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mohammadSorooshfar/golang-http-monitoring/handler"
	"github.com/mohammadSorooshfar/golang-http-monitoring/store"
)

func main() {
	app := echo.New()

	var userStore store.Auth
	{
		// userStore = store.NewAuthInMemory()
	}
	ha := handler.Auth{
		Store: userStore,
	}

	ha.Register(app.Group(""))
	app.Logger.Fatal(app.Start(":3000"))
}
