package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAlert(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusCreated, id)
}
