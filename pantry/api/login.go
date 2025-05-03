package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (config *Config) Login(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)
	loginRequest := UserLoginRequest{}
	err := decoder.Decode(&loginRequest)
	if err != nil {
		respondWithError(writer, http.StatusBadRequest, err.Error())
		return
	}

	user, errUser := config.Db.GetUserByEmail(request.Context(), loginRequest.Email)

	if errUser != nil {
		respondWithError(writer, http.StatusInternalServerError, errUser.Error())
		return
		return
	}
	if CheckPasswordHash(loginRequest.Password, user.PasswordHash) != nil {
		respondWithError(writer, http.StatusUnauthorized, "Wrong Password, try again")
		return
	}
	log.Println("User details after login. \n- User:", user.ID, "\n- Email:", user.Email, "\n- Hashed Password:", user.PasswordHash, "\n- Created At:", user.CreatedAt, "\n- Updated At:", user.UpdatedAt)

	userJWTToken, errJWTToken := MakeJWT(user.ID, config.Secret, time.Second*3600*6)
	if errJWTToken != nil {
		respondWithError(writer, http.StatusUnauthorized, errJWTToken.Error())
		return
	}
	log.Println("JWT Token Created with Success during login:", userJWTToken)

	loginResponse := UserLoginResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		JWTToken:  &userJWTToken,
	}
	data, err := json.Marshal(loginResponse)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, data)
}
