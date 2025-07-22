package api

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
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
var goodEmail = "admin@admin.com"
var badEmail = "notadmin@admin.com"
var goodPass = "admin"
var badPass = "notadmin"

func BuildLogin(email, password string) url.Values {
	form := url.Values{}
	form.Set("email", email)
	form.Set("password", password)
	return form
}

func TestMain(m *testing.M) {
	log.Println("Setting up Database connection for the test suite...")
	DB, err := sql.Open("sqlite3", TestConfig.DBUrl)
	if err != nil {
		os.Exit(1)
	}
	TestConfig.Db = database.New(DB)

	code := m.Run()

	defer CloseDB(DB)
	log.Println("Closing Database connection for the test suite...")
	os.Exit(code)
}

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
	goodForm := BuildLogin(goodEmail, goodPass)
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(goodForm.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	TestConfig.Login(writer, request)
	if writer.Result().Status == "200" {
		t.Errorf("Expected 200 status code, got %s.", writer.Result().Status)
	}
}
