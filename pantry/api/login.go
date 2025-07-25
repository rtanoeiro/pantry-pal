package api

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (config *Config) Index(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("HX-Replace-Url", "/login")
	writer.WriteHeader(http.StatusOK)
	_ = config.Renderer.Render(writer, "index", nil)
}

func (config *Config) SignUp(writer http.ResponseWriter, request *http.Request) {
	_ = config.Renderer.Render(writer, "signup", nil)
}

func (config *Config) Logout(writer http.ResponseWriter, request *http.Request) {
	var renderLogout SuccessErrorResponse
	http.SetCookie(writer, &http.Cookie{
		Name:     "JWTToken",
		Value:    "",
		Expires:  time.Now().Add(6 * time.Hour),
		HttpOnly: true,
	})
	writer.Header().Set("HX-Redirect", "/login")
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		renderLogout.ErrorMessage = fmt.Sprintf("Unable to retrieve user data. Error: %s", errUser.Error())
		writer.WriteHeader(http.StatusUnauthorized)
		_ = config.Renderer.Render(writer, "ResponseMessage", renderLogout)
		return
	}

	writer.WriteHeader(http.StatusOK)
	log.Printf("User %s logged out at %s", userID, time.Now())
}

func (config *Config) Login(writer http.ResponseWriter, request *http.Request) {
	var returnResponse SuccessErrorResponse
	email := request.FormValue("email")
	password := request.FormValue("password")

	user, errEmail := config.Db.GetUserByEmail(request.Context(), email)

	if errEmail != nil {
		log.Printf("Email %s failed login at %s:", email, time.Now())
		returnResponse.ErrorMessage = "Invalid Email"
		_ = config.Renderer.Render(writer, "errorLogin", returnResponse)
		return
	}
	if CheckPasswordHash(password, user.PasswordHash) != nil {
		log.Printf("Email %s failed login with wrong password at %s:", email, time.Now())
		returnResponse.ErrorMessage = "Wrong Password"
		_ = config.Renderer.Render(writer, "errorLogin", returnResponse)
		return
	}

	userJWTToken, errJWTToken := MakeJWT(user.ID, config.Secret, time.Second*3600*6)
	if errJWTToken != nil {
		log.Printf("Failed creating JWT Tokenat %s:", time.Now())
		returnResponse.ErrorMessage = "Error request on getting user, please try again"
		_ = config.Renderer.Render(writer, "errorLogin", returnResponse)
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

func (config *Config) Home(writer http.ResponseWriter, request *http.Request) {
	var returnPantry CurrentUserRequest
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		returnPantry.ErrorMessage = fmt.Sprintf("Unable to retrieve user Pantry Items. Error: %s", errUser.Error())
		_ = config.Renderer.Render(writer, "ResponseMessage", returnPantry)
		return
	}

	user, userError := config.Db.GetUserById(request.Context(), userID)
	if userError != nil {
		returnPantry.ErrorMessage = fmt.Sprintf("Unable to retrieve user Pantry Items. Error: %s", userError.Error())
		_ = config.Renderer.Render(writer, "ResponseMessage", returnPantry)
		writer.Header().Set("HX-Redirect", "/")
		return
	}

	returnPantry.UserName = user.Name
	_ = config.Renderer.Render(writer, "home", returnPantry)
	log.Println("User entered Home Page...")
}
