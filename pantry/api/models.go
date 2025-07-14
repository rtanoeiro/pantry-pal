package api

import (
	"html/template"
	"time"
)

type User struct {
	UserID    string
	Name      string
	Email     string
	UserAdmin int64
}

type Templates struct {
	templates *template.Template
}

// To be used when we use javascript
type CreateUserRequest struct {
	Name     string
	Email    string
	Password string
}

type CreateUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserInfoRequest struct {
	ID             string `json:"id"`
	UserName       string `json:"name"`
	UserEmail      string `json:"email"`
	IsAdmin        bool   `json:"is_admin"`
	Users          []User `json:"users"`
	ErrorMessage   string `json:"error_message"`
	SuccessMessage string `json:"success_message"`
}

type LoginUserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	JWTToken  *string   `json:"token,omitempty"`
}

type AddItemRequest struct {
	UserID   string `json:"user_id"`
	ItemName string `json:"item_name"`
	Quantity int64  `json:"quantity"`
	ExpiryAt string `json:"expiry_at"`
}

type AddItemResponse struct {
	ItemID   string `json:"item_id"`
	UserID   string `json:"user_id"`
	ItemName string `json:"item_name"`
	Quantity int64  `json:"quantity"`
	ExpiryAt string `json:"expiry_at"`
}

type UpdateItemRequest struct {
	ItemID            string
	UserID            string
	ItemName          string
	QuantityAvailable int64
	QuantityToAdd     int64
	ExpiryAt          string
}

type UpdateItemResponse struct {
	ItemID   string `json:"item_id"`
	UserID   string `json:"user_id"`
	ItemName string `json:"item_name"`
	Quantity int64  `json:"quantity"`
	ExpiryAt string `json:"expiry_at"`
}

type PantryItem struct {
	ItemID   string `json:"item_id"`
	ItemName string `json:"item_name"`
	Quantity int    `json:"quantity"`
	ExpiryAt string `json:"expiry_at"`
}

type PantryItems struct {
	SuccessMessage string       `json:"success_message"`
	ErrorMessage   string       `json:"error_message"`
	Items          []PantryItem `json:"items"`
}

type PantryStats struct {
	ExpiringSoon []PantryItem
}

type SuccessErrorResponse struct {
	SuccessMessage string `json:"success_message"`
	ErrorMessage   string `json:"error_message"`
}
