package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pantry-pal/pantry/database"
	"strconv"
	"time"
)

func (config *Config) GetAllShopping(writer http.ResponseWriter, request *http.Request) {
	var cartInfo CartInfo
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		cartInfo.ErrorMessage = fmt.Sprintf("Unable to get current user data. Error: %s", errUser.Error())
		writer.WriteHeader(http.StatusForbidden)
		_ = config.Renderer.Render(writer, "ResponseMessage", cartInfo)
		return
	}
	log.Printf("Loading all cart items for user %s", userID)
	cartItems, errCart := config.Db.GetAllShopping(context.Background(), userID)
	if errCart != nil {
		cartInfo.ErrorMessage = fmt.Sprintf("Unable to get current user cart info. Error: %s ", errCart.Error())
		writer.WriteHeader(http.StatusForbidden)
		_ = config.Renderer.Render(writer, "ResponseMessage", cartInfo)
	}

	cartInfo.CartItems = cartItems
	writer.WriteHeader(http.StatusOK)
	_ = config.Renderer.Render(writer, "shoppingCart", cartInfo)
}

func (config *Config) AddItemShopping(writer http.ResponseWriter, request *http.Request) {
	var cartInfo CartInfo
	itemName := request.FormValue("itemName")
	itemQuantity, errQuantity := strconv.Atoi(request.FormValue("itemQuantity"))
	if errQuantity != nil {
		cartInfo.ErrorMessage = "Invalid quantity"
		writer.WriteHeader(http.StatusBadRequest)
		_ = config.Renderer.Render(writer, "ResponseMessage", cartInfo)
		return
	}
	if itemName == "" || itemQuantity == 0 {
		cartInfo.ErrorMessage = "Please provide valid Name and Quantity for all fields"
		writer.WriteHeader(http.StatusBadRequest)
		_ = config.Renderer.Render(writer, "ResponseMessage", cartInfo)
		return
	}

	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		cartInfo.ErrorMessage = fmt.Sprintf("Unable to retrieve user. Error: %s", errUser.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		_ = config.Renderer.Render(writer, "ResponseMessage", cartInfo)
		return
	}
	log.Printf("User %s is trying to add %d of %s into its shopping cart", userID, itemQuantity, itemName)

	findItem := database.FindItemShoppingParams{
		ItemName: itemName,
		UserID:   userID,
	}
	currentItem, errFind := config.Db.FindItemShopping(request.Context(), findItem)
	if errFind == nil {
		newQuantity := currentItem.Quantity + int64(itemQuantity)
		config.UpdateItemShopping(writer, request, newQuantity, currentItem.ItemName, userID)
		return
	}

	addItem := database.AddItemShoppingParams{
		UserID:   userID,
		ItemName: itemName,
		Quantity: int64(itemQuantity),
	}
	errAdd := config.Db.AddItemShopping(request.Context(), addItem)
	if errAdd != nil {
		cartInfo.ErrorMessage = fmt.Sprintf("Unable to add items to your Shopping Cart. Error: %s", errAdd.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		_ = config.Renderer.Render(writer, "ResponseMessage", cartInfo)
		return
	}

	writer.WriteHeader(http.StatusOK)
	cartInfo.SuccessMessage = fmt.Sprintf("Successfully added x%d - %s", addItem.Quantity, addItem.ItemName)
	_ = config.Renderer.Render(writer, "shoppingCart", cartInfo)
}

func (config *Config) AddOneItemShopping(writer http.ResponseWriter, request *http.Request) {
	var cartInfo SuccessErrorResponse
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		cartInfo.ErrorMessage = fmt.Sprintf("Unable to retrieve user. Error: %s", errUser.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		_ = config.Renderer.Render(writer, "ResponseMessage", cartInfo)
		return
	}
	log.Printf("User %s is trying to add one item of %s", userID, request.FormValue("itemName"))
	request.Form.Add("itemName", request.FormValue("itemName"))
	request.Form.Add("itemQuantity", "1")
	config.AddItemShopping(writer, request)
}

func (config *Config) UpdateItemShopping(writer http.ResponseWriter, request *http.Request, newQuantity int64, itemName, userID string) {
	var cartInfo SuccessErrorResponse
	updateItem := database.UpdateItemShoppingParams{
		Quantity: newQuantity,
		ItemName: itemName,
		UserID:   userID,
	}
	errUpdate := config.Db.UpdateItemShopping(request.Context(), updateItem)
	if errUpdate != nil {
		cartInfo.ErrorMessage = "Invalid quantity"
		writer.WriteHeader(http.StatusBadRequest)
		_ = config.Renderer.Render(writer, "ResponseMessage", cartInfo)
		return
	}
	cartInfo.SuccessMessage = fmt.Sprintf("Added %d items of %s", updateItem.Quantity, updateItem.ItemName)
	writer.WriteHeader(http.StatusOK)
	_ = config.Renderer.Render(writer, "shoppingCart", cartInfo)
}

func (config *Config) RemoveItemShopping(writer http.ResponseWriter, request *http.Request) {
	var cartInfo SuccessErrorResponse
	itemID := request.PathValue("ItemID")
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		log.Printf("Error on deleting item from user %s at %s. Invalid token", userID, time.Now())
		cartInfo.ErrorMessage = fmt.Sprintf("Unable to retrieve user Pantry Items. Error: %s", errUser.Error())
		writer.WriteHeader(http.StatusForbidden)
		_ = config.Renderer.Render(writer, "ResponseMessage", cartInfo)
		return
	}
	log.Printf("User %s is trying to remove %s from shopping cart at %s.", userID, itemID, time.Now())

	removeParams := database.RemoveItemShoppingParams{
		UserID:   userID,
		ItemName: userID,
	}
	errRemove := config.Db.RemoveItemShopping(request.Context(), removeParams)
	if errRemove != nil {
		log.Printf("Error on deleting item %s from user %s shopping cart at %s. Error: %s", removeParams.ItemName, userID, time.Now(), errRemove.Error())
		cartInfo.ErrorMessage = fmt.Sprintf("Error on deleting item from shopping cart. Error: %s", errRemove.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		_ = config.Renderer.Render(writer, "ResponseMessage", cartInfo)
		return
	}
	log.Printf("User %s removed %s from its shopping cart at %s", userID, removeParams.ItemName, time.Now())
	cartInfo.SuccessMessage = fmt.Sprintf("Successfulyl deleted %s: ", removeParams.ItemName)
	writer.WriteHeader(http.StatusOK)
	_ = config.Renderer.Render(writer, "shoppingCart", cartInfo)
}
