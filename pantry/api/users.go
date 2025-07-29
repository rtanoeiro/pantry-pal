package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"pantry-pal/pantry/database"

	"github.com/google/uuid"
)

func (config *Config) CreateUser(writer http.ResponseWriter, request *http.Request) {
	var userInfo UserInfoRequest
	name := request.FormValue("name")
	password := request.FormValue("password")

	if name == "" || password == "" {
		userInfo.ErrorMessage = "Please provide valid date for all fields"
		writer.WriteHeader(http.StatusBadRequest)
		_ = config.Renderer.Render(writer, "signup", userInfo)
		return
	}
	config.validateUniqueName(name, &userInfo, writer)

	hashedPassword, errPwd := HashPassword(password)
	if errPwd != nil {
		userInfo.ErrorMessage = "Server error, please try again"
		writer.WriteHeader(http.StatusInternalServerError)
		_ = config.Renderer.Render(writer, "signup", userInfo)
		return
	}
	createUser := database.CreateUserParams{
		ID:           uuid.New().String(),
		Name:         name,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	userAdd, errAdd := config.Db.CreateUser(request.Context(), createUser)
	if errAdd != nil {
		userInfo.ErrorMessage = "Failed adding user to our systems, please try again"
		writer.WriteHeader(http.StatusInternalServerError)
		_ = config.Renderer.Render(writer, "signup", userInfo)
		return
	}
	log.Printf(
		"User added with success at %s- UserID:%s \n-Name:%s",
		time.Now(),
		userAdd.ID,
		userAdd.Name,
	)
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("HX-Push-Url", "/login")
	_ = config.Renderer.Render(writer, "index", nil)

}

func (config *Config) GetUserInfo(writer http.ResponseWriter, request *http.Request) {
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	var userInfo UserInfoRequest
	// TODO: Create fuction to redirect to a 401 page
	if errUser != nil {
		userInfo.ErrorMessage = fmt.Sprintf("Unable to load user info. Error: %s", errUser.Error())
		writer.WriteHeader(http.StatusForbidden)
		_ = config.Renderer.Render(writer, "user", userInfo)
		return
	}
	userInfo = config.getUserInformation(userID, userInfo, writer)

	config.GetAllOtherUsers(writer, request, &userInfo)
	writer.WriteHeader(http.StatusOK)
	_ = config.Renderer.Render(writer, "user", userInfo)
}

func (config *Config) GetAllOtherUsers(
	writer http.ResponseWriter,
	request *http.Request,
	userInfo *UserInfoRequest,
) []User {
	var usersSlice []User
	if userInfo.IsAdmin == 1 {
		allOtherUsers, _ := config.Db.GetAllUsers(request.Context(), userInfo.ID)
		for _, user := range allOtherUsers {
			usersSlice = append(usersSlice, User{
				UserID:      user.ID,
				UserName:    user.Name,
				IsUserAdmin: user.IsAdmin.Int64,
			})
		}
	}
	userInfo.Users = usersSlice
	log.Println("Users available in application: ", usersSlice)
	return usersSlice
}

func (config *Config) DeleteUser(writer http.ResponseWriter, request *http.Request) {
	var userInfo UserInfoRequest
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		userInfo.ErrorMessage = fmt.Sprintf("Unable to get current user data. Error: %s", errUser.Error())
		writer.WriteHeader(http.StatusForbidden)
		_ = config.Renderer.Render(writer, "ResponseMessage", userInfo)
		return
	}
	userInfo = config.getUserInformation(userID, userInfo, writer)
	// All data fetched above is just to render part of the data on the final screen
	// It's needed because All GetAllOtherUsers function checks if the user making this request is an Admin

	userIDToDelete := request.PathValue("UserID")
	errDelete := config.Db.DeleteUser(request.Context(), userIDToDelete)
	if errDelete != nil {
		userInfo.ErrorMessage = fmt.Sprintf("Unable to delete user. Please try again. Error: %s", errDelete.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		_ = config.Renderer.Render(writer, "Admin", userInfo)
		return
	}
	config.GetAllOtherUsers(writer, request, &userInfo)
	writer.WriteHeader(http.StatusOK)
	_ = config.Renderer.Render(writer, "Admin", userInfo)
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
	var userInfo UserInfoRequest
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		log.Println("Unable to get user ID from Token")
		userInfo.ErrorMessage = fmt.Sprintf("Unable to get current user data. Error: %s", errUser.Error())
		writer.WriteHeader(http.StatusBadRequest)
		_ = config.Renderer.Render(writer, "ResponseMessage", userInfo)
		return
	}
	userInfo = config.getUserInformation(userID, userInfo, writer)

	switch updateType {
	case "password":
		log.Printf("Updating password for userID: %s at %s", userID, time.Now())
		config.handlePassword(writer, request, &userInfo, updateData)
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
		writer.WriteHeader(http.StatusBadRequest)
		_ = config.Renderer.Render(writer, "user", userInfo)
	}
	userInfo.SuccessMessage = "Name updated with success!"
	userInfo.UserName = updateData
	writer.WriteHeader(http.StatusOK)
	_ = config.Renderer.Render(writer, "UserInformation", userInfo)
}

