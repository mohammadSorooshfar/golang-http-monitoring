package routes

import (
	"github.com/labstack/echo/v4"
	controllers "github.com/mohammadSorooshfar/golang-http-monitoring/internal/Controllers"
)

func HandleUrlRoutes(group *echo.Group) {
	group.POST("/create", controllers.CreateUrl)
}
