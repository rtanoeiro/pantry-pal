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

	userInfo := UserInfoRequest{
		ID:             userData.ID,
		UserName:       userData.Name,
		UserEmail:      userData.Email,
		IsAdmin:        userData.IsAdmin.Valid,
		Users:          []User{},
		ErrorMessage:   "",
		SuccessMessage: "",
	}
	userInfo.Users = config.GetAllOtherUsers(writer, request, userInfo)
	config.Renderer.Render(writer, "user", userInfo)
}

func (config *Config) GetAllOtherUsers(writer http.ResponseWriter, request *http.Request, userInfo UserInfoRequest) []User {

	var usersSlice []User
	if userInfo.IsAdmin {
		allOtherUsers, _ := config.Db.GetAllUsers(request.Context(), userInfo.ID)
		for _, user := range allOtherUsers {
			usersSlice = append(usersSlice, User{
				UserID:    user.ID,
				Name:      user.Name,
				Email:     user.Email,
				UserAdmin: user.IsAdmin.Int64,
			})
		}
	}
	log.Println("All other available users:", usersSlice)
	return usersSlice
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
	adminUserID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return
	}
	userData, errUser := config.Db.GetUserById(request.Context(), adminUserID)
	if errUser != nil {
		respondWithJSON(writer, http.StatusBadRequest, errUser.Error())
		return
	}

	userInfo := UserInfoRequest{
		ID:             userData.ID,
		UserName:       userData.Name,
		UserEmail:      userData.Email,
		IsAdmin:        userData.IsAdmin.Valid,
		Users:          []User{},
		ErrorMessage:   "",
		SuccessMessage: "",
	}
	userInfo.Users = config.GetAllOtherUsers(writer, request, userInfo)
	userIDToDelete := request.PathValue("userID")
	errDelete := config.Db.DeleteUser(request.Context(), userIDToDelete)
	if errDelete != nil {
		respondWithJSON(writer, http.StatusBadRequest, errDelete.Error())
		userInfo.ErrorMessage = "Error on deleting user"
		config.Renderer.Render(writer, "Admin", userInfo)
		return
	}
	config.Renderer.Render(writer, "Admin", userInfo)

}

// TODO: Find a way to improve the replacements of data when rendeding HTML, currently rendering everything
func (config *Config) UpdateUser(writer http.ResponseWriter, request *http.Request, updateType, updateData string) {

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
	userInfo := UserInfoRequest{
		ID:             userData.ID,
		UserName:       userData.Name,
		UserEmail:      userData.Email,
		IsAdmin:        userData.IsAdmin.Valid,
		Users:          []User{},
		ErrorMessage:   "",
		SuccessMessage: "",
	}

	switch updateType {
	case "password":
		log.Printf("Updating password for userID: %s", userID)
		config.handlePassword(writer, request, &userInfo, updateData)
		return
	case "email":
		log.Printf("Updating email for userID: %s", userID)
		config.handleEmail(writer, request, &userInfo, updateData)
		return
	case "name":
		log.Printf("Updating name for userID: %s", userID)
		config.handleName(writer, request, &userInfo, updateData)
		return
	default:
		log.Println("Wrong parameters in request")
		writer.Header().Set("HX-Redirect", "/user")
		return
	}
}

func (config *Config) handleName(writer http.ResponseWriter, request *http.Request, userInfo *UserInfoRequest, updateData string) {
	data := database.UpdateUserNameParams{
		Name: updateData,
		ID:   userInfo.ID,
	}
	errUpdate := config.Db.UpdateUserName(request.Context(), data)
	if errUpdate != nil {
		userInfo.ErrorMessage = "Error on updating user Name"
		config.Renderer.Render(writer, "user", userInfo)
	}
	userInfo.Users = config.GetAllOtherUsers(writer, request, *userInfo)

	userInfo.SuccessMessage = "Name updated with success!"
	userInfo.UserName = updateData
	config.Renderer.Render(writer, "user", userInfo)
}

