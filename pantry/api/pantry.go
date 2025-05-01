package api

import (
	"encoding/json"
	"net/http"
	"pantry-pal/pantry/database"

	"github.com/google/uuid"
)

func (config *Config) HandleNewItem(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	item := ItemAdd{}
	err := decoder.Decode(&item)

	if err != nil {
		respondWithError(writer, http.StatusBadRequest, err.Error())
	}

	findItem := database.FindItemByNameParams{
		UserID:   item.UserID,
		ItemName: item.ItemName,
	}
	items, errItem := config.Db.FindItemByName(request.Context(), findItem)
	if errItem != nil {
		respondWithError(writer, http.StatusInternalServerError, errItem.Error())
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
			config.UpdateItem(writer, request, toUpdate)
		} else {
			config.AddItem(writer, request, item)
		}
	}

}

func (config *Config) UpdateItem(writer http.ResponseWriter, request *http.Request, toUpdate ItemUpdate) {

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

func (config *Config) AddItem(writer http.ResponseWriter, request *http.Request, toAdd ItemAdd) {
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
	updateReponse := AddItemResponse{
		UserID:   addedItem.UserID,
		ItemName: addedItem.ItemName,
		Quantity: addedItem.Quantity,
		ExpiryAt: addedItem.ExpiryAt,
	}
	data, err := json.Marshal(updateReponse)
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
	_, errItem := config.Db.FindItemByName(request.Context(), findItem)

	if errItem != nil {
		respondWithError(writer, http.StatusInternalServerError, errItem.Error())
	}

}
