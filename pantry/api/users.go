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

	// Update to check for admin privileges?
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

func GetUserIDFromToken(request *http.Request, writer http.ResponseWriter, config *Config) (string, error) {
	token, errTk := GetBearerToken(request.Header)
	log.Println("Token from header:", token)
	if errTk != nil {
		return "", errTk
	}

	userID, errJWT := ValidateJWT(token, config.Secret)
	if errJWT != nil {
		return "", errJWT
	}
	log.Println("User ID from token:", userID)

	return userID, nil
}

func (config *Config) GetUserInfo(writer http.ResponseWriter, request *http.Request) {

	userID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithError(writer, http.StatusUnauthorized, errUser.Error())
		return
	}

	userData, errUser := config.Db.GetUserByIdOrEmail(request.Context(), userID)

	if errUser != nil {
		respondWithError(writer, http.StatusBadRequest, errUser.Error())
		return
	}

	returnUser := CreateUserResponse{
		ID:    userData.ID,
		Name:  userData.Name,
		Email: userData.Email,
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
	user := CreateUserRequest{}
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
		return
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

func UpdateUser[T interface{}](writer http.ResponseWriter, request *http.Request, dbFunc func(context.Context, T) error, config *Config) {

	userID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithError(writer, http.StatusUnauthorized, errUser.Error())
		return
	}
	log.Println("User ID from token on Update User Call:", userID)

	var updateParams T
	err := json.NewDecoder(request.Body).Decode(&updateParams)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, err.Error())
		return
	}

	switch params := any(updateParams).(type) {
	case database.UpdateUserPasswordParams:
		log.Println("Updating user password")
		hashedPassword, errPWD := HashPassword(params.PasswordHash)
		if errPWD != nil {
			respondWithError(writer, http.StatusInternalServerError, errPWD.Error())
			return
		}
		params.PasswordHash = hashedPassword
		params.ID = userID
		updateParams = any(params).(T)
	case database.UpdateUserEmailParams:
		log.Println("Updating user email")
		params.ID = userID
		updateParams = any(params).(T)
	case database.UpdateUserNameParams:
		log.Println("Updating user name")
		params.ID = userID
		updateParams = any(params).(T)
	default:
		log.Println("Wrong parameters in request")
		respondWithError(writer, http.StatusMethodNotAllowed, "Wrong Parameters in Request")
	}
	errUpdate := dbFunc(request.Context(), updateParams)

	if errUpdate != nil {
		respondWithError(writer, http.StatusInternalServerError, errUpdate.Error())
		return
	}

	respondWithJSON(writer, http.StatusAccepted, []byte{})
}

func (config *Config) UpdateUserEmail(writer http.ResponseWriter, request *http.Request) {
	UpdateUser(writer, request, config.Db.UpdateUserEmail, config)
}

func (config *Config) UpdateUserName(writer http.ResponseWriter, request *http.Request) {
	UpdateUser(writer, request, config.Db.UpdateUserName, config)
}

func (config *Config) UpdateUserPassword(writer http.ResponseWriter, request *http.Request) {
	UpdateUser(writer, request, config.Db.UpdateUserPassword, config)
}

func (config *Config) UserAdmin(writer http.ResponseWriter, request *http.Request) {
	log.Println("User Admin endpoint called")
	userID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithError(writer, http.StatusUnauthorized, errUser.Error())
		return
	}

	log.Println("User ID from token:", userID)
	errADmin := config.Db.UpdateUserAdmin(request.Context(), userID)
	if errADmin != nil {
		respondWithError(writer, http.StatusInternalServerError, errADmin.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, []byte("User successfully upgraded to Admin"))
}
