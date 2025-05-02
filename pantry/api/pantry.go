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

	decoder := json.NewDecoder(request.Body)
	item := ItemAdd{}
	err := decoder.Decode(&item)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("Received item \n- Name: %s\n- Quantity: %d", item.ItemName, item.Quantity)

	findItem := database.FindItemByNameParams{
		UserID:   item.UserID,
		ItemName: strings.ToLower(item.ItemName),
	}
	items, errItem := config.Db.FindItemByName(request.Context(), findItem)
	if errItem != nil {
		respondWithError(writer, http.StatusInternalServerError, errItem.Error())
		return
	}
	for _, item := range items {
		log.Printf("Found item \n- Name: %s\n- Quantity: %d", item.ItemName, item.Quantity)
	}

	for _, currentItem := range items {
		if currentItem.ItemName == item.ItemName && currentItem.ExpiryAt == item.ExpiryAt {
			toUpdate := ItemUpdate{
				ItemID:            currentItem.ID,
				UserID:            item.UserID,
				ItemName:          item.ItemName,
				QuantityAvailable: currentItem.Quantity,
				QuantityToAdd:     item.Quantity,
				ExpiryAt:          currentItem.ExpiryAt,
			}
			config.ItemUpdate(writer, request, toUpdate)
			break
		}
	}
	// if the function hasn't returned yet, the item is new, so we add it
	config.ItemAdd(writer, request, item)

}

func (config *Config) ItemUpdate(writer http.ResponseWriter, request *http.Request, toUpdate ItemUpdate) {

	itemToUpdate := database.UpdateItemQuantityParams{
		Quantity: toUpdate.QuantityAvailable + toUpdate.QuantityToAdd,
		ID:       toUpdate.ItemID,
		UserID:   toUpdate.UserID,
	}
	updatedItem, errUpdate := config.Db.UpdateItemQuantity(request.Context(), itemToUpdate)
	if errUpdate != nil {
		respondWithError(writer, http.StatusInternalServerError, errUpdate.Error())
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
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, data)

}

func (config *Config) ItemAdd(writer http.ResponseWriter, request *http.Request, toAdd ItemAdd) {

	itemToAdd := database.AddItemParams{
		ID:       uuid.NewString(),
		UserID:   toAdd.UserID,
		ItemName: toAdd.ItemName,
		Quantity: toAdd.Quantity,
		ExpiryAt: toAdd.ExpiryAt,
	}
	addedItem, errUpdate := config.Db.AddItem(request.Context(), itemToAdd)
	if errUpdate != nil {
		respondWithError(writer, http.StatusInternalServerError, errUpdate.Error())
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
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, data)
}

func (config *Config) GetItemByName(writer http.ResponseWriter, request *http.Request) {
	itemName := request.PathValue("itemName")

	decoder := json.NewDecoder(request.Body)
	user := UserRequests{}
	err := decoder.Decode(&user)

	if err != nil {
		respondWithError(writer, http.StatusBadRequest, err.Error())
	}

	findItem := database.FindItemByNameParams{
		UserID:   user.ID,
		ItemName: itemName,
	}
	items, errItem := config.Db.FindItemByName(request.Context(), findItem)

	if errItem != nil {
		respondWithError(writer, http.StatusInternalServerError, errItem.Error())
	}

	data, err := json.Marshal(items)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, data)
}

func (config *Config) GetAllPantryItems(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	user := UserRequests{}
	err := decoder.Decode(&user)

	if err != nil {
		respondWithError(writer, http.StatusBadRequest, err.Error())
		return
	}

	allPantryItems, errAll := config.Db.GetAllItems(request.Context(), user.ID)

	if errAll != nil {
		respondWithError(writer, http.StatusBadRequest, errAll.Error())
		return
	}

	for _, item := range allPantryItems {
		log.Printf("Found item \n- Name: %s\n- Quantity: %d", item.ItemName, item.Quantity)
	}
	data, err := json.Marshal(allPantryItems)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, data)

}

func (config *Config) DeleteItem(writer http.ResponseWriter, request *http.Request) {
	itemID := request.PathValue("itemID")
	log.Println("Trying to remove: ", itemID)

	item, errRemove := config.Db.RemoveItem(request.Context(), itemID)

	if errRemove != nil {
		respondWithError(writer, http.StatusInternalServerError, errRemove.Error())
		return
	}
	log.Printf("Successfully remove %d item(s) %s, with Expiry date at %s", item.Quantity, item.ItemName, item.ExpiryAt)
	respondWithJSON(writer, http.StatusOK, []byte{})
}