func (config *Config) handleEmail(writer http.ResponseWriter, request *http.Request, userInfo *UserInfoRequest, updateData string) {
	data := database.UpdateUserEmailParams{
		Email: updateData,
		ID:    userInfo.ID,
	}
	userInfo.Users = config.GetAllOtherUsers(writer, request, *userInfo)

	errUpdate := config.Db.UpdateUserEmail(request.Context(), data)
	if errUpdate != nil {
		userInfo.ErrorMessage = "Error on updating user email"
		config.Renderer.Render(writer, "user", userInfo)
		return
	}
	userInfo.SuccessMessage = "Email updated with success!"
	userInfo.UserEmail = updateData
	config.Renderer.Render(writer, "user", userInfo)
}

func (config *Config) handlePassword(writer http.ResponseWriter, request *http.Request, userInfo *UserInfoRequest, updateData string) {
	hashedPassword, errPWD := HashPassword(updateData)
	if errPWD != nil {
		respondWithJSON(writer, http.StatusInternalServerError, errPWD.Error())
		userInfo.ErrorMessage = "Error on hashing password"
		config.Renderer.Render(writer, "user", userInfo)
	}
	data := database.UpdateUserPasswordParams{
		PasswordHash: hashedPassword,
		ID:           userInfo.ID,
	}
	userInfo.Users = config.GetAllOtherUsers(writer, request, *userInfo)

	errUpdate := config.Db.UpdateUserPassword(request.Context(), data)
	if errUpdate != nil {
		userInfo.ErrorMessage = "Error on updating password"
		config.Renderer.Render(writer, "user", userInfo)
		return
	}

	userInfo.SuccessMessage = "Password updated with success!"
	config.Renderer.Render(writer, "user", userInfo)
}

func (config *Config) UpdateUserEmail(writer http.ResponseWriter, request *http.Request) {
	log.Println("Update User Email endpoint called")
	email := request.FormValue("email")
	config.UpdateUser(writer, request, "email", email)
}

func (config *Config) UpdateUserName(writer http.ResponseWriter, request *http.Request) {
	log.Println("Update User Name endpoint called")
	name := request.FormValue("name")
	config.UpdateUser(writer, request, "name", name)
}

func (config *Config) UpdateUserPassword(writer http.ResponseWriter, request *http.Request) {
	log.Println("Update User Password endpoint called")
	password := request.FormValue("password")
	config.UpdateUser(writer, request, "password", password)
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

	userData, errUser := config.Db.GetUserById(request.Context(), AdminUserID)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return
	}
	returnUser := UserInfoRequest{
		ID:             userData.ID,
		UserName:       userData.Name,
		UserEmail:      userData.Email,
		IsAdmin:        true,
		ErrorMessage:   "",
		SuccessMessage: "",
		Users:          []User{},
	}

	errADmin := config.Db.MakeUserAdmin(request.Context(), toUpdateuserID)
	if errADmin != nil {
		returnUser.ErrorMessage = "Unable to Assign Admin to User!"
		config.Renderer.Render(writer, "ResponseMessage", returnUser)
		return
	}

	returnUser.Users = config.GetAllOtherUsers(writer, request, returnUser)
	returnUser.SuccessMessage = "User Added as Admin with Success!"
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

	userData, errUser := config.Db.GetUserById(request.Context(), AdminUserID)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return
	}
	returnUser := UserInfoRequest{
		ID:             userData.ID,
		UserName:       userData.Name,
		UserEmail:      userData.Email,
		IsAdmin:        true,
		ErrorMessage:   "",
		SuccessMessage: "",
		Users:          []User{},
	}
	errADmin := config.Db.RemoveUserAdmin(request.Context(), toUpdateuserID)
	if errADmin != nil {
		returnUser.ErrorMessage = "Unable to Remove Admin to user!"
		config.Renderer.Render(writer, "ResponseMessage", returnUser)
		return
	}

	returnUser.Users = config.GetAllOtherUsers(writer, request, returnUser)
	returnUser.SuccessMessage = "User Removed as Admin with Success!"
	config.Renderer.Render(writer, "Admin", returnUser)
}