func (config *Config) handlePassword(
	writer http.ResponseWriter,
	request *http.Request,
	userInfo *UserInfoRequest,
	updateData string,
) {
	hashedPassword, errPWD := HashPassword(updateData)
	if errPWD != nil {
		userInfo.ErrorMessage = "Error on changing password, please try again"
		writer.WriteHeader(http.StatusBadRequest)
		_ = config.Renderer.Render(writer, "user", userInfo)
	}
	data := database.UpdateUserPasswordParams{
		PasswordHash: hashedPassword,
		ID:           userInfo.ID,
	}

	errUpdate := config.Db.UpdateUserPassword(request.Context(), data)
	if errUpdate != nil {
		userInfo.ErrorMessage = "Error on updating password"
		writer.WriteHeader(http.StatusInternalServerError)
		_ = config.Renderer.Render(writer, "user", userInfo)
		return
	}

	userInfo.SuccessMessage = "Password updated with success!"
	writer.WriteHeader(http.StatusOK)
	_ = config.Renderer.Render(writer, "UserInformation", userInfo)
}

func (config *Config) AddUserAdmin(writer http.ResponseWriter, request *http.Request) {
	returnUser, userToAddAdmin := config.prepareUserAdmin(request, writer)
	errADmin := config.Db.MakeUserAdmin(request.Context(), userToAddAdmin)
	if errADmin != nil {
		returnUser.ErrorMessage = "Unable to Assign Admin to User!"
		_ = config.Renderer.Render(writer, "ResponseMessage", returnUser)
		return
	}
	config.GetAllOtherUsers(writer, request, &returnUser)
	log.Printf("User %s has been assigned adming privileges at %s from %s", userToAddAdmin, time.Now(), returnUser.UserName)
	returnUser.SuccessMessage = "User Added as Admin with Success!"
	_ = config.Renderer.Render(writer, "Admin", returnUser)
}

func (config *Config) RevokeUserAdmin(writer http.ResponseWriter, request *http.Request) {
	returnUser, userToRemoveAdmin := config.prepareUserAdmin(request, writer)
	errADmin := config.Db.RemoveUserAdmin(request.Context(), userToRemoveAdmin)
	if errADmin != nil {
		returnUser.ErrorMessage = "Unable to Remove Admin to user!"
		_ = config.Renderer.Render(writer, "ResponseMessage", returnUser)
		return
	}
	config.GetAllOtherUsers(writer, request, &returnUser)
	log.Printf("User %s has lost adming privileges at %s from %s", userToRemoveAdmin, time.Now(), returnUser.UserName)
	returnUser.SuccessMessage = "User Removed as Admin with Success!"
	_ = config.Renderer.Render(writer, "Admin", returnUser)
}

func (config *Config) prepareUserAdmin(request *http.Request, writer http.ResponseWriter) (UserInfoRequest, string) {
	UserID := request.PathValue("UserID")
	var userInfo UserInfoRequest
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)

	if errUser != nil {
		userInfo.ErrorMessage = "Unable to get current user data" + errUser.Error()
		_ = config.Renderer.Render(writer, "ResponseMessage", userInfo)
		return UserInfoRequest{}, ""
	}

	userInfo = config.getUserInformation(userID, userInfo, writer)
	return userInfo, UserID
}

func (config *Config) getUserInformation(userID string, userInfo UserInfoRequest, writer http.ResponseWriter) UserInfoRequest {
	userData, errUser := config.Db.GetUserById(context.Background(), userID)
	if errUser != nil {
		userInfo.ErrorMessage = fmt.Sprintf("Unable to get current user data. Error: %s", errUser.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		_ = config.Renderer.Render(writer, "ResponseMessage", userInfo)
	}
	userInfo.ID = userData.ID
	userInfo.UserName = userData.Name
	userInfo.IsAdmin = userData.IsAdmin.Int64
	return userInfo
}

func (config *Config) validateUniqueName(name string, userInfo *UserInfoRequest, writer http.ResponseWriter) {
	_, userError := config.Db.GetUserByName(context.Background(), name)
	if userError == nil {
		userInfo.ErrorMessage = "Name already registered, please try again"
		writer.WriteHeader(http.StatusBadRequest)
		_ = config.Renderer.Render(writer, "signup", userInfo)
		return
	}
}
