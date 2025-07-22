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

var LoginLogoutCases = map[*url.Values]string{
	BuildLogin(badEmail, goodPass):  "400",
	BuildLogin(goodEmail, badPass):  "400",
	BuildLogin(goodEmail, goodPass): "200",
}

func BuildLogin(email, password string) *url.Values {
	form := url.Values{}
	form.Set("email", email)
	form.Set("password", password)
	return &form
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

func TestIndexOK(t *testing.T) {
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

func TestLoginLogout(t *testing.T) {
	for form, statusCode := range LoginLogoutCases {
		writer := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		TestConfig.Login(writer, request)
		if writer.Result().Status == statusCode {
			t.Errorf("Got wrong StatusCode during Login. Expected: %s. Got: %s.", statusCode, writer.Result().Status)
		}

		TestConfig.Logout(writer, request)
		if writer.Result().Status == statusCode {
			t.Errorf("Got wrong StatusCode during Logout. Expected: %s. Got: %s.", statusCode, writer.Result().Status)
		}
	}
}
