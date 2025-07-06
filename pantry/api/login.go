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

	log.Println("Email from form:", email)
	log.Println("Password from form:", password)

	user, errEmail := config.Db.GetUserByEmail(request.Context(), email)

	if errEmail != nil {
		config.Renderer.Render(writer, "errorLogin", CreateErrorMessageInterfaces("Invalid Email"))
		log.Println("Invalid email during login:", errEmail)
		return
	}
	if CheckPasswordHash(password, user.PasswordHash) != nil {
		config.Renderer.Render(writer, "errorLogin", CreateErrorMessageInterfaces("Wrong Password"))
		log.Println("Invalid password during login")
		return
	}
	log.Println(
		"User details after login. \n- User:",
		user.ID,
		"\n- Email:",
		user.Email,
		"\n- Hashed Password:",
		user.PasswordHash,
		"\n- Created At:",
		user.CreatedAt,
		"\n- Updated At:",
		user.UpdatedAt,
	)

	userJWTToken, errJWTToken := MakeJWT(user.ID, config.Secret, time.Second*3600*6)
	if errJWTToken != nil {
		config.Renderer.Render(
			writer,
			"errorLogin",
			CreateErrorMessageInterfaces("Error request on getting user, please try again"),
		)
		log.Println("Error on making JWT During Login:", errJWTToken)
		return
	}
	log.Println("JWT Token Created with Success during login:", userJWTToken)

	http.SetCookie(writer, &http.Cookie{
		Name:     "JWTToken",
		Value:    userJWTToken,
		Expires:  time.Now().Add(6 * time.Hour),
		HttpOnly: true,
	})
	writer.Header().Set("HX-Redirect", "/home")
	log.Println("User logged in with success. Redirecting to Home Page...")
}

func (config *Config) Logout(writer http.ResponseWriter, request *http.Request) {
	http.SetCookie(writer, &http.Cookie{
		Name:     "JWTToken",
		Value:    "",
		Expires:  time.Now().Add(6 * time.Hour),
		HttpOnly: true,
	})
	writer.Header().Set("HX-Redirect", "/")
	log.Println("User logged out")
}

func (config *Config) SignUp(writer http.ResponseWriter, request *http.Request) {
	config.Renderer.Render(writer, "signup", nil)
	log.Println("User entered SignUp Page...")
}

// TODO: Add Check UserID from JWT
func (config *Config) Home(writer http.ResponseWriter, request *http.Request) {
	returnPantry := map[string]interface{}{
		"ErrorMessage":   "",
		"SuccessMessage": "",
	}

	userID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		returnPantry["ErrorMessage"] = "Unable to retrieve user Pantry Items"
		config.Renderer.Render(writer, "ResponseMessage", returnPantry)
		return
	}

	user, userError := config.Db.GetUserById(request.Context(), userID)
	if userError != nil {
		respondWithJSON(writer, http.StatusUnauthorized, userError.Error())
		returnPantry["ErrorMessage"] = "Unable to retrieve user Pantry Items"
		config.Renderer.Render(writer, "ResponseMessage", returnPantry)
		writer.Header().Set("HX-Redirect", "/")
		return
	}

	returnPantry["UserName"] = user.Name
	config.Renderer.Render(writer, "home", returnPantry)
	log.Println("User entered Home Page...")
}
