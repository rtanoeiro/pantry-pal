package api

import (
	"log"
	"net/http"
	"time"

	"pantry-pal/pantry/database"

	"github.com/google/uuid"
)

func (config *Config) CreateUser(writer http.ResponseWriter, request *http.Request) {
	email := request.FormValue("email")
	name := request.FormValue("name")
	password := request.FormValue("password")

	_, userError := config.Db.GetUserByEmail(request.Context(), email)

	if userError == nil {
		config.Renderer.Render(
			writer,
			"signup",
			CreateErrorMessageInterfaces("User already exists"),
		)
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
	log.Printf(
		"User added with success at %s- UserID:%s \n-Name:%s\n-Email: %s",
		userAdd.ID,
		userAdd.Name,
		userAdd.Email,
		time.Now(),
	)
	config.Index(writer, request)
}

func (config *Config) GetUserInfo(writer http.ResponseWriter, request *http.Request) {
	userID, errUser := GetUserIDFromToken(request, writer, config)
	var UserPageData UserInfoRequest
	// TODO: Create fuction to redirect to a 401 page
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return
	}
	userData, errUser := config.Db.GetUserById(request.Context(), userID)
	if errUser != nil {
		UserPageData.ErrorMessage = "Unable to load user info"
		config.Renderer.Render(writer, "user", UserPageData)
		return
	}
	UserPageData.ID = userData.ID
	UserPageData.UserName = userData.Name
	UserPageData.UserEmail = userData.Email
	UserPageData.IsAdmin = userData.IsAdmin.Valid

	config.GetAllOtherUsers(writer, request, &UserPageData)
	config.Renderer.Render(writer, "user", UserPageData)
}

func (config *Config) GetAllOtherUsers(
	writer http.ResponseWriter,
	request *http.Request,
	userInfo *UserInfoRequest,
) []User {
	var usersSlice []User
	if userInfo.IsAdmin {
		allOtherUsers, _ := config.Db.GetAllUsers(request.Context(), userInfo.ID)
		for _, user := range allOtherUsers {
			usersSlice = append(usersSlice, User{
				ID:        user.ID,
				UserName:  user.Name,
				UserEmail: user.Email,
				IsAdmin:   user.IsAdmin.Valid,
			})
		}
	}
	userInfo.Users = usersSlice
	return usersSlice
}

func (config *Config) DeleteUser(writer http.ResponseWriter, request *http.Request) {
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
	userIDToDelete := request.PathValue("userID")
	errDelete := config.Db.DeleteUser(request.Context(), userIDToDelete)
	if errDelete != nil {
		respondWithJSON(writer, http.StatusBadRequest, errDelete.Error())
		userInfo.ErrorMessage = "Error on deleting user"
		config.Renderer.Render(writer, "Admin", userInfo)
		return
	}
	config.GetAllOtherUsers(writer, request, &userInfo)
	config.Renderer.Render(writer, "Admin", userInfo)
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

func (config *Config) UpdateUser(
	writer http.ResponseWriter,
	request *http.Request,
	updateType, updateData string,
) {
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
	// TODO: No need for the whole list of other users, just the updated user data
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
		log.Printf("Updating password for userID: %s at %s", userID, time.Now())
		config.handlePassword(writer, request, &userInfo, updateData)
		return
	case "email":
		log.Printf("Updating email for userID: %s at %s", userID, time.Now())
		config.handleEmail(writer, request, &userInfo, updateData)
		return
	case "name":
		log.Printf("Updating name for userID:%s at %s", userID, time.Now())
		config.handleName(writer, request, &userInfo, updateData)
		return
	default:
		log.Println("Wrong parameters in request")
		writer.Header().Set("HX-Redirect", "/user")
		return
	}
}

