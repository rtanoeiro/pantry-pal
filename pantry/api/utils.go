package api

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (tmplt *Templates) Render(writer io.Writer, name string, data interface{}) error {
	return tmplt.templates.ExecuteTemplate(writer, name, data)
}

func MyTemplates() *Templates {
	templates, _ := template.ParseGlob("static/*.html")
	return &Templates{
		templates: templates,
	}
}

func respondWithJSON(writer http.ResponseWriter, code int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(code)
	results, _ := json.Marshal(data)
	writer.Write(results)
}

func CreateErrorMessageInterfaces(message string) map[string]interface{} {
	return map[string]interface{}{
		"Message": message,
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID string, tokenSecret string, expiresIn time.Duration) (string, error) {
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

	log.Println("JWTToken Name:", jwtToken.Name)
	log.Println("JWTToken Value:", jwtToken.Value)

	return jwtToken.Value, nil
}

func checkDate(givenDate string) bool {

	dateLayout := "2006-01-02"
	formattedDate, errParse := time.Parse(dateLayout, givenDate)

	if errParse != nil {
		return false
	}
	results := formattedDate.Compare(time.Now())
	return results != -1
}
