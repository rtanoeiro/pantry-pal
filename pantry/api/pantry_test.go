package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

var goodItem = "chocolate"
var goodQuantity = "2"
var badQuantity = "2x"
var negativeQuantity = "-20"
var goodExpiryDate = "2100-12-31"
var badExpiryDate = "2000-01-01"

type ItemCases struct {
	ItemName       string
	ItemQuantity   string
	ItemExpiryDate string
}

func BuildItem(itemName, itemQuantity, itemExpiryDate string) *url.Values {
	form := url.Values{}
	form.Set("itemName", itemName)
	form.Set("itemQuantity", itemQuantity)
	form.Set("itemExpiryDate", itemExpiryDate)
	return &form
}

func AttachItemToRequest(itemName, itemQuantity, itemExpiryDate, method, endpoint string) (*httptest.ResponseRecorder, *http.Request) {
	itemForm := BuildItem(itemName, itemQuantity, itemExpiryDate)
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(method, endpoint, strings.NewReader(itemForm.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return writer, request
}

func TestHandleNewItem(t *testing.T) {
	loginWriter, _ := Login(goodEmail, goodPass)
	expectedStatusCode := 200

	pantryWriter, pantryRequest := AttachItemToRequest(goodItem, goodQuantity, goodExpiryDate, http.MethodPost, "/pantry")
	// After login cookies are added to the writer and passed to the request. So we need to add them back to following requests.
	for _, cookie := range loginWriter.Result().Cookies() {
		pantryRequest.AddCookie(cookie)
	}
	TestConfig.HandleNewItem(pantryWriter, pantryRequest)
	if pantryWriter.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode during item addition. Expected %d. Got: %d.", expectedStatusCode, pantryWriter.Result().StatusCode)
	}
}

func TestHandleNewItemQuantityError(t *testing.T) {
	loginWriter, _ := Login(goodEmail, goodPass)
	expectedStatusCode := 400
	pantryWriter, pantryRequest := AttachItemToRequest(goodItem, badQuantity, goodExpiryDate, http.MethodPost, "/pantry")
	// After login cookies are added to the writer and passed to the request. So we need to add them back to following requests.
	for _, cookie := range loginWriter.Result().Cookies() {
		pantryRequest.AddCookie(cookie)
	}
	TestConfig.HandleNewItem(pantryWriter, pantryRequest)
	if pantryWriter.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode during item addition. Expected %d. Got: %d.", expectedStatusCode, pantryWriter.Result().StatusCode)
	}
}

func TestHandleNewItemDateError(t *testing.T) {
	loginWriter, _ := Login(goodEmail, goodPass)
	expectedStatusCode := 400
	pantryWriter, pantryRequest := AttachItemToRequest(goodItem, goodQuantity, badExpiryDate, http.MethodPost, "/pantry")
	// After login cookies are added to the writer and passed to the request. So we need to add them back to following requests.
	for _, cookie := range loginWriter.Result().Cookies() {
		pantryRequest.AddCookie(cookie)
	}
	TestConfig.HandleNewItem(pantryWriter, pantryRequest)
	if pantryWriter.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode during item addition. Expected %d. Got: %d.", expectedStatusCode, pantryWriter.Result().StatusCode)
	}
}

func TestHandleItemUpdateMoreThanitExists(t *testing.T) {
	loginWriter, _ := Login(goodEmail, goodPass)
	firstExpectedStatusCode := 200
	secondExpectedStatusCode := 400
	pantryWriter, pantryRequest := AttachItemToRequest(goodItem, goodQuantity, goodExpiryDate, http.MethodPost, "/pantry")
	// After login cookies are added to the writer and passed to the request. So we need to add them back to following requests.
	for _, cookie := range loginWriter.Result().Cookies() {
		pantryRequest.AddCookie(cookie)
	}

	TestConfig.HandleNewItem(pantryWriter, pantryRequest)
	if pantryWriter.Result().StatusCode != firstExpectedStatusCode {
		t.Errorf("Got wrong StatusCode during item update. Expected %d. Got: %d.", firstExpectedStatusCode, pantryWriter.Result().StatusCode)
	}

	secondPantryWriter, secondPantryRequest := AttachItemToRequest(goodItem, negativeQuantity, goodExpiryDate, http.MethodPost, "/pantry")
	for _, cookie := range loginWriter.Result().Cookies() {
		secondPantryRequest.AddCookie(cookie)
	}

	TestConfig.HandleNewItem(secondPantryWriter, secondPantryRequest)
	if secondPantryWriter.Result().StatusCode != secondExpectedStatusCode {
		t.Errorf("Got wrong StatusCode during item update. Expected %d. Got: %d.", secondExpectedStatusCode, secondPantryWriter.Result().StatusCode)
	}
}

