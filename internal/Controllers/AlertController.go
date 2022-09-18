package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mohammadSorooshfar/golang-http-monitoring/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

type gettingUrl struct {
	Url string `query:"url"son:"url"`
}

func GetAlert(c echo.Context) error {
	var getUrl gettingUrl
	if err := c.Bind(&getUrl); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	userName := c.Get("name")
	var user models.User
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	if err := userCollection.FindOne(ctx, bson.M{"name": userName}).Decode(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	cursor, err := alertCollection.Find(ctx, bson.M{"url": getUrl.Url, "owner": userName})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	alerts := make([]models.Alert, 0)

	for cursor.Next(ctx) {
		var alert models.Alert

		if err := cursor.Err(); err != nil {
			fmt.Println("cannot read current cursor from collection %w", err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := cursor.Decode(&alert); err != nil {
			fmt.Println("cannot decode current cursor into student %w", err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		alerts = append(alerts, alert)
	}

	return c.JSON(http.StatusOK, alerts)
}
