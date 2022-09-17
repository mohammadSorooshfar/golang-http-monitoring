package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	models "github.com/mohammadSorooshfar/golang-http-monitoring/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func RequestToUrl(username string, url models.Url, index int) {
	for true {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		if err := userCollection.FindOne(ctx, bson.M{"name": username}).Decode(&user); err != nil {
			fmt.Println("url : ", url.Link, " we have error in finding user in database ", " text of error is : ", err.Error())
		}
		defer cancel()
		resp, err := http.Get(url.Link)
		if err != nil {
			fmt.Println("url : ", url.Link, " we have error in request to url in time : ", time.Now(), " text of error is : ", err.Error())
		}
		fmt.Printf("url: %s  statuscode: %d\n", url.Link, resp.StatusCode)
		n := time.Now().Format("2006-01-02")
		if resp.StatusCode < 200 && resp.StatusCode > 299 {
			user.Urls[index].Allfailed++
			if val, ok := user.Urls[index].Failed[n]; ok {
				user.Urls[index].Failed[n] = val + 1
			}
			if user.Urls[index].Allfailed == user.Urls[index].Threshold {
				fmt.Println("url : ", url.Link, " alert trigered in time : ", time.Now(), "for user : ", username)
				user.Urls[index].Allfailed = 0
			}

		} else {
			if val, ok := user.Urls[index].Failed[n]; ok {
				user.Urls[index].Failed[n] = val + user.Urls[index].Failed[n]
			}
			user.Urls[index].Allsuccess++
		}
		filter := bson.M{"name": username}
		userCollection.ReplaceOne(ctx, filter, user)
		time.Sleep(10 * time.Second)
	}
}
