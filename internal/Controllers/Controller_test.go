package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mohammadSorooshfar/golang-http-monitoring/internal/database"
	"github.com/stretchr/testify/assert"
)

var (
	userJSON = `{"name":"JonSnowoaa","password":"123459876"}`
	urlJson  = `{"link":"www.google.com","period":5}`
	token    = ""
)

func TestCreateUser(t *testing.T) {
	// Setup
	userCollection = database.ConnectToCollection(database.CreateClient(), "Users")
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/User/signup", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, SignUp(c)) {
		fmt.Println("         //////////////////start of create user response/////////////////////")
		fmt.Println("Response is: ", rec.Body.String())
		fmt.Println("         //////////////////end of create user response///////////////////////")
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

	fmt.Println("//////////////////////////////////////////////////////////////////////////////////")
	fmt.Println()
	fmt.Println()
	fmt.Println()

}

func TestLoginUser(t *testing.T) {
	// Setup
	userCollection = database.ConnectToCollection(database.CreateClient(), "Users")
	urlCollection = database.ConnectToCollection(database.Client, "Urls")
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/User/login", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, Login(c)) {
		var testtoken UserToken
		fmt.Println("                 //////////////////start of login response/////////////////////")
		fmt.Println("Response is: ", rec.Body.String())
		fmt.Println("                //////////////////end of login response///////////////////////")
		err := json.Unmarshal(rec.Body.Bytes(), &testtoken)
		if err != nil {
			fmt.Println("            //////////////////start of login error/////////////////////")
			fmt.Println(err.Error())
			fmt.Println("            //////////////////end of login error/////////////////////")
			return
		}
		token = "Bearer " + testtoken.Token
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	fmt.Println("//////////////////////////////////////////////////////////////////////////////////")
	fmt.Println()
	fmt.Println()
	fmt.Println()

}

func TestCreateUrl(t *testing.T) {
	// Setup
	userCollection = database.ConnectToCollection(database.CreateClient(), "Users")
	urlCollection = database.ConnectToCollection(database.Client, "Urls")
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/Url/create", strings.NewReader(urlJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Add("Authorization", token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, Auth(echo.HandlerFunc(CreateUrl))(c)) {
		fmt.Println("           //////////////////start of create url  response/////////////////////")
		fmt.Println("Response is: ", rec.Body.String())
		fmt.Println("           //////////////////end of login response///////////////////////")
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

}
