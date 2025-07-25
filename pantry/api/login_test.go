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
	DBUrl:    "../../data/pantry_pal_dev.db",
	Port:     "8080",
	Renderer: &MockRenderer{},
}
var goodEmail = "admin@admin.com"
var goodPass = "admin"
var badEmail = "notadmin@admin.com"
var badPass = "notadmin"

type LoginCases struct {
	Email    string
	Password string
}

var LoginLogoutCases = map[LoginCases]int{
	{Email: badEmail, Password: badPass}:   400,
	{Email: badEmail, Password: goodPass}:  400,
	{Email: goodEmail, Password: goodPass}: 200,
}

func BuildLogin(email, password string) *url.Values {
	form := url.Values{}
	form.Set("email", email)
	form.Set("password", password)
	return &form
}

func Login(email, password string) (*httptest.ResponseRecorder, *http.Request) {
	loginForm := BuildLogin(email, password)
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(loginForm.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	TestConfig.Login(writer, request)
	return writer, request
}

func TestMain(m *testing.M) {
	log.Println("Setting up Database connection for the test suite...")
	DB, err := sql.Open("sqlite3", TestConfig.DBUrl)
	if err != nil {
		os.Exit(1)
	}
	TestConfig.Db = database.New(DB)
	code := m.Run()

	log.Println("Closing Database connection for the test suite...")
	defer CloseDB(DB)
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
	for login, statusCode := range LoginLogoutCases {
		writer, request := Login(login.Email, login.Password)
		if writer.Result().StatusCode != statusCode {
			t.Errorf("Got wrong StatusCode during Login. Expected %d. Got: %d.", statusCode, writer.Result().StatusCode)
		}

		TestConfig.Logout(writer, request)
		if writer.Result().StatusCode != statusCode {
			t.Errorf("Got wrong StatusCode during Logout. Expected %d. Got: %d.", statusCode, writer.Result().StatusCode)
		}
	}
}

func TestHome(t *testing.T) {
	writer, request := Login(goodEmail, goodPass)
	expectedStatusCode := 200
	if writer.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode during Login. Expected %d. Got: %d.", expectedStatusCode, writer.Result().StatusCode)
	}

	TestConfig.Home(writer, request)
	if writer.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode on Home. Expected %d. Got: %d.", expectedStatusCode, writer.Result().StatusCode)
	}
}

func TestHomeError(t *testing.T) {
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/home", nil)
	TestConfig.Home(writer, request)
	expectedStatusCode := 401
	if writer.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode on home screen. Expected %d. Got: %d.", expectedStatusCode, writer.Result().StatusCode)
	}
}
