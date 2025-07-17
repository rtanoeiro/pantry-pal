package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"pantry-pal/pantry/database"

	"github.com/google/uuid"
)

func (config *Config) CreateUser(writer http.ResponseWriter, request *http.Request) {
	var returnResponse SuccessErrorResponse
	email := request.FormValue("email")
	name := request.FormValue("name")
	password := request.FormValue("password")

	_, userError := config.Db.GetUserByEmail(request.Context(), email)

	if userError == nil {
		returnResponse.ErrorMessage = "Email already registered, please try again"
		config.Renderer.Render(writer, "signup", returnResponse)
		return
	}

	hashedPassword, errPwd := HashPassword(password)
	if errPwd != nil {
		returnResponse.ErrorMessage = "Server error, please try again"
		config.Renderer.Render(writer, "signup", returnResponse)
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
		returnResponse.ErrorMessage = "Failed adding user to our systems, please try again"
		config.Renderer.Render(writer, "signup", returnResponse)
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
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	var UserPageData UserInfoRequest
	// TODO: Create fuction to redirect to a 401 page
	if errUser != nil {
		UserPageData.ErrorMessage = fmt.Sprintf("Unable to load user info. Error: %s", errUser.Error())
		config.Renderer.Render(writer, "user", UserPageData)
		return
	}
	userData, errUser := config.Db.GetUserById(request.Context(), userID)
	if errUser != nil {
		UserPageData.ErrorMessage = fmt.Sprintf("Unable to load user info. Error: %s", errUser.Error())
		config.Renderer.Render(writer, "user", UserPageData)
		return
	}
	UserPageData.ID = userData.ID
	UserPageData.UserName = userData.Name
	UserPageData.UserEmail = userData.Email
	UserPageData.IsAdmin = userData.IsAdmin.Int64

	config.GetAllOtherUsers(writer, request, &UserPageData)
	config.Renderer.Render(writer, "user", UserPageData)
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
				UserEmail:   user.Email,
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
	adminUserID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		userInfo.ErrorMessage = fmt.Sprintf("Unable to get current user data. Error: %s", errUser.Error())
		config.Renderer.Render(writer, "ResponseMessage", userInfo)
		return
	}
	userData, errUser := config.Db.GetUserById(request.Context(), adminUserID)
	if errUser != nil {
		userInfo.ErrorMessage = fmt.Sprintf("Unable to get current user data. Error: %s", errUser.Error())
		config.Renderer.Render(writer, "ResponseMessage", userInfo)
		return
	}
	userInfo.ID = userData.ID
	userInfo.UserName = userData.Name
	userInfo.UserEmail = userData.Email
	userInfo.IsAdmin = userData.IsAdmin.Int64
	userInfo.Users = []User{}

	userIDToDelete := request.PathValue("UserID")
	errDelete := config.Db.DeleteUser(request.Context(), userIDToDelete)
	if errDelete != nil {
		userInfo.ErrorMessage = fmt.Sprintf("Unable to delete user. Please try again. Error: %s", errDelete.Error())
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
	var userInfo CurrentUserRequest
	userID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		userInfo.ErrorMessage = fmt.Sprintf("Unable to get current user data. Error: %s", errUser.Error())
		config.Renderer.Render(writer, "ResponseMessage", userInfo)
		return
	}
	userData, errUser := config.Db.GetUserById(request.Context(), userID)
	if errUser != nil {
		userInfo.ErrorMessage = fmt.Sprintf("Unable to get current user data. Error: %s", errUser.Error())
		config.Renderer.Render(writer, "ResponseMessage", userInfo)
		return
	}
	userInfo.ID = userData.ID
	userInfo.UserName = userData.Name
	userInfo.UserEmail = userData.Email
	userInfo.IsAdmin = userData.IsAdmin.Int64

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
	userInfo *CurrentUserRequest,
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
	userInfo *CurrentUserRequest,
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
	userInfo *CurrentUserRequest,
	updateData string,
) {
	hashedPassword, errPWD := HashPassword(updateData)
	if errPWD != nil {
		userInfo.ErrorMessage = "Error on changing password, please try again"
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
	log.Printf("User %s has been assigned adming privileges at %s from %s", userToAddAdmin, time.Now(), returnUser.UserName)
	returnUser.SuccessMessage = "User Added as Admin with Success!"
	config.Renderer.Render(writer, "Admin", returnUser)
}

func (config *Config) RevokeUserAdmin(writer http.ResponseWriter, request *http.Request) {
	returnUser, userToRemoveAdmin := config.prepareUserAdmin(request, writer)
	errADmin := config.Db.RemoveUserAdmin(request.Context(), userToRemoveAdmin)
	if errADmin != nil {
		returnUser.ErrorMessage = "Unable to Remove Admin to user!"
		config.Renderer.Render(writer, "ResponseMessage", returnUser)
		return
	}
	config.GetAllOtherUsers(writer, request, &returnUser)
	log.Printf("User %s has lost adming privileges at %s from %s", userToRemoveAdmin, time.Now(), returnUser.UserName)
	returnUser.SuccessMessage = "User Removed as Admin with Success!"
	config.Renderer.Render(writer, "Admin", returnUser)
}

func (config *Config) prepareUserAdmin(request *http.Request, writer http.ResponseWriter) (UserInfoRequest, string) {
	UserID := request.PathValue("UserID")
	var userInfo UserInfoRequest
	AdminUserID, errUser := GetUserIDFromTokenAndValidate(request, config)
	if errUser != nil {
		userInfo.ErrorMessage = "Unable to get current user data" + errUser.Error()
		config.Renderer.Render(writer, "ResponseMessage", userInfo)
		return UserInfoRequest{}, ""
	}

	userData, errUser := config.Db.GetUserById(request.Context(), AdminUserID)
	if errUser != nil {
		userInfo.ErrorMessage = "Unable to get current user data" + errUser.Error()
		config.Renderer.Render(writer, "ResponseMessage", userInfo)
		return UserInfoRequest{}, ""
	}
	userInfo.ID = userData.ID
	userInfo.UserName = userData.Name
	userInfo.UserEmail = userData.Email
	userInfo.IsAdmin = 1
	return userInfo, UserID

}
