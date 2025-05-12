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
		respondWithJSON(writer, http.StatusUnauthorized, "Unable to perform this action in this environment")
		return
	}

	errReset := config.Db.ResetTable(request.Context())
	if errReset != nil {
		respondWithJSON(writer, http.StatusInternalServerError, errReset.Error())
		return
	}

	respondWithJSON(writer, http.StatusAccepted, []byte{})

}

func GetUserIDFromToken(request *http.Request, writer http.ResponseWriter, config *Config) (string, error) {
	token, errTk := GetJWTFromCookie(request)
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

	log.Println("User Info endpoint called")
	userID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return
	}
	writer.Header().Set("HX-Replace-Url", "/user")
	userData, errUser := config.Db.GetUserById(request.Context(), userID)

	if errUser != nil {
		respondWithJSON(writer, http.StatusBadRequest, errUser.Error())
		return
	}

	returnUser := CreateUserResponse{
		ID:    userData.ID,
		Name:  userData.Name,
		Email: userData.Email,
	}
	config.Renderer.Render(writer, "user", returnUser)

}

func (config *Config) CreateUser(writer http.ResponseWriter, request *http.Request) {

	email := request.FormValue("email")
	name := request.FormValue("name")
	password := request.FormValue("password")

	log.Println("Email from form:", email)
	log.Println("Name from form: ", name)
	log.Println("Password from form:", password)

	_, userError := config.Db.GetUserByEmail(request.Context(), email)

	if userError == nil {
		config.Renderer.Render(writer, "signup", CreateErrorMessageInterfaces("User already exists"))
		return
	}

	hashedPassword, errPwd := HashPassword(password)
	if errPwd != nil {
		respondWithJSON(writer, http.StatusInternalServerError, errPwd.Error())
		return
	}
	createUser := database.CreateUserParams{
		ID:           uuid.New().String(),
		Email:        email,
		Name:         name,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	userAdd, errAdd := config.Db.CreateUser(request.Context(), createUser)
	if errAdd != nil {
		respondWithJSON(writer, http.StatusInternalServerError, errAdd.Error())
		return
	}
	log.Printf("User added with success - UserID:%s \n-Name:%s\n-Email: %s", userAdd.ID, userAdd.Name, userAdd.Email)
	config.Index(writer, request)

}

func UpdateUser[T interface{}](writer http.ResponseWriter, request *http.Request, dbFunc func(context.Context, T) error, config *Config) {

	userID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return
	}
	log.Println("User ID from token on Update User Call:", userID)

	var updateParams T
	err := json.NewDecoder(request.Body).Decode(&updateParams)
	if err != nil {
		respondWithJSON(writer, http.StatusBadRequest, err.Error())
		return
	}

	switch params := any(updateParams).(type) {
	case database.UpdateUserPasswordParams:
		log.Println("Updating user password")
		hashedPassword, errPWD := HashPassword(params.PasswordHash)
		if errPWD != nil {
			respondWithJSON(writer, http.StatusInternalServerError, errPWD.Error())
			return
		}
		params.PasswordHash = hashedPassword
		params.ID = userID
		updateParams = any(params).(T)
	case database.UpdateUserEmailParams:
		log.Println("Updating user email")
		_, userError := config.Db.GetUserByEmail(request.Context(), params.Email)

		// Make that onto the My Account Endpoint
		if userError == nil {
			config.Renderer.Render(writer, "TOBEDEFINED", CreateErrorMessageInterfaces("User already exists"))
			return
		}
		params.ID = userID
		updateParams = any(params).(T)
	case database.UpdateUserNameParams:
		log.Println("Updating user name")
		params.ID = userID
		updateParams = any(params).(T)
	default:
		log.Println("Wrong parameters in request")
		respondWithJSON(writer, http.StatusMethodNotAllowed, "Wrong Parameters in Request")
	}
	errUpdate := dbFunc(request.Context(), updateParams)

	if errUpdate != nil {
		respondWithJSON(writer, http.StatusInternalServerError, errUpdate.Error())
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
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return
	}

	log.Println("User ID from token:", userID)
	errADmin := config.Db.UpdateUserAdmin(request.Context(), userID)
	if errADmin != nil {
		respondWithJSON(writer, http.StatusInternalServerError, errADmin.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, "User successfully upgraded to Admin")
}