func (config *Config) handleName(
	writer http.ResponseWriter,
	request *http.Request,
	userInfo *UserInfoRequest,
	updateData string,
) {
	data := database.UpdateUserNameParams{
		Name: updateData,
		ID:   userInfo.ID,
	}
	errUpdate := config.Db.UpdateUserName(request.Context(), data)
	if errUpdate != nil {
		userInfo.ErrorMessage = "Error on updating user Name"
		config.Renderer.Render(writer, "user", userInfo)
	}
	userInfo.SuccessMessage = "Name updated with success!"
	userInfo.UserName = updateData
	config.Renderer.Render(writer, "UserInformation", userInfo)
}

func (config *Config) handleEmail(
	writer http.ResponseWriter,
	request *http.Request,
	userInfo *UserInfoRequest,
	updateData string,
) {
	data := database.UpdateUserEmailParams{
		Email: updateData,
		ID:    userInfo.ID,
	}
	errUpdate := config.Db.UpdateUserEmail(request.Context(), data)
	if errUpdate != nil {
		userInfo.ErrorMessage = "Error on updating user email"
		config.Renderer.Render(writer, "user", userInfo)
		return
	}
	userInfo.SuccessMessage = "Email updated with success!"
	userInfo.UserEmail = updateData
	config.Renderer.Render(writer, "UserInformation", userInfo)
}

func (config *Config) handlePassword(
	writer http.ResponseWriter,
	request *http.Request,
	userInfo *UserInfoRequest,
	updateData string,
) {
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

	errUpdate := config.Db.UpdateUserPassword(request.Context(), data)
	if errUpdate != nil {
		userInfo.ErrorMessage = "Error on updating password"
		config.Renderer.Render(writer, "user", userInfo)
		return
	}

	userInfo.SuccessMessage = "Password updated with success!"
	config.Renderer.Render(writer, "user", userInfo)
}

func (config *Config) AddUserAdmin(writer http.ResponseWriter, request *http.Request) {
	returnUser, userToAddAdmin := config.prepareUserAdmin(request, writer)

	errADmin := config.Db.MakeUserAdmin(request.Context(), userToAddAdmin)
	if errADmin != nil {
		returnUser.ErrorMessage = "Unable to Assign Admin to User!"
		config.Renderer.Render(writer, "ResponseMessage", returnUser)
		return
	}

	config.GetAllOtherUsers(writer, request, &returnUser)
	returnUser.SuccessMessage = "User Added as Admin with Success!"
	config.Renderer.Render(writer, "Admin", returnUser)
}

func (config *Config) RemoveUserAdmin(writer http.ResponseWriter, request *http.Request) {
	returnUser, userToRemoveAdmin := config.prepareUserAdmin(request, writer)

	errADmin := config.Db.RemoveUserAdmin(request.Context(), userToRemoveAdmin)
	if errADmin != nil {
		returnUser.ErrorMessage = "Unable to Remove Admin to user!"
		config.Renderer.Render(writer, "ResponseMessage", returnUser)
		return
	}

	config.GetAllOtherUsers(writer, request, &returnUser)
	returnUser.SuccessMessage = "User Removed as Admin with Success!"
	config.Renderer.Render(writer, "Admin", returnUser)
}

func (config *Config) prepareUserAdmin(request *http.Request, writer http.ResponseWriter) (UserInfoRequest, string) {
	UserID := request.PathValue("userID")
	log.Printf("User Add Admin endpoint called. Path Value User ID %s", UserID)

	AdminUserID, errUser := GetUserIDFromToken(request, writer, config)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return UserInfoRequest{}, ""
	}

	userData, errUser := config.Db.GetUserById(request.Context(), AdminUserID)
	if errUser != nil {
		respondWithJSON(writer, http.StatusUnauthorized, errUser.Error())
		return UserInfoRequest{}, ""
	}
	// TODO: No need for the whole user info request, just a list of users
	returnUser := UserInfoRequest{
		ID:             userData.ID,
		UserName:       userData.Name,
		UserEmail:      userData.Email,
		IsAdmin:        true,
		ErrorMessage:   "",
		SuccessMessage: "",
		Users:          []User{},
	}
	return returnUser, UserID

}
