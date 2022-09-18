package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mohammadSorooshfar/golang-http-monitoring/internal/database"
	models "github.com/mohammadSorooshfar/golang-http-monitoring/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var alertCollection *mongo.Collection = database.ConnectToCollection(database.Client, "Alerts")

func onRequestError(err error, ctx context.Context, userId primitive.ObjectID, username string, url models.Url, user models.User, index int, n string) {
	fmt.Println("url : ", url.Link, " we have error in request to url in time : ", time.Now(), " text of error is : ", err.Error())
	user.Urls[index].Allfailed++
	if val, ok := user.Urls[index].Failed[n]; ok {
		user.Urls[index].Failed[n] = val + 1
	}
	if user.Urls[index].Allfailed == user.Urls[index].Threshold {
		fmt.Println("url : ", url.Link, " alert trigered in time : ", time.Now(), "for user : ", username)
		var urlAlert models.Alert
		urlAlert.ID = primitive.NewObjectID()
		urlAlert.Name = username
		urlAlert.Time = time.Now().Format("2006-01-02 15:04:05.000000")
		urlAlert.UserId = userId
		urlAlert.Url = url.Link
		if _, insertErr := alertCollection.InsertOne(ctx, urlAlert); insertErr != nil {
			fmt.Println("database error: ", insertErr)
		}
		user.Urls[index].Allfailed = 0
	}
}

func RequestToUrl(userId primitive.ObjectID, username string, url models.Url, index int) {
	for true {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		if err := userCollection.FindOne(ctx, bson.M{"name": username}).Decode(&user); err != nil {
			fmt.Println("url : ", url.Link, " we have error in finding user in database ", " text of error is : ", err.Error())
		}
		defer cancel()
		client := http.Client{
			Timeout: 10 * time.Second,
		}
		n := time.Now().Format("2006-01-02")
		resp, err := client.Get(url.Link)
		if err != nil {
			onRequestError(err, ctx, userId, username, url, user, index, n)
		} else {
			fmt.Printf("url: %s  statuscode: %d\n", url.Link, resp.StatusCode)
			if resp.StatusCode < 200 && resp.StatusCode > 299 {
				onRequestError(err, ctx, userId, username, url, user, index, n)
			} else {
				if val, ok := user.Urls[index].Success[n]; ok {
					user.Urls[index].Success[n] = val + 1
				}
				user.Urls[index].Allsuccess++
			}
		}
		filter := bson.M{"name": username}
		userCollection.ReplaceOne(ctx, filter, user)
		time.Sleep(10 * time.Second)
	}
}
