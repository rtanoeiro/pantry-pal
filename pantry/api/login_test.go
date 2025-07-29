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
var goodUser = "Admin"
var goodPass = "admin"
var badUser = "notadmin"
var badPass = "notadmin"

type LoginCases struct {
	User     string
	Password string
}

var LoginLogoutCases = map[LoginCases]int{
	{User: badUser, Password: badPass}:   400,
	{User: badUser, Password: goodPass}:  400,
	{User: goodUser, Password: badPass}:  400,
	{User: goodUser, Password: goodPass}: 200,
}

func BuildLogin(username, password string) *url.Values {
	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)
	return &form
}

func Login(username, password string) (*httptest.ResponseRecorder, *http.Request) {
	loginForm := BuildLogin(username, password)
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
		writer, _ := Login(login.User, login.Password)
		if writer.Result().StatusCode != statusCode {
			t.Errorf("Got wrong StatusCode during Login. Expected %d. Got: %d.", statusCode, writer.Result().StatusCode)
		}

		homeWriter := httptest.NewRecorder()
		homeRequest := httptest.NewRequest(http.MethodGet, "/home", nil)
		for _, cookie := range writer.Result().Cookies() {
			homeRequest.AddCookie(cookie)
		}

		TestConfig.Logout(homeWriter, homeRequest)
		if writer.Result().StatusCode != statusCode {
			t.Errorf("Got wrong StatusCode during Logout. Expected %d. Got: %d.", statusCode, writer.Result().StatusCode)
		}
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

func TestHome(t *testing.T) {
	writer, _ := Login(goodUser, goodPass)
	expectedStatusCode := 200
	if writer.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode during Login. Expected %d. Got: %d.", expectedStatusCode, writer.Result().StatusCode)
	}

	//We have to build a new request, as the one from Login is pointing to a POST/login
	request := httptest.NewRequest(http.MethodGet, "/home", nil)
	for _, cookie := range writer.Result().Cookies() {
		request.AddCookie(cookie)
	}

	TestConfig.Home(writer, request)
	if writer.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode on Home. Expected %d. Got: %d.", expectedStatusCode, writer.Result().StatusCode)
	}
}
