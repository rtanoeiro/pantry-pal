package api

import (
	"log"
	"net/http"
	"time"
)

func (config *Config) Login(writer http.ResponseWriter, request *http.Request) {

	email := request.FormValue("email")
	password := request.FormValue("password")

	log.Println("Email from form:", email)
	log.Println("Password from form:", password)

	user, errEmail := config.Db.GetUserByEmail(request.Context(), email)

	if errEmail != nil {
		respondWithError(writer, http.StatusInternalServerError, errEmail.Error())
		return
	}
	if CheckPasswordHash(password, user.PasswordHash) != nil {
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

	http.SetCookie(writer, &http.Cookie{
		Name:     "JWTToken",
		Value:    userJWTToken,
		Expires:  time.Now().Add(6 * time.Hour),
		HttpOnly: true,
	})

	http.Redirect(writer, request, "/home", http.StatusFound) // Redirect to /home
}
