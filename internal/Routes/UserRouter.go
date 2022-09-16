package routes

import (
	"github.com/labstack/echo/v4"
	controllers "github.com/mohammadSorooshfar/golang-http-monitoring/internal/Controllers"
)

func HandleUserRoutes(group *echo.Group) {
	group.POST("/signup", controllers.SignUp)
}
