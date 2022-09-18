package routes

import (
	"github.com/labstack/echo/v4"
	controllers "github.com/mohammadSorooshfar/golang-http-monitoring/internal/Controllers"
)

func HandleAlertRoutes(group *echo.Group) {
	group.GET("/getAlert", controllers.GetAlert)
}
