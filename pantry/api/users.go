package api

import (
	"encoding/json"
	"log"
	"net/http"
	"pantry-pal/pantry/database"
	"time"

	"context"

	"github.com/google/uuid"
)

func (config *Config) ResetUsers(writer http.ResponseWriter, request *http.Request) {

	user, errUser := GetUserFromToken(request, writer, config)
	if errUser != nil {
		respondWithError(writer, http.StatusUnauthorized, errUser.Error())
		return
	}
	if config.Env != "dev" {
		respondWithError(writer, http.StatusUnauthorized, "Unable to perform this action in this environment")
		return
	}
	if user.IsAdmin.Int64 == 0 {
		respondWithError(writer, http.StatusUnauthorized, "User is not admin")
		return
	}

	errReset := config.Db.ResetTable(request.Context())
	if errReset != nil {
		respondWithError(writer, http.StatusInternalServerError, errReset.Error())
		return
	}

	respondWithJSON(writer, http.StatusAccepted, []byte{})

}

func GetUserFromToken(request *http.Request, writer http.ResponseWriter, config *Config) (database.User, error) {
	token, errTk := GetBearerToken(request.Header)
	log.Println("Token from header:", token)
	if errTk != nil {
		return database.User{}, errTk
	}

	userID, errJWT := ValidateJWT(token, config.Secret)
	if errJWT != nil {
		return database.User{}, errJWT
	}

	user, errUser := config.Db.GetUserByIdOrEmail(request.Context(), userID)
	if errUser != nil {
		return database.User{}, errUser
	}
	return user, nil
}

func (config *Config) GetUserInfo(writer http.ResponseWriter, request *http.Request) {

	info := request.PathValue("userInfo")
	log.Println("Trying to get data from users with:", info)
	user, errUser := config.Db.GetUserByIdOrEmail(request.Context(), info)

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

	hashedPassword, errPwd := HashPassword(user.Password)
	if errPwd != nil {
		respondWithError(writer, http.StatusInternalServerError, errPwd.Error())
	}
	createUser := database.CreateUserParams{
		ID:           uuid.New().String(),
		Email:        user.Email,
		Name:         user.Name,
		PasswordHash: hashedPassword,
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

func UpdateUser[T interface{}](writer http.ResponseWriter, request *http.Request, dbFunc func(context.Context, T) error) {
	var updateParams T

	err := json.NewDecoder(request.Body).Decode(&updateParams)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, err.Error())
		return
	}

	switch params := any(updateParams).(type) {
	case database.UpdateUserPasswordParams:
		hashedPassword, errPWD := HashPassword(params.PasswordHash)
		if errPWD != nil {
			respondWithError(writer, http.StatusInternalServerError, errPWD.Error())
			return
		}
		params.PasswordHash = hashedPassword
		updateParams = any(params).(T)
	}
	errUpdate := dbFunc(request.Context(), updateParams)

	if errUpdate != nil {
		respondWithError(writer, http.StatusInternalServerError, errUpdate.Error())
		return
	}

	respondWithJSON(writer, http.StatusAccepted, []byte{})
}

func (config *Config) UpdateUserEmail(writer http.ResponseWriter, request *http.Request) {
	UpdateUser(writer, request, config.Db.UpdateUserEmail)
}

func (config *Config) UpdateUserName(writer http.ResponseWriter, request *http.Request) {
	UpdateUser(writer, request, config.Db.UpdateUserName)
}

func (config *Config) UpdateUserPassword(writer http.ResponseWriter, request *http.Request) {
	UpdateUser(writer, request, config.Db.UpdateUserPassword)
}
