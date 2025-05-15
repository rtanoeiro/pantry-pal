package api

import (
	"encoding/json"
	"log"
	"net/http"
	"pantry-pal/pantry/database"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func (config *Config) HandleNewItem(writer http.ResponseWriter, request *http.Request) {

	itemName := request.FormValue("itemName")
	itemName = strings.TrimSpace(itemName)
	itemQuantity, errQuantity := strconv.Atoi(request.FormValue("itemQuantity"))
	if errQuantity != nil {
		respondWithJSON(writer, http.StatusBadRequest, "Invalid quantity")
		return
	}
	itemExpiry := request.FormValue("itemExpiryDate")
	log.Printf("Received item \n- Name: %s\n- Quantity: %d\n- Expity Date: %s", itemName, itemQuantity, itemExpiry)
	returnPantry := map[string]interface{}{
		"ErrorMessage":   "",
		"SuccessMessage": "",
	}

	if !checkDate(itemExpiry) {
		respondWithJSON(writer, http.StatusForbidden, "Invalid Date. Please send in the Format YYYY-MM-DD or Date is already expired")
		returnPantry["ErrorMessage"] = "Invalid Date. Please send in the Format YYYY-MM-DD or Date is already expired"
		config.Renderer.Render(writer, "ResponseMessage", returnPantry)
		return
	}

	userID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		returnPantry["ErrorMessage"] = "Unable to retrieve user Pantry Items"
		config.Renderer.Render(writer, "ResponseMessage", returnPantry)
		return
	}

	findItem := database.FindItemByNameParams{
		UserID:   userID,
		ItemName: strings.ToLower(itemName),
	}
	items, errItem := config.Db.FindItemByName(request.Context(), findItem)
	if errItem != nil {
		respondWithJSON(writer, http.StatusInternalServerError, errItem.Error())
		returnPantry["ErrorMessage"] = "Failed to get current items in pantry, please try again"
		config.Renderer.Render(writer, "ResponseMessage", returnPantry)
		return
	}

	for _, currentItem := range items {
		if currentItem.ItemName == itemName && currentItem.ExpiryAt == itemExpiry {
			toUpdate := UpdateItemRequest{
				ItemID:            currentItem.ID,
				UserID:            userID,
				ItemName:          itemName,
				QuantityAvailable: currentItem.Quantity,
				QuantityToAdd:     int64(itemQuantity),
				ExpiryAt:          currentItem.ExpiryAt,
			}
			config.ItemUpdate(writer, request, toUpdate)
			config.GetPantryStats(writer, request)
			return
		}
	}
	addItem := AddItemRequest{
		UserID:   userID,
		ItemName: itemName,
		Quantity: int64(itemQuantity),
		ExpiryAt: itemExpiry,
	}
	config.ItemAdd(writer, request, addItem)
	writer.Header().Set("HX-Redirect", "/home")
}

func (config *Config) ItemUpdate(writer http.ResponseWriter, request *http.Request, toUpdate UpdateItemRequest) {
	log.Println("Updating item in pantry")
	returnPantry := map[string]interface{}{
		"ErrorMessage":   "",
		"SuccessMessage": "",
	}
	itemToUpdate := database.UpdateItemQuantityParams{
		Quantity: toUpdate.QuantityAvailable + toUpdate.QuantityToAdd,
		ID:       toUpdate.ItemID,
		UserID:   toUpdate.UserID,
	}

	if itemToUpdate.Quantity < 0 {
		respondWithJSON(writer, http.StatusForbidden, "unable to remove more items than available")
		returnPantry["ErrorMessage"] = "Unable to remove more items than available"
		config.Renderer.Render(writer, "ResponseMessage", returnPantry)
		return
	}
	updatedItem, errUpdate := config.Db.UpdateItemQuantity(request.Context(), itemToUpdate)
	if errUpdate != nil {
		respondWithJSON(writer, http.StatusInternalServerError, errUpdate.Error())
		returnPantry["ErrorMessage"] = "Failed to update items to Pantry, please try again"
		config.Renderer.Render(writer, "ResponseMessage", returnPantry)
		return
	}
	returnPantry["SuccessMessage"] = updatedItem.ItemName + "updated on Pantry"
	config.Renderer.Render(writer, "ResponseMessage", returnPantry)
}

