package api

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func BuildAddItemShopping(itemName, itemQuantity string) *url.Values {
	form := url.Values{}
	form.Set("itemName", itemName)
	form.Set("itemQuantity", itemQuantity)
	return &form
}

func AttachItemToShoppingRequest(itemName, itemQuantity, method, endpoint string) (*httptest.ResponseRecorder, *http.Request) {
	itemForm := BuildAddItemShopping(itemName, itemQuantity)
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(method, endpoint, strings.NewReader(itemForm.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return writer, request
}

func TestGetAllShopping(t *testing.T) {
	loginWriter, _ := Login(goodUser, goodPass)
	expectedStatusCode := 200

	writer, request := AttachItemToShoppingRequest("beans", "1", http.MethodPost, "/shopping")
	for _, cookie := range loginWriter.Result().Cookies() {
		request.AddCookie(cookie)
	}

	TestConfig.AddItemShopping(writer, request)
	if writer.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode when getting all items from User Shopping Cart. Expected %d. Got: %d.", expectedStatusCode, writer.Result().StatusCode)
	}
}

func TestAddItemShopping(t *testing.T) {
	loginWriter, _ := Login(goodUser, goodPass)
	expectedStatusCode := 200

	writer, request := AttachItemToShoppingRequest("tomato", "1", http.MethodPost, "/shopping")
	for _, cookie := range loginWriter.Result().Cookies() {
		request.AddCookie(cookie)
	}
	TestConfig.AddItemShopping(writer, request)
	if writer.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode when Adding Item to User Shopping Cart. Expected %d. Got: %d.", expectedStatusCode, writer.Result().StatusCode)
	}
}

func TestRemoveItemShopping(t *testing.T) {
	loginWriter, _ := Login(goodUser, goodPass)
	expectedStatusCode := 200

	writer, request := AttachItemToShoppingRequest("rice", "1", http.MethodPost, "/shopping")
	for _, cookie := range loginWriter.Result().Cookies() {
		request.AddCookie(cookie)
	}
	TestConfig.AddItemShopping(writer, request)

	removeWriter, removeRequest := AttachItemToShoppingRequest("rice", "0", http.MethodDelete, "/shopping/rice")
	for _, cookie := range loginWriter.Result().Cookies() {
		removeRequest.AddCookie(cookie)
	}

	TestConfig.RemoveItemShopping(removeWriter, removeRequest)
	if removeWriter.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode when Adding Item to User Shopping Cart. Expected %d. Got: %d.", expectedStatusCode, writer.Result().StatusCode)
	}
}
