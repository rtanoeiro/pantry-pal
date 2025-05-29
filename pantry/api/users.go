package api

import (
	"log"
	"net/http"
	"pantry-pal/pantry/database"
	"time"

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
	userData, errUser := config.Db.GetUserById(request.Context(), userID)

	if errUser != nil {
		respondWithJSON(writer, http.StatusBadRequest, errUser.Error())
		return
	}

	var usersSlice []User
	returnUser := map[string]interface{}{
		"UserID":    userData.ID,
		"UserName":  userData.Name,
		"UserEmail": userData.Email,
		"IsAdmin":   userData.IsAdmin.Valid,
		"Users":     []User{},
	}

	if userData.IsAdmin.Int64 == 1 {
		allOtherUsers, _ := config.Db.GetAllUsers(request.Context(), userData.ID)
		for _, user := range allOtherUsers {
			usersSlice = append(usersSlice, User{
				UserID:    user.ID,
				Name:      user.Name,
				Email:     user.Email,
				UserAdmin: user.IsAdmin.Int64,
			})
		}
		returnUser["Users"] = usersSlice
	}
	log.Println("User data to be rendered:", returnUser)

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

func (config *Config) DeleteUser(writer http.ResponseWriter, request *http.Request) {
	log.Println("Delete User endpoint called")
	AdminUserID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return
	}
	var usersSlice []User
	returnUser := map[string]interface{}{
		"ErrorMessage":   "",
		"SuccessMessage": "",
		"IsAdmin":        int64(1),
		"Users":          []User{},
	}

	userID := request.PathValue("userID")
	errDelete := config.Db.DeleteUser(request.Context(), userID)
	if errDelete != nil {
		respondWithJSON(writer, http.StatusBadRequest, errDelete.Error())
		returnUser["ErrorMessage"] = "Error on deleting user"
		config.Renderer.Render(writer, "Admin", returnUser)
		return
	}
	allOtherUsers, _ := config.Db.GetAllUsers(request.Context(), AdminUserID)

	for _, user := range allOtherUsers {
		usersSlice = append(usersSlice, User{
			UserID:    user.ID,
			Name:      user.Name,
			Email:     user.Email,
			UserAdmin: user.IsAdmin.Int64,
		},
		)
		returnUser["Users"] = usersSlice
	}

	returnUser["SuccessMessage"] = "User Deleted With Success!"
	config.Renderer.Render(writer, "Admin", returnUser)
}

// TODO: Find a way to improve the replacements of data when rendeding HTML, currently rendering everything
func UpdateUser(writer http.ResponseWriter, request *http.Request, updateType, updateData string, config *Config) {

	userID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return
	}
	log.Println("User ID from token on Update User Call:", userID)
	userData, errUser := config.Db.GetUserById(request.Context(), userID)
	if errUser != nil {
		respondWithJSON(writer, http.StatusBadRequest, errUser.Error())
		return
	}
	returnUser := map[string]interface{}{
		"UserName":       userData.Name,
		"UserEmail":      userData.Email,
		"ErrorMessage":   "",
		"SuccessMessage": "",
	}

	switch updateType {
	case "password":
		log.Println("Updating user password")
		hashedPassword, errPWD := HashPassword(updateData)
		if errPWD != nil {
			respondWithJSON(writer, http.StatusInternalServerError, errPWD.Error())
			returnUser["ErrorMessage"] = "Error on hashing password"
			config.Renderer.Render(writer, "user", returnUser)
		}
		data := database.UpdateUserPasswordParams{
			PasswordHash: hashedPassword,
			ID:           userID,
		}
		errUpdate := config.Db.UpdateUserPassword(request.Context(), data)

		if errUpdate != nil {
			returnUser["ErrorMessage"] = "Error on updating password"
			config.Renderer.Render(writer, "user", returnUser)
			return
		}
		returnUser["SuccessMessage"] = "Password updated with success!"
		config.Renderer.Render(writer, "user", returnUser)
	case "email":
		log.Println("Updating user email")
		data := database.UpdateUserEmailParams{
			Email: updateData,
			ID:    userID,
		}
		errUpdate := config.Db.UpdateUserEmail(request.Context(), data)
		if errUpdate != nil {
			returnUser["ErrorMessage"] = "Error on updating user email"
			config.Renderer.Render(writer, "user", returnUser)
			return
		}
		returnUser["SuccessMessage"] = "Email updated with success!"
		returnUser["UserEmail"] = updateData
		config.Renderer.Render(writer, "user", returnUser)
	case "name":
		log.Println("Updating user name")
		data := database.UpdateUserNameParams{
			Name: updateData,
			ID:   userID,
		}
		errUpdate := config.Db.UpdateUserName(request.Context(), data)
		if errUpdate != nil {
			returnUser["ErrorMessage"] = "Error on updating user Name"
			config.Renderer.Render(writer, "user", returnUser)
			return
		}
		returnUser["SuccessMessage"] = "Name updated with success!"
		returnUser["UserName"] = updateData
		config.Renderer.Render(writer, "user", returnUser)
	default:
		log.Println("Wrong parameters in request")
		writer.Header().Set("HX-Redirect", "/user")
	}
}

