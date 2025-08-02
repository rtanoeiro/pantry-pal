package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pantry-pal/pantry/database"
	"strconv"
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
		currentItem.Quantity = currentItem.Quantity + int64(itemQuantity)
		log.Printf("Item already exists. User %s is updating %s into its shopping cart to %d items", userID, itemName, currentItem.Quantity)
		config.UpdateItemShopping(writer, request, currentItem)
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
		config.UpdateItemShopping(writer, request, currentItem)
		return
	}

	writer.WriteHeader(http.StatusBadRequest)
	cartInfo.SuccessMessage = fmt.Sprintf("Successfully added x%d - %s", addItem.Quantity, addItem.ItemName)
	_ = config.Renderer.Render(writer, "ResponseMessage", cartInfo)
}

func (config *Config) UpdateItemShopping(writer http.ResponseWriter, request *http.Request, item database.FindItemShoppingRow) {
	var returnPantry SuccessErrorResponse
	updateItem := database.UpdateItemShoppingParams{
		Quantity: item.Quantity,
		ItemName: item.ItemName,
	}
	errUpdate := config.Db.UpdateItemShopping(request.Context(), updateItem)
	if errUpdate != nil {
		returnPantry.ErrorMessage = "Invalid quantity"
		writer.WriteHeader(http.StatusBadRequest)
		_ = config.Renderer.Render(writer, "ResponseMessage", returnPantry)
		return
	}
	returnPantry.SuccessMessage = fmt.Sprintf("Added %d items of %s", updateItem.Quantity, updateItem.ItemName)
	_ = config.Renderer.Render(writer, "ResponseMessage", returnPantry)
}

func (config *Config) RemoveItemShopping(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusBadRequest)
}
