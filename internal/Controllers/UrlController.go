package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v4"
	models "github.com/mohammadSorooshfar/golang-http-monitoring/internal/Models"
	"github.com/mohammadSorooshfar/golang-http-monitoring/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var urlCollection *mongo.Collection = database.ConnectToCollection(database.Client, "Urls")

type UrlJson struct {
	Username string `json:"name" xml:"name"`
	Url      string `json:"url" xml:"mongoid"`
	Message  string `json:"message"`
}

func CreateUrl(c echo.Context) error {
	userName := c.Get("name")
	fmt.Println(userName)
	var user models.User
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	if err := userCollection.FindOne(ctx, bson.M{"name": userName}).Decode(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var url models.Url
	if err := c.Bind(&url); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	fmt.Println(url.Link)
	if len(user.Urls) >= 20 {
		return echo.NewHTTPError(http.StatusBadRequest, "you can add 20 urls only!!!")
	}

	for _, x := range user.Urls {
		if x.Link == url.Link {
			return echo.NewHTTPError(http.StatusBadRequest, "duplicae url!!!")
		}
	}
	if validationErr := validation.ValidateStruct(&url,
		validation.Field(&url.Period,
			validation.Required, validation.Min(1)),
		validation.Field(&url.Threshold, is.Digit, validation.Min(5)),
		validation.Field(&url.Link, validation.Required, is.URL),
	); validationErr != nil {
		fmt.Println(validationErr)
		return echo.NewHTTPError(http.StatusBadRequest, validationErr.Error())

	}
	if url.Threshold == 0 {
		url.Threshold = 5
	}
	url.Failed = 0
	user.Urls = append(user.Urls, url)
	url.ID = primitive.NewObjectID()
	if _, insertErr := urlCollection.InsertOne(ctx, url); insertErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, insertErr.Error())
	}
	_, err := userCollection.ReplaceOne(ctx, bson.M{"name": userName}, user)
	if err != nil {
		if _, removeErr := urlCollection.DeleteOne(ctx, bson.M{"_id": url.ID}); removeErr != nil {
			log.Fatal("url document saved in url collection but not in user urls")
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	u := &UrlJson{
		Username: user.Name,
		Url:      url.Link,
		Message:  "url added",
	}
	return c.JSON(http.StatusCreated, u)

}