func (config *Config) UpdateUserEmail(writer http.ResponseWriter, request *http.Request) {
	log.Println("Update User Email endpoint called")
	email := request.FormValue("email")
	UpdateUser(writer, request, "email", email, config)
}

func (config *Config) UpdateUserName(writer http.ResponseWriter, request *http.Request) {
	log.Println("Update User Name endpoint called")
	name := request.FormValue("name")
	UpdateUser(writer, request, "name", name, config)
}

func (config *Config) UpdateUserPassword(writer http.ResponseWriter, request *http.Request) {
	log.Println("Update User Password endpoint called")
	password := request.FormValue("password")
	UpdateUser(writer, request, "password", password, config)
}

// TODO: Modify these admin functions to be more generic
func (config *Config) AddUserAdmin(writer http.ResponseWriter, request *http.Request) {
	toUpdateuserID := request.PathValue("userID")
	log.Println("User Add Admin endpoint called")

	AdminUserID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return
	}

	var usersSlice []User
	returnUser := map[string]interface{}{
		"ErrorMessage":   "",
		"SuccessMessage": "",
		"IsAdmin":        int64(1),
		"Users":          []User{},
	}
	errADmin := config.Db.MakeUserAdmin(request.Context(), toUpdateuserID)
	if errADmin != nil {
		returnUser["ErrorMessage"] = "Unable to Assign Admin to User!"
		config.Renderer.Render(writer, "ResponseMessage", returnUser)
		return
	}
	allOtherUsers, _ := config.Db.GetAllUsers(request.Context(), AdminUserID)
	for _, user := range allOtherUsers {
		usersSlice = append(usersSlice, User{
			UserID:    user.ID,
			Name:      user.Name,
			Email:     user.Email,
			UserAdmin: user.IsAdmin.Int64,
		},
		)
		returnUser["Users"] = usersSlice
	}
	returnUser["SuccessMessage"] = "User Added as Admin with Success!"
	config.Renderer.Render(writer, "Admin", returnUser)
}

func (config *Config) RemoveUserAdmin(writer http.ResponseWriter, request *http.Request) {
	toUpdateuserID := request.PathValue("userID")
	log.Println("User Remove Admin endpoint called")

	AdminUserID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return
	}
	var usersSlice []User
	returnUser := map[string]interface{}{
		"ErrorMessage":   "",
		"SuccessMessage": "",
		"IsAdmin":        int64(1),
		"Users":          []User{},
	}

	errADmin := config.Db.RemoveUserAdmin(request.Context(), toUpdateuserID)
	if errADmin != nil {
		returnUser["ErrorMessage"] = "Unable to Remove Admin to ser!"
		config.Renderer.Render(writer, "ResponseMessage", returnUser)
		return
	}

	allOtherUsers, _ := config.Db.GetAllUsers(request.Context(), AdminUserID)
	for _, user := range allOtherUsers {
		usersSlice = append(usersSlice, User{
			UserID:    user.ID,
			Name:      user.Name,
			Email:     user.Email,
			UserAdmin: user.IsAdmin.Int64,
		},
		)
		returnUser["Users"] = usersSlice
	}

	returnUser["SuccessMessage"] = "User Removed as Admin with Success!"
	config.Renderer.Render(writer, "Admin", returnUser)
}
