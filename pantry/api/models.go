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

// To be used when we use javascript
type LoginUserRequest struct {
	Email    string
	Password string
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

// To Be used in HomePage
type PantryItem struct {
	ItemName string `json:"item_name"`
	Quantity int    `json:"quantity"`
	ExpiryAt string `json:"expiry_at"`
}

type ItemShopping struct {
	ItemName string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

type PantryStats struct {
	NumItems     int
	ExpiringSoon []PantryItem
	ShoppingList []ItemShopping
}
