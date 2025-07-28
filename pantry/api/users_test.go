package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var goodUserName = "John Doe"
var goodUserPass = "testpass"
var badUserName = ""
var badUserPass = ""
var badUserEmail = ""
var newName = "New John"
var newPass = "newtestpass"

type CreateUserCases struct {
	Name     string
	Password string
}

func BuildPerson(name, password string) *url.Values {
	form := url.Values{}
	form.Set("name", name)
	form.Set("password", password)
	return &form
}

func AttachUserToRequest(name, password, method, endpoint string) (*httptest.ResponseRecorder, *http.Request) {
	itemForm := BuildPerson(name, password)
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(method, endpoint, strings.NewReader(itemForm.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return writer, request
}

func AddUser(name, password string) (*httptest.ResponseRecorder, *http.Request) {
	writer, request := AttachUserToRequest(name, password, http.MethodPost, "/signup")
	TestConfig.CreateUser(writer, request)
	return writer, request
}

func TestCreateUser(t *testing.T) {
	writer, _ := Login(goodUser, goodPass)
	expectedStatusCodes := map[CreateUserCases]int{
		{Name: goodUserName, Password: goodUserPass}: 200,
		{Name: goodUserName, Password: badUserPass}:  400,
		{Name: badUserName, Password: badUserPass}:   400,
		{Name: badUserName, Password: badUserPass}:   400,
		{Name: goodUserName, Password: badUserPass}:  400,
		{Name: goodUserName, Password: badUserPass}:  400,
	}

	for userCase, expectedCode := range expectedStatusCodes {
		userWriter, userRequest := AddUser(userCase.Name, userCase.Password)
		if userWriter.Result().StatusCode != expectedCode {
			t.Errorf("Got wrong StatusCode during user addition. Expected %d. Got: %d.", expectedCode, userWriter.Result().StatusCode)
		}
		if userWriter.Result().StatusCode == 200 {
			addedUser, _ := TestConfig.Db.GetUserByName(userRequest.Context(), userCase.Name)
			DeleteWriter := httptest.NewRecorder()
			DeleteRequest := httptest.NewRequest(http.MethodDelete, "/user", nil)
			DeleteRequest.SetPathValue("UserID", addedUser.ID)
			for _, cookie := range writer.Result().Cookies() {
				DeleteRequest.AddCookie(cookie)
			}

			TestConfig.DeleteUser(DeleteWriter, DeleteRequest)
		}
	}
}

func TestGetUserInfo(t *testing.T) {
	writer, request := Login(goodUser, goodPass)
	expectedCode := 200
	for _, cookie := range writer.Result().Cookies() {
		request.AddCookie(cookie)
	}
	TestConfig.GetUserInfo(writer, request)
	if writer.Result().StatusCode != expectedCode {
		t.Errorf("Got wrong StatusCode when getting user Data. Expected %d. Got: %d.", expectedCode, writer.Result().StatusCode)
	}
}

func TestGetUserInfoError(t *testing.T) {
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/user", nil)
	expectedCode := 403
	TestConfig.GetUserInfo(writer, request)
	if writer.Result().StatusCode != expectedCode {
		t.Errorf("Got wrong StatusCode when getting user Data. Expected %d. Got: %d.", expectedCode, writer.Result().StatusCode)
	}
}

func TestCreateUpdateDeleteUser(t *testing.T) {
	expectedCodeCreate := 200
	addWriter, _ := AddUser(goodUserName, goodUserPass)
	if addWriter.Result().StatusCode != expectedCodeCreate {
		t.Errorf("Got wrong StatusCode during user addition. Expected %d. Got: %d.", expectedCodeCreate, addWriter.Result().StatusCode)
	}

	loginWriter, _ := Login(goodUser, goodUserPass)
	expectedCodeUpdate := 200
	expectedCodeDelete := 200

	updateWriter, updateRequest := AttachUserToRequest(newName, "", http.MethodPost, "/user/name")
	for _, cookie := range loginWriter.Result().Cookies() {
		updateRequest.AddCookie(cookie)
	}
	TestConfig.UpdateUserName(updateWriter, updateRequest)
	if updateWriter.Result().StatusCode != expectedCodeUpdate {
		t.Errorf("Got wrong StatusCode when Updatin User Name. Expected %d. Got: %d.", expectedCodeDelete, updateWriter.Result().StatusCode)
	}

	updateWriter, updateRequest = AttachUserToRequest("", newPass, http.MethodPost, "/user/password")
	for _, cookie := range loginWriter.Result().Cookies() {
		updateRequest.AddCookie(cookie)
	}
	TestConfig.UpdateUserPassword(updateWriter, updateRequest)
	if updateWriter.Result().StatusCode != expectedCodeUpdate {
		t.Errorf("Got wrong StatusCode when Updatin User Password. Expected %d. Got: %d.", expectedCodeDelete, updateWriter.Result().StatusCode)
	}

	addedUser, _ := TestConfig.Db.GetUserByName(updateRequest.Context(), newName)
	DeleteWriter := httptest.NewRecorder()
	DeleteRequest := httptest.NewRequest(http.MethodDelete, "/user", nil)
	DeleteRequest.SetPathValue("UserID", addedUser.ID)
	for _, cookie := range loginWriter.Result().Cookies() {
		DeleteRequest.AddCookie(cookie)
	}

	TestConfig.DeleteUser(DeleteWriter, DeleteRequest)
	if DeleteWriter.Result().StatusCode != expectedCodeDelete {
		t.Errorf("Got wrong StatusCode when Deleting User. Expected %d. Got: %d.", expectedCodeDelete, DeleteWriter.Result().StatusCode)
	}

}

func TestAddRevokeAdmin(t *testing.T) {
	userWriter, userRequest := AddUser(goodUserName, goodUserPass)
	for _, cookie := range userWriter.Result().Cookies() {
		userRequest.AddCookie(cookie)
	}

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/user", nil)

	loginWriter, _ := Login(goodUser, goodPass)
	for _, cookie := range loginWriter.Result().Cookies() {
		request.AddCookie(cookie)
	}
	addedUser, _ := TestConfig.Db.GetUserByName(context.Background(), goodUserName)
	request.SetPathValue("UserID", addedUser.ID)
	TestConfig.AddUserAdmin(writer, request)
	TestConfig.RevokeUserAdmin(writer, request)
}
