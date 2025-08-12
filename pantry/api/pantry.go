package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pantry-pal/pantry/database"

	"github.com/google/uuid"
)

func (config *Config) HandleNewItem(writer http.ResponseWriter, request *http.Request) {
	var returnPantry SuccessErrorResponse
	itemName := request.FormValue("itemName")
	itemName = strings.TrimSpace(itemName)

	itemQuantity, errQuantity := strconv.Atoi(request.FormValue("itemQuantity"))
	if errQuantity != nil {
		returnPantry.ErrorMessage = "Invalid quantity"
		writer.WriteHeader(http.StatusBadRequest)
		_ = config.Renderer.Render(writer, "HomeResponseMessage", returnPantry)
		return
	}
	itemExpiry := request.FormValue("itemExpiryDate")

	if !ValidateDate(itemExpiry) {
		returnPantry.ErrorMessage = "Invalid Date. Please send in the Format YYYY-MM-DD or Date is already expired"
		writer.WriteHeader(http.StatusBadRequest)
		_ = config.Renderer.Render(writer, "HomeResponseMessage", returnPantry)
		return
	}

	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		returnPantry.ErrorMessage = fmt.Sprintf("Unable to retrieve user Pantry Items. Error: %s", errUser.Error())
		_ = config.Renderer.Render(writer, "HomeResponseMessage", returnPantry)
		return
	}

	findItem := database.FindItemByNameParams{
		UserID:   userID,
		ItemName: strings.ToLower(itemName),
	}
	items, errItem := config.Db.FindItemByName(request.Context(), findItem)
	if errItem != nil {
		returnPantry.ErrorMessage = "Failed to get current items in pantry, please try again"
		_ = config.Renderer.Render(writer, "HomeResponseMessage", returnPantry)
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
}

func (config *Config) ItemUpdate(
	writer http.ResponseWriter,
	request *http.Request,
	toUpdate UpdateItemRequest,
) {
	log.Printf("User %s is trying to add %d of %s items into pantry at %s.", toUpdate.UserID, toUpdate.QuantityToAdd, toUpdate.ItemName, time.Now())
	var returnPantry SuccessErrorResponse

	itemToUpdate := database.UpdateItemQuantityParams{
		Quantity: toUpdate.QuantityAvailable + toUpdate.QuantityToAdd,
		ID:       toUpdate.ItemID,
		UserID:   toUpdate.UserID,
	}

	if itemToUpdate.Quantity <= 0 {
		log.Printf("User %s trying to delete %s item from pantry at %s.", toUpdate.UserID, toUpdate.ItemName, time.Now())
		request.Form.Add("ItemID", toUpdate.ItemID)
		config.DeleteItem(writer, request)
		return
	}
	errUpdate := config.Db.UpdateItemQuantity(request.Context(), itemToUpdate)
	if errUpdate != nil {
		log.Printf("User %s failed to add %d of %s items into pantry at %s. Failed update", toUpdate.UserID, toUpdate.QuantityToAdd, toUpdate.ItemName, time.Now())
		returnPantry.ErrorMessage = fmt.Sprintf("Failed to update items to Pantry, please try again. Error: %s", errUpdate.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		_ = config.Renderer.Render(writer, "HomeResponseMessage", returnPantry)
		return
	}
	returnPantry.SuccessMessage = fmt.Sprintf("%s updated on Pantry", toUpdate.ItemName)
	writer.WriteHeader(http.StatusOK)
	log.Printf("User %s successfully updated %d of %s items into pantry at %s.", toUpdate.UserID, toUpdate.QuantityToAdd, toUpdate.ItemName, time.Now())
	_ = config.Renderer.Render(writer, "HomeResponseMessage", returnPantry)
}

func (config *Config) HandleAddOnePantry(writer http.ResponseWriter, request *http.Request) {
	var returnPantry SuccessErrorResponse
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		returnPantry.ErrorMessage = fmt.Sprintf("Unable to retrieve user Pantry Items. Error: %s", errUser.Error())
		_ = config.Renderer.Render(writer, "HomeResponseMessage", returnPantry)
		return
	}
	itemID := request.PathValue("ItemID")
	findItem := database.FindItemByIDParams{UserID: userID, ID: itemID}
	log.Printf("User is trying to add one of ItemID %s", itemID)
	currentItem, errItem := config.Db.FindItemByID(request.Context(), findItem)
	if errItem != nil {
		returnPantry.ErrorMessage = "Failed to get current items in pantry, please try again"
		_ = config.Renderer.Render(writer, "HomeResponseMessage", returnPantry)
		return
	}

	toUpdate := UpdateItemRequest{
		ItemID:            itemID,
		UserID:            userID,
		ItemName:          currentItem.ItemName,
		QuantityAvailable: currentItem.Quantity,
		QuantityToAdd:     1,
		ExpiryAt:          currentItem.ExpiryAt,
	}
	config.ItemUpdate(writer, request, toUpdate)
}

