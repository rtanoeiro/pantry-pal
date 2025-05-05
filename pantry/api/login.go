package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (config *Config) Login(writer http.ResponseWriter, request *http.Request) {

	userRequest := LoginUserRequest{}

	err := json.NewDecoder(request.Body).Decode(&userRequest)
	if err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Println("Email from form:", userRequest.Email)
	log.Println("Password from form:", userRequest.Password)

	user, errUser := config.Db.GetUserByEmail(request.Context(), userRequest.Email)

	if errUser != nil {
		respondWithError(writer, http.StatusInternalServerError, errUser.Error())
		return
	}
	if CheckPasswordHash(userRequest.Password, user.PasswordHash) != nil {
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

	loginResponse := LoginUserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		//JWTToken:  &userJWTToken,
	}

	data, err := json.Marshal(loginResponse)
	if err != nil {
		respondWithError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(writer, http.StatusOK, data)
}