func (config *Config) ItemAdd(writer http.ResponseWriter, request *http.Request, toAdd AddItemRequest) {
	log.Println("Adding item to pantry")
	returnPantry := map[string]interface{}{
		"ErrorMessage":   "",
		"SuccessMessage": "",
	}
	if toAdd.Quantity < 0 {
		respondWithJSON(writer, http.StatusForbidden, "Unable to add negative items")
		returnPantry["ErrorMessage"] = "Unable to add negative items"
		config.Renderer.Render(writer, "ResponseMessage", returnPantry)
		return
	}
	itemToAdd := database.AddItemParams{
		ID:       uuid.NewString(),
		UserID:   toAdd.UserID,
		ItemName: toAdd.ItemName,
		Quantity: toAdd.Quantity,
		ExpiryAt: toAdd.ExpiryAt,
	}

	log.Printf("Item to add: \n- UserID: %s \n- ItemName: %s \n- Quantity: %d \n- Expiry Date: %s", toAdd.UserID, toAdd.ItemName, toAdd.Quantity, toAdd.ExpiryAt)
	addedItem, errUpdate := config.Db.AddItem(request.Context(), itemToAdd)
	if errUpdate != nil {
		respondWithJSON(writer, http.StatusInternalServerError, errUpdate.Error())
		returnPantry["ErrorMessage"] = "Failed to add items to Pantry, please try again"
		config.Renderer.Render(writer, "ResponseMessage", returnPantry)
		return
	}
	returnPantry["SuccessMessage"] = addedItem.ItemName + " added to pantry"
	config.Renderer.Render(writer, "ResponseMessage", returnPantry)
}

func (config *Config) GetItemByName(writer http.ResponseWriter, request *http.Request) {
	itemName := request.PathValue("itemName")

	userID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return
	}

	findItem := database.FindItemByNameParams{
		UserID:   userID,
		ItemName: itemName,
	}
	items, errItem := config.Db.FindItemByName(request.Context(), findItem)

	if errItem != nil {
		respondWithJSON(writer, http.StatusInternalServerError, errItem.Error())
	}

	data, err := json.Marshal(items)
	if err != nil {
		respondWithJSON(writer, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, data)
}

func (config *Config) GetAllPantryItems(writer http.ResponseWriter, request *http.Request) {

	log.Println("User entered it's Pantry")
	returnPantry := map[string]interface{}{
		"ErrorMessage":   "",
		"SuccessMessage": "",
		"Items":          []PantryItem{},
	}

	userID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		returnPantry["ErrorMessage"] = "Unable to retrieve user Pantry Items"
		config.Renderer.Render(writer, "pantrypage", returnPantry)
		return
	}

	allPantryItems, errAll := config.Db.GetAllItems(request.Context(), userID)
	if errAll != nil {
		respondWithJSON(writer, http.StatusBadRequest, errAll.Error())
		returnPantry["ErrorMessage"] = "Unable to retrieve user Pantry Items"
		config.Renderer.Render(writer, "pantrypage", returnPantry)
		return
	}

	var PantrySlice []PantryItem
	for _, item := range allPantryItems {
		log.Printf("Found item \n- Name: %s\n- Quantity: %d", item.ItemName, item.Quantity)
		toAppend := PantryItem{
			ItemName: item.ItemName,
			Quantity: int(item.Quantity),
			ExpiryAt: item.ExpiryAt,
		}
		PantrySlice = append(PantrySlice, toAppend)
	}
	returnPantry["Items"] = PantrySlice
	config.Renderer.Render(writer, "pantrypage", returnPantry)
}

func (config *Config) DeleteItem(writer http.ResponseWriter, request *http.Request) {

	userID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return
	}

	itemID := request.PathValue("itemID")
	log.Println("Trying to remove: ", itemID)

	removeParams := database.RemoveItemParams{
		ID:     itemID,
		UserID: userID,
	}
	item, errRemove := config.Db.RemoveItem(request.Context(), removeParams)

	if errRemove != nil {
		respondWithJSON(writer, http.StatusInternalServerError, errRemove.Error())
		return
	}
	log.Printf("Successfully remove %d item(s) %s, with Expiry date at %s", item.Quantity, item.ItemName, item.ExpiryAt)
	respondWithJSON(writer, http.StatusOK, []byte{})
}

func (config *Config) GetPantryStats(writer http.ResponseWriter, request *http.Request) {
	userID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return
	}

	totalItems, errItem := config.Db.GetTotalNumberOfItems(request.Context(), userID)
	if errItem != nil {
		respondWithJSON(writer, http.StatusInternalServerError, errItem.Error())
		return
	}
	expiringSoon, errExpiring := config.Db.GetExpiringSoon(request.Context(), userID)
	if errExpiring != nil {
		respondWithJSON(writer, http.StatusInternalServerError, errExpiring.Error())
		return
	}

	expiringSoonItems := make([]PantryItem, len(expiringSoon))
	for index, item := range expiringSoon {
		expiringSoonItems[index] = PantryItem{
			ItemName: item.ItemName,
			Quantity: int(item.Quantity),
			ExpiryAt: item.ExpiryAt,
		}
	}
	pantryStats := map[string]interface{}{
		"NumItems":     int(totalItems),
		"ExpiringSoon": expiringSoonItems,
		"ShoppingList": []ItemShopping{},
	}
	log.Println("Current Pantry Stats: ", pantryStats)
	config.Renderer.Render(writer, "pantryStats", pantryStats)
}