func (config *Config) HandleRemoveOnePantry(writer http.ResponseWriter, request *http.Request) {
	var returnPantry SuccessErrorResponse
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		returnPantry.ErrorMessage = fmt.Sprintf("Unable to retrieve user Pantry Items. Error: %s", errUser.Error())
		_ = config.Renderer.Render(writer, "HomeResponseMessage", returnPantry)
		return
	}
	itemID := request.PathValue("ItemID")
	findItem := database.FindItemByIDParams{UserID: userID, ID: itemID}
	log.Printf("User is trying to remove one of ItemID %s", itemID)
	currentItem, errItem := config.Db.FindItemByID(request.Context(), findItem)
	if errItem != nil {
		returnPantry.ErrorMessage = "Failed to get current items in pantry, please try again"
		_ = config.Renderer.Render(writer, "HomeResponseMessage", returnPantry)
		return
	}

	toUpdate := UpdateItemRequest{
		ItemID:            itemID,
		UserID:            userID,
		ItemName:          currentItem.ItemName,
		QuantityAvailable: currentItem.Quantity,
		QuantityToAdd:     -1,
		ExpiryAt:          currentItem.ExpiryAt,
	}
	config.ItemUpdate(writer, request, toUpdate)
}

func (config *Config) ItemAdd(
	writer http.ResponseWriter,
	request *http.Request,
	toAdd AddItemRequest,
) {
	log.Println("Adding item to pantry")
	var returnPantry SuccessErrorResponse
	if toAdd.Quantity < 0 {
		log.Printf("User %s failed to add %d of %s items into pantry at %s. Negative quantity", toAdd.UserID, toAdd.Quantity, toAdd.ItemName, time.Now())
		returnPantry.ErrorMessage = "Unable to add negative items"
		writer.WriteHeader(http.StatusBadRequest)
		_ = config.Renderer.Render(writer, "HomeResponseMessage", returnPantry)
		return
	}
	itemToAdd := database.AddItemParams{
		ID:       uuid.NewString(),
		UserID:   toAdd.UserID,
		ItemName: toAdd.ItemName,
		Quantity: toAdd.Quantity,
		ExpiryAt: toAdd.ExpiryAt,
	}

	addedItem, errUpdate := config.Db.AddItem(request.Context(), itemToAdd)
	if errUpdate != nil {
		log.Printf("User %s failed to add %d of %s items into pantry at %s. Server error", toAdd.UserID, toAdd.Quantity, toAdd.ItemName, time.Now())
		returnPantry.ErrorMessage = fmt.Sprintf("Failed to add items to Pantry, please try again. Error: %s", errUpdate.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		_ = config.Renderer.Render(writer, "HomeResponseMessage", returnPantry)
		return
	}
	returnPantry.SuccessMessage = fmt.Sprintf("%s added to pantry", addedItem.ItemName)
	log.Printf("User %s successfully added %d of %s items into pantry at %s.", toAdd.UserID, toAdd.Quantity, toAdd.ItemName, time.Now())
	_ = config.Renderer.Render(writer, "HomeResponseMessage", returnPantry)
	writer.WriteHeader(http.StatusOK)
}

// TODO: RIght now when hitting this API, we're returning only portion of the HTML, which is causing CSS issues when rendering it.
// TODO: Adjust how data is returned when rendering this, so the whole page load doesn't fail.
func (config *Config) RenderPantryPage(writer http.ResponseWriter, request *http.Request) {
	var returnPantry SuccessErrorResponse
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		log.Printf("Unable to retrieve pantry items from user %s at %s. Error on User Token", userID, time.Now())
		returnPantry.ErrorMessage = fmt.Sprintf("Unable to retrieve user Pantry Items. Error: %s", errUser.Error())
		writer.WriteHeader(http.StatusForbidden)
		_ = config.Renderer.Render(writer, "pantry", returnPantry)
		return
	}
	writer.WriteHeader(http.StatusOK)
	_ = config.Renderer.Render(writer, "pantry", returnPantry)
}

