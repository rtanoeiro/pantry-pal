package api

import (
	"html/template"
	"pantry-pal/pantry/database"
	"time"
)

type User struct {
	UserID      string
	UserName    string
	IsUserAdmin int64
}

type Templates struct {
	templates *template.Template
}

type CreateUserRequest struct {
	Name     string
	Password string
}

type CreateUserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserInfoRequest struct {
	ID             string `json:"id"`
	UserName       string `json:"name"`
	IsAdmin        int64  `json:"is_admin"`
	Users          []User `json:"users"`
	ErrorMessage   string `json:"error_message"`
	SuccessMessage string `json:"success_message"`
}

type LoginUserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
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
	Items          []PantryItem `json:"items"`
	SuccessMessage string       `json:"success_message"`
	ErrorMessage   string       `json:"error_message"`
}

type PantryStats struct {
	ExpiringSoon   []PantryItem
	SuccessMessage string
	ErrorMessage   string
}

type SuccessErrorResponse struct {
	SuccessMessage string `json:"success_message"`
	ErrorMessage   string `json:"error_message"`
}

type CartInfo struct {
	CartItems      []database.CartItem `json:"cart_items"`
	ErrorMessage   string              `json:"success_message"`
	SuccessMessage string              `json:"error_message"`
}
