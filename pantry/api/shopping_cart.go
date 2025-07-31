package api

import (
	"context"
	"fmt"
	"net/http"
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

	cartItems, errCart := config.Db.GetAllShopping(context.Background(), userID)
	if errCart != nil {
		cartInfo.ErrorMessage = fmt.Sprintf("Unable to get current user cart info. Error: %s", errUser.Error())
		writer.WriteHeader(http.StatusForbidden)
		_ = config.Renderer.Render(writer, "ResponseMessage", cartInfo)

	}

	cartInfo.CartItems = cartItems
	_ = config.Renderer.Render(writer, "shoppingCartBlock", cartInfo)
	writer.WriteHeader(http.StatusOK)
}

func (config *Config) AddItemShopping(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusBadRequest)
}

func (config *Config) RemoveItemShopping(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusBadRequest)
}
