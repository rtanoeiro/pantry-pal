package api

import (
	"log"
	"net/http"
	"time"
)

func (config *Config) Index(writer http.ResponseWriter, request *http.Request) {
	log.Println("User entered login page")
	writer.Header().Set("HX-Replace-Url", "/")
	config.Renderer.Render(writer, "index", nil)
}

func (config *Config) Login(writer http.ResponseWriter, request *http.Request) {
	email := request.FormValue("email")
	password := request.FormValue("password")

	user, errEmail := config.Db.GetUserByEmail(request.Context(), email)

	if errEmail != nil {
		config.Renderer.Render(writer, "errorLogin", CreateErrorMessageInterfaces("Invalid Email"))
		log.Printf("Email %s failed login at %s:", email, time.Now())
		return
	}
	if CheckPasswordHash(password, user.PasswordHash) != nil {
		config.Renderer.Render(writer, "errorLogin", CreateErrorMessageInterfaces("Wrong Password"))
		log.Printf("Email %s failed login with wrong password at %s:", email, time.Now())
		return
	}

	userJWTToken, errJWTToken := MakeJWT(user.ID, config.Secret, time.Second*3600*6)
	if errJWTToken != nil {
		config.Renderer.Render(
			writer,
			"errorLogin",
			CreateErrorMessageInterfaces("Error request on getting user, please try again"),
		)
		log.Printf("Failed creating JWT Tokenat %s:", time.Now())
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:     "JWTToken",
		Value:    userJWTToken,
		Expires:  time.Now().Add(6 * time.Hour),
		HttpOnly: true,
	})
	writer.Header().Set("HX-Redirect", "/home")
	log.Printf("User %s logged in with success. Redirecting to Home Page...", email)
}

func (config *Config) Logout(writer http.ResponseWriter, request *http.Request) {
	http.SetCookie(writer, &http.Cookie{
		Name:     "JWTToken",
		Value:    "",
		Expires:  time.Now().Add(6 * time.Hour),
		HttpOnly: true,
	})
	writer.Header().Set("HX-Redirect", "/login")
	userID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		config.Renderer.Render(writer, "ResponseMessage", CreateErrorMessageInterfaces("Unable to retrieve user data"))
		return
	}

	log.Printf("User %s logged out at %s", userID, time.Now())
}

func (config *Config) SignUp(writer http.ResponseWriter, request *http.Request) {
	config.Renderer.Render(writer, "signup", nil)
}

func (config *Config) Home(writer http.ResponseWriter, request *http.Request) {
	var returnPantry CurrentUserRequest
	userID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		returnPantry.ErrorMessage = "Unable to retrieve user Pantry Items"
		config.Renderer.Render(writer, "ResponseMessage", returnPantry)
		return
	}

	user, userError := config.Db.GetUserById(request.Context(), userID)
	if userError != nil {
		respondWithJSON(writer, http.StatusUnauthorized, userError.Error())
		returnPantry.ErrorMessage = "Unable to retrieve user Pantry Items"
		config.Renderer.Render(writer, "ResponseMessage", returnPantry)
		writer.Header().Set("HX-Redirect", "/")
		return
	}

	returnPantry.UserName = user.Name
	config.Renderer.Render(writer, "home", returnPantry)
	log.Println("User entered Home Page...")
}
