package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"pantry-pal/pantry/database"
	"testing"
)

var ItemCardTestCases = map[database.GetAllShoppingRow]int{
	{ItemName: "Test item", Quantity: 2}:          200,
	{ItemName: "Another Test Item", Quantity: -1}: 400,
}

func BuildAddItemShopping(itemName, itemQuantity string) *url.Values {
	form := url.Values{}
	form.Set("itemName", itemName)
	form.Set("itemQuantity", itemQuantity)
	return &form
}

func TestGetAllShopping(t *testing.T) {
	loginWriter, _ := Login(goodUser, goodPass)
	expectedStatusCode := 200

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/shopping", nil)
	for _, cookie := range loginWriter.Result().Cookies() {
		request.AddCookie(cookie)
	}

	TestConfig.AddItemShopping(writer, request)
	if writer.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode when Updatin User Password. Expected %d. Got: %d.", expectedStatusCode, writer.Result().StatusCode)
	}
}

func TestAddItemShopping(t *testing.T) {
	loginWriter, _ := Login(goodUser, goodPass)
	expectedStatusCode := 200

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/shopping", nil)
	for _, cookie := range loginWriter.Result().Cookies() {
		request.AddCookie(cookie)
	}
	TestConfig.AddItemShopping(writer, request)
	if writer.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode when Updatin User Password. Expected %d. Got: %d.", expectedStatusCode, writer.Result().StatusCode)
	}
}

func TestRemoveItemShopping(t *testing.T) {
	loginWriter, _ := Login(goodUser, goodPass)
	expectedStatusCode := 200

	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/shopping", nil)
	for _, cookie := range loginWriter.Result().Cookies() {
		request.AddCookie(cookie)
	}
	TestConfig.RemoveItemShopping(writer, request)
	if writer.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode when Updatin User Password. Expected %d. Got: %d.", expectedStatusCode, writer.Result().StatusCode)
	}
}
