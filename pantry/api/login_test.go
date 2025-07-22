package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"pantry-pal/pantry/database"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var TestConfig = Config{
	Secret:   "SuperTestSecret",
	DBUrl:    "data/pantry_pal_dev.db",
	Port:     "8080",
	Renderer: &MockRenderer{},
}
var adminEmail = "admin@admin.com"
var adminPass = "admin"

func TestIndex(t *testing.T) {
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/login", nil)

	TestConfig.Index(writer, request)

	if writer.Header().Get("HX-Replace-Url") != "/login" {
		t.Errorf("Expected HX-Replace-Url to be set")
	}

	if writer.Result().StatusCode != 200 {
		t.Errorf("Expected 200 status code. Got: %d", writer.Result().StatusCode)
	}
}

func TestSignup(t *testing.T) {
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/signup", nil)
	TestConfig.SignUp(writer, request)
}

func TestLogin(t *testing.T) {
	DB, dbError := sql.Open("sqlite3", TestConfig.DBUrl)
	if dbError != nil {
		log.Fatal("Unable to open database. Closing app. Error: ", dbError)
	}
	defer CloseDB(DB)
	TestConfig.Db = database.New(DB)

	writer := httptest.NewRecorder()
	form := url.Values{}
	form.Set("email", adminEmail)
	form.Set("password", adminPass)
	fmt.Println(form)

	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	TestConfig.Login(writer, request)
	if writer.Result().Status == "200" {
		t.Errorf("Expected 200 status code, got %s.", writer.Result().Status)
	}

}
