package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/mohammadSorooshfar/golang-http-monitoring/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slices"
)

type gettingUrl struct {
	Url string `query:"url"json:"url"`
}

func GetAlert(c echo.Context) error {
	var getUrl gettingUrl
	if err := c.Bind(&getUrl); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if validationErr := validation.ValidateStruct(&getUrl,
		validation.Field(&getUrl.Url,
			validation.Required),
	); validationErr != nil {
		fmt.Println(validationErr)
		return echo.NewHTTPError(http.StatusBadRequest, validationErr.Error())
	}
	userName := c.Get("name")
	var user models.User
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	if err := userCollection.FindOne(ctx, bson.M{"name": userName}).Decode(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	idx := slices.IndexFunc(user.Urls, func(c models.Url) bool { return c.Link == getUrl.Url })
	if idx == -1 {
		return echo.NewHTTPError(http.StatusBadRequest, "this url is not in your saved urls")
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
