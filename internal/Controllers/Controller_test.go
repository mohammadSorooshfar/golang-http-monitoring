package controllers

import (
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
	userJSON = `{"name":"JonSnow2","password":"123459876"}`
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
		fmt.Println(rec.Body.String())
		assert.Equal(t, http.StatusCreated, rec.Code)
	}

}