func (config *Config) GetPantryItems(writer http.ResponseWriter, request *http.Request) {
	var returnPantry PantryItems
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		log.Printf("Unable to retrieve pantry items from user %s at %s. Error on User Token", userID, time.Now())
		returnPantry.ErrorMessage = fmt.Sprintf("Unable to retrieve user Pantry Items. Error: %s", errUser.Error())
		writer.WriteHeader(http.StatusForbidden)
		_ = config.Renderer.Render(writer, "pantryItems", returnPantry)
		return
	}

	allPantryItems, errAll := config.Db.GetAllItems(request.Context(), userID)
	if errAll != nil {
		log.Printf("Unable to retrieve pantry items from user %s at %s. Error on getting items", userID, time.Now())
		returnPantry.ErrorMessage = fmt.Sprintf("Unable to retrieve user Pantry Items. Error: %s", errAll.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		_ = config.Renderer.Render(writer, "pantryItems", returnPantry)
		return
	}

	var PantrySlice []PantryItem
	for _, item := range allPantryItems {
		toAppend := PantryItem{
			ItemID:   item.ID,
			ItemName: item.ItemName,
			Quantity: int(item.Quantity),
			ExpiryAt: item.ExpiryAt,
		}
		PantrySlice = append(PantrySlice, toAppend)
	}
	returnPantry.Items = PantrySlice
	writer.WriteHeader(http.StatusOK)
	_ = config.Renderer.Render(writer, "pantryItems", returnPantry)

}

func (config *Config) DeleteItem(writer http.ResponseWriter, request *http.Request) {
	var returnMessage SuccessErrorResponse
	itemID := request.PathValue("ItemID")
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		log.Printf("Error on deleting item from user %s at %s. Invalid token", userID, time.Now())
		returnMessage.ErrorMessage = fmt.Sprintf("Unable to retrieve user Pantry Items. Error: %s", errUser.Error())
		writer.WriteHeader(http.StatusForbidden)
		_ = config.Renderer.Render(writer, "pantry", returnMessage)
		return
	}
	log.Printf("User %s is trying to remove %s from pantry at %s.", userID, itemID, time.Now())

	removeParams := database.RemoveItemParams{
		ID:     itemID,
		UserID: userID,
	}
	item, errRemove := config.Db.RemoveItem(request.Context(), removeParams)

	if errRemove != nil {
		log.Printf("Error on deleting item from user %s at %s. Error: %s", userID, time.Now(), errRemove.Error())
		returnMessage.ErrorMessage = fmt.Sprintf("Error on deleting item . Error: %s", errRemove.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		_ = config.Renderer.Render(writer, "pantry", returnMessage)
		return
	}
	log.Printf("User %s removed x%d from ItemID %s at %s", userID, item.Quantity, removeParams.ID, time.Now())
	returnMessage.SuccessMessage = fmt.Sprintf("Successfulyl deleted x%d %s: ", item.Quantity, item.ItemName)
	writer.WriteHeader(http.StatusOK)

	_ = config.Renderer.Render(writer, "pantry", returnMessage)
}

func (config *Config) RenderExpiringSoon(writer http.ResponseWriter, request *http.Request) {
	var pantryItems PantryStats
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		log.Printf("Unable to get expiring soon items. unauthorised User ID at %s", time.Now())
		pantryItems.ErrorMessage = fmt.Sprintf("Unable to retrieve expiring soon items. Error: %s", errUser.Error())
		writer.WriteHeader(http.StatusForbidden)
		_ = config.Renderer.Render(writer, "expiringSoonBlock", pantryItems)
		return
	}

	expiringSoon, errExpiring := config.Db.GetExpiringSoon(request.Context(), userID)
	if errExpiring != nil {
		log.Printf("Unable to get expiring soon items. Failed to read data from database at %s", time.Now())
		_ = config.Renderer.Render(writer, "expiringSoonBlock", pantryItems)
		writer.WriteHeader(http.StatusInternalServerError)
		pantryItems.ErrorMessage = fmt.Sprintf("Unable to retrieve expiring soon items. Error: %s", errExpiring.Error())
		return
	}

	expiringSoonItems := make([]PantryItem, len(expiringSoon))
	for index, item := range expiringSoon {
		expiringSoonItems[index] = PantryItem{
			ItemID:   item.ID,
			ItemName: item.ItemName,
			Quantity: int(item.Quantity),
			ExpiryAt: item.ExpiryAt,
		}
	}
	log.Printf("Loading all expiring items for user %s", userID)
	pantryItems.ExpiringSoon = expiringSoonItems
	writer.WriteHeader(http.StatusOK)
	_ = config.Renderer.Render(writer, "expiringSoonBlock", pantryItems)
}
