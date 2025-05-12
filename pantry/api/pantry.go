package api

import (
	"encoding/json"
	"log"
	"net/http"
	"pantry-pal/pantry/database"
	"strings"

	"github.com/google/uuid"
)

func (config *Config) HandleNewItem(writer http.ResponseWriter, request *http.Request) {

	userID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return
	}

	decoder := json.NewDecoder(request.Body)
	item := AddItemRequest{}
	err := decoder.Decode(&item)
	if err != nil {
		respondWithJSON(writer, http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("Received item \n- Name: %s\n- Quantity: %d", item.ItemName, item.Quantity)

	findItem := database.FindItemByNameParams{
		UserID:   userID,
		ItemName: strings.ToLower(item.ItemName),
	}
	items, errItem := config.Db.FindItemByName(request.Context(), findItem)
	if errItem != nil {
		respondWithJSON(writer, http.StatusInternalServerError, errItem.Error())
		return
	}
	for _, item := range items {
		log.Printf("Found item \n- Name: %s\n- Quantity: %d", item.ItemName, item.Quantity)
	}

	for _, currentItem := range items {
		if currentItem.ItemName == item.ItemName && currentItem.ExpiryAt == item.ExpiryAt {
			toUpdate := UpdateItemRequest{
				ItemID:            currentItem.ID,
				UserID:            userID,
				ItemName:          item.ItemName,
				QuantityAvailable: currentItem.Quantity,
				QuantityToAdd:     item.Quantity,
				ExpiryAt:          currentItem.ExpiryAt,
			}
			config.ItemUpdate(writer, request, toUpdate)
			return
		}
	}
	// if the function hasn't returned yet, the item is new, so we add it
	item.UserID = userID
	config.ItemAdd(writer, request, item)

}

func (config *Config) ItemUpdate(writer http.ResponseWriter, request *http.Request, toUpdate UpdateItemRequest) {
	log.Println("Updating item in pantry")
	itemToUpdate := database.UpdateItemQuantityParams{
		Quantity: toUpdate.QuantityAvailable + toUpdate.QuantityToAdd,
		ID:       toUpdate.ItemID,
		UserID:   toUpdate.UserID,
	}

	if itemToUpdate.Quantity < 0 {
		respondWithJSON(writer, http.StatusForbidden, "unable to remove more items than available")
		return
	}
	updatedItem, errUpdate := config.Db.UpdateItemQuantity(request.Context(), itemToUpdate)
	if errUpdate != nil {
		respondWithJSON(writer, http.StatusInternalServerError, errUpdate.Error())
		return
	}
	updateReponse := UpdateItemResponse{
		ItemID:   updatedItem.ID,
		UserID:   toUpdate.UserID,
		ItemName: toUpdate.ItemName,
		Quantity: updatedItem.Quantity,
		ExpiryAt: toUpdate.ExpiryAt,
	}
	data, err := json.Marshal(updateReponse)
	if err != nil {
		respondWithJSON(writer, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, data)

}

func (config *Config) ItemAdd(writer http.ResponseWriter, request *http.Request, toAdd AddItemRequest) {
	log.Println("Adding item to pantry")

	if toAdd.Quantity < 0 {
		respondWithJSON(writer, http.StatusForbidden, "unable to (add negative items")
		return
	}
	itemToAdd := database.AddItemParams{
		ID:       uuid.NewString(),
		UserID:   toAdd.UserID,
		ItemName: toAdd.ItemName,
		Quantity: toAdd.Quantity,
		ExpiryAt: toAdd.ExpiryAt,
	}

	if !checkDate(toAdd.ExpiryAt) {
		respondWithJSON(writer, http.StatusForbidden, "Invalid Date. Please send in the Format YYYY-MM-DD or Date is already expired")
		return
	}
	log.Printf("Item to add: \n- UserID: %s \n- ItemName: %s \n- Quantity: %d \n- Expiry Date: %s", toAdd.UserID, toAdd.ItemName, toAdd.Quantity, toAdd.ExpiryAt)
	addedItem, errUpdate := config.Db.AddItem(request.Context(), itemToAdd)
	if errUpdate != nil {
		respondWithJSON(writer, http.StatusInternalServerError, errUpdate.Error())
		return
	}

	addResponse := AddItemResponse{
		ItemID:   addedItem.ID,
		UserID:   toAdd.UserID,
		ItemName: addedItem.ItemName,
		Quantity: addedItem.Quantity,
		ExpiryAt: addedItem.ExpiryAt,
	}

	data, err := json.Marshal(addResponse)
	if err != nil {
		respondWithJSON(writer, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, data)
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

	userID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return
	}

	allPantryItems, errAll := config.Db.GetAllItems(request.Context(), userID)
	if errAll != nil {
		respondWithJSON(writer, http.StatusBadRequest, errAll.Error())
		return
	}

	for _, item := range allPantryItems {
		log.Printf("Found item \n- Name: %s\n- Quantity: %d", item.ItemName, item.Quantity)
	}
	data, err := json.Marshal(allPantryItems)
	if err != nil {
		respondWithJSON(writer, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, data)

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
	pantryStats := PantryStats{
		NumItems:     int(totalItems),
		ExpiringSoon: expiringSoonItems,
		ShoppingList: []ItemShopping{},
	}
	log.Println("Current Pantry Stats: ", pantryStats)

	data, err := json.Marshal(pantryStats)
	if err != nil {
		respondWithJSON(writer, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, data)
}
