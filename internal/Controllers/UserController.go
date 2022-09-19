package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/mohammadSorooshfar/golang-http-monitoring/internal/database"
	helpermethods "github.com/mohammadSorooshfar/golang-http-monitoring/internal/helperMethods"
	models "github.com/mohammadSorooshfar/golang-http-monitoring/internal/models"
	"github.com/mohammadSorooshfar/golang-http-monitoring/internal/request"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Name    string `json:"name" xml:"name"`
	Mongoid string `json:"mongoid" xml:"mongoid"`
	Message string `json:"message"`
}

type UserToken struct {
	Name    string `json:"name" xml:"name"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

const (
	StudentIDLen    = 7
	FirstNameMaxLen = 255
	FirstNameMinLen = 1
	LastNameMaxLen  = 255
	LastNameMinLen  = 1
)

var userCollection *mongo.Collection = database.ConnectToCollection(database.Client, "Users")

func SignUp(c echo.Context) error {
	if userCollection == nil {

	}
	var user models.User

	if err := c.Bind(&user); err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if validationErr := validation.ValidateStruct(&user,
		validation.Field(&user.Name,
			validation.Required, validation.Length(FirstNameMinLen, FirstNameMaxLen), is.UTFLetterNumeric),
		validation.Field(&user.Password,
			validation.Required, validation.Length(LastNameMinLen, LastNameMaxLen)),
	); validationErr != nil {
		fmt.Println(validationErr)
		return echo.NewHTTPError(http.StatusBadRequest, validationErr.Error())
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	count, err := userCollection.CountDocuments(ctx, bson.M{"name": user.Name})
	if err != nil {
		log.Panic(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if count > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "user with this username already exist!!!")
	}
	password := helpermethods.HashPassword(user.Password)
	user.Password = password
	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	_, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user cannot created pls try again")
	}
	for_send := user.ID.Hex()
	fmt.Println(for_send)
	u := &User{
		Name:    user.Name,
		Mongoid: for_send,
		Message: "created",
	}
	return c.JSON(http.StatusCreated, u)
}

func Login(c echo.Context) error {
	var req request.Login
	var user models.User
	if err := c.Bind(&req); err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	err := userCollection.FindOne(ctx, bson.M{"name": req.Username}).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "username or password is wrong!!!")
	}

	passwordIsValid, _ := helpermethods.VerifyPassword(req.Password, user.Password)
	if passwordIsValid != true {
		return echo.NewHTTPError(http.StatusBadRequest, "username or password is wrong!!!")
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    "students-summer-2022",
		Subject:   req.Username,
		Audience:  []string{"admin"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		NotBefore: jwt.NewNumericDate(time.Now()),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        user.ID.Hex(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return echo.ErrInternalServerError
	}
	u := &UserToken{
		Name:    user.Name,
		Token:   tokenString,
		Message: "created",
	}

	return c.JSON(http.StatusOK, u)
}
