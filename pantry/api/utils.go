package api

import (
	"database/sql"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (tmplt *Templates) Render(writer io.Writer, name string, data interface{}) {
	errorTemplate := tmplt.templates.ExecuteTemplate(writer, name, data)
	if errorTemplate != nil {
		log.Println("Error rendering template:", errorTemplate)
	}
}

func MyTemplates() *Templates {
	templates, _ := template.ParseGlob("static/*.html")
	return &Templates{
		templates: templates,
	}
}

func CloseDB(dbConn *sql.DB) {
	if err := dbConn.Close(); err != nil {
		log.Println("Error closing database connection:", err)
	} else {
		log.Println("Database connection closed successfully.")
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "pantry-pal",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	// subject is the userID setup in the JWT
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}
	return subject, nil
}

func GetJWTFromCookie(request *http.Request) (string, error) {
	jwtToken, errorJwt := request.Cookie("JWTToken")
	if errorJwt != nil {
		return "", errorJwt
	}
	return jwtToken.Value, nil
}

func ValidateDate(givenDate string) bool {
	dateLayout := "2006-01-02"
	formattedDate, errParse := time.Parse(dateLayout, givenDate)

	if errParse != nil {
		return false
	}
	return !formattedDate.Before(time.Now())
}

func GetUserIDFromTokenAndValidate(
	request *http.Request,
	config *Config,
) (string, error) {
	token, errTk := GetJWTFromCookie(request)
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
