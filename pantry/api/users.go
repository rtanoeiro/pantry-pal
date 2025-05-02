package api

import (
	"encoding/json"
	"log"
	"net/http"
	"pantry-pal/pantry/database"
	"time"

	"github.com/google/uuid"
)

func (config *Config) ResetUsers(writer http.ResponseWriter, request *http.Request) {
	if config.Env != "dev" {
		respondWithError(writer, http.StatusUnauthorized, "Unable to perform this action in this environment")
		return
	}
	errReset := config.Db.ResetTable(request.Context())

	if errReset != nil {
		respondWithError(writer, http.StatusInternalServerError, errReset.Error())
		return
	}

	respondWithJSON(writer, http.StatusAccepted, []byte{})

}

func (config *Config) GetUserInfo(writer http.ResponseWriter, request *http.Request) {

	info := request.PathValue("userInfo")
	log.Println("Trying to get data from users with:", info)
	infoType := IdOrEmail(info)

	var user database.User
	var errUser error

	if infoType == "ID" {
		user, errUser = config.Db.GetUserById(request.Context(), info)
	} else {
		user, errUser = config.Db.GetUserByEmail(request.Context(), info)
	}

	if errUser != nil {
		respondWithError(writer, http.StatusBadRequest, errUser.Error())
		return
	}

	returnUser := CreateUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	data, err := json.Marshal(returnUser)

	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(writer, http.StatusCreated, data)

}

func (config *Config) CreateUser(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	user := UserAdd{}
	err := decoder.Decode(&user)

	if err != nil {
		respondWithError(writer, http.StatusBadRequest, err.Error())
	}

	_, userError := config.Db.GetUserByEmail(request.Context(), user.Email)

	if userError == nil {
		respondWithError(writer, http.StatusInternalServerError, "User already exists")
		return
	}

	createUser := database.CreateUserParams{
		ID:           uuid.New().String(),
		Email:        user.Email,
		Name:         user.Name,
		PasswordHash: user.Password,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	userAdd, errAdd := config.Db.CreateUser(request.Context(), createUser)
	if errAdd != nil {
		respondWithError(writer, http.StatusInternalServerError, errAdd.Error())
		return
	}
	log.Printf("User added with success - UserID:%s \n-Name:%s\n-Email: %s", userAdd.ID, userAdd.Name, userAdd.Email)

	returnUser := CreateUserResponse{
		ID:    userAdd.ID,
		Name:  userAdd.Name,
		Email: userAdd.Email,
	}
	data, err := json.Marshal(returnUser)

	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(writer, http.StatusCreated, data)

}