func TestHandleItemUpdate(t *testing.T) {
	loginWriter, _ := Login(goodEmail, goodPass)
	firstExpectedStatusCode := 200
	secondExpectedStatusCode := 200
	pantryWriter, pantryRequest := AttachItemToRequest(goodItem, goodQuantity, goodExpiryDate, http.MethodPost, "/pantry")
	// After login cookies are added to the writer and passed to the request. So we need to add them back to following requests.
	for _, cookie := range loginWriter.Result().Cookies() {
		pantryRequest.AddCookie(cookie)
	}

	TestConfig.HandleNewItem(pantryWriter, pantryRequest)
	if pantryWriter.Result().StatusCode != firstExpectedStatusCode {
		t.Errorf("Got wrong StatusCode during item update. Expected %d. Got: %d.", firstExpectedStatusCode, pantryWriter.Result().StatusCode)
	}

	secondPantryWriter, secondPantryRequest := AttachItemToRequest(goodItem, goodQuantity, goodExpiryDate, http.MethodPost, "/pantry")
	for _, cookie := range loginWriter.Result().Cookies() {
		secondPantryRequest.AddCookie(cookie)
	}
	TestConfig.HandleNewItem(secondPantryWriter, secondPantryRequest)

	if secondPantryWriter.Result().StatusCode != secondExpectedStatusCode {
		t.Errorf("Got wrong StatusCode during item update. Expected %d. Got: %d.", secondExpectedStatusCode, secondPantryWriter.Result().StatusCode)
	}
}

func TestGetAllPantryItems(t *testing.T) {
	loginWriter, loginRequest := Login(goodEmail, goodPass)
	expectedStatusCode := 200
	// After login cookies are added to the writer and passed to the request. So we need to add them back to following requests.
	for _, cookie := range loginWriter.Result().Cookies() {
		loginRequest.AddCookie(cookie)
	}
	TestConfig.GetAllPantryItems(loginWriter, loginRequest)
	if loginWriter.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode during item update. Expected %d. Got: %d.", expectedStatusCode, loginWriter.Result().StatusCode)
	}

}

func TestDeleteItem(t *testing.T) {
	loginWriter, _ := Login(goodEmail, goodPass)
	expectedStatusCode := 200

	pantryWriter, pantryRequest := AttachItemToRequest(goodItem, goodQuantity, goodExpiryDate, http.MethodPost, "/pantry")
	// After login cookies are added to the writer and passed to the request. So we need to add them back to following requests.
	for _, cookie := range loginWriter.Result().Cookies() {
		pantryRequest.AddCookie(cookie)
	}

	TestConfig.HandleNewItem(pantryWriter, pantryRequest)
	if pantryWriter.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode during item addition. Expected %d. Got: %d.", expectedStatusCode, pantryWriter.Result().StatusCode)
	}

	userID, _ := GetUserIDFromTokenAndValidate(pantryRequest, &TestConfig)
	pantryItems, _ := TestConfig.Db.GetAllItems(context.Background(), userID)
	pantryRequest.SetPathValue("ItemID", pantryItems[0].ID)
	TestConfig.DeleteItem(pantryWriter, pantryRequest)
	if pantryWriter.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode during item deletion. Expected %d. Got: %d.", expectedStatusCode, pantryWriter.Result().StatusCode)
	}

	userID, _ = GetUserIDFromTokenAndValidate(pantryRequest, &TestConfig)
	pantryItems, _ = TestConfig.Db.GetAllItems(context.Background(), userID)
	if len(pantryItems) != 0 {
		t.Errorf("Item was not deleted. Expected 0 items. Got: %d.", len(pantryItems))
	}
}

func TestRenderExpiringSoon(t *testing.T) {
	loginWriter, loginRequest := Login(goodEmail, goodPass)
	expectedStatusCode := 200

	TestConfig.RenderExpiringSoon(loginWriter, loginRequest)
	if loginWriter.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode during item update. Expected %d. Got: %d.", expectedStatusCode, loginWriter.Result().StatusCode)
	}
}

func TestRenderExpiringSoonFail(t *testing.T) {
	loginWriter, loginRequest := Login(goodEmail, goodPass)
	expectedStatusCode := 200

	TestConfig.RenderExpiringSoon(loginWriter, loginRequest)
	if loginWriter.Result().StatusCode != expectedStatusCode {
		t.Errorf("Got wrong StatusCode during item update. Expected %d. Got: %d.", expectedStatusCode, loginWriter.Result().StatusCode)
	}
}
