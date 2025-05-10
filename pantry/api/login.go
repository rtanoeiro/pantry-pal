package api

import (
	"log"
	"net/http"
	"time"
)

func (config *Config) Index(writer http.ResponseWriter, request *http.Request) {
	config.Renderer.Render(writer, "index", nil, request.Context())
}

func (config *Config) Login(writer http.ResponseWriter, request *http.Request) {

	log.Println("User entered Login Page...")
	email := request.FormValue("email")
	password := request.FormValue("password")

	log.Println("Email from form:", email)
	log.Println("Password from form:", password)

	user, errEmail := config.Db.GetUserByEmail(request.Context(), email)

	if errEmail != nil {
		respondWithError(writer, http.StatusOK, "Invalid Email")
		config.Renderer.Render(writer, "errorMessage", "Invalid Email", request.Context())
		return
	}
	if CheckPasswordHash(password, user.PasswordHash) != nil {
		respondWithError(writer, http.StatusOK, "Wrong Password")
		config.Renderer.Render(writer, "errorMessage", "Wrong Password", request.Context())
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

	writer.Header().Add("HX-Redirect", "/home")
}

func (config *Config) Logout(writer http.ResponseWriter, request *http.Request) {

	log.Println("User logged out")
	http.SetCookie(writer, &http.Cookie{
		Name:     "JWTToken",
		Value:    "",
		Expires:  time.Now().Add(6 * time.Hour),
		HttpOnly: true,
	})
	writer.Header().Add("HX-Redirect", "/")
	writer.Header().Add("HX-Replace-Url", "/")
}

func (config *Config) SignUp(writer http.ResponseWriter, request *http.Request) {

	config.Renderer.Render(writer, "signup", "", request.Context())
	writer.Header().Add("HX-Replace-Url", "/signup")
	writer.Header().Add("HX-Redirect", "/signup")

}

func (config *Config) Home(writer http.ResponseWriter, request *http.Request) {

	config.Renderer.Render(writer, "home", "", request.Context())
	writer.Header().Add("HX-Replace-Url", "/home")
	writer.Header().Add("HX-Redirect", "/home")

}
