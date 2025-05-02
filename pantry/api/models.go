package api

type UserAdd struct {
	Name     string
	Email    string
	Password string
}

// Probably replacable once we add JWT and login
type UserRequests struct {
	ID string
}

type ItemAdd struct {
	UserID   string `json:"user_id"`
	ItemName string `json:"item_name"`
	Quantity int64  `json:"quantity"`
	ExpiryAt string `json:"expiry_at"`
}

type ItemUpdate struct {
	ItemID            string
	UserID            string
	ItemName          string
	QuantityAvailable int64
	QuantityToAdd     int64
	ExpiryAt          string
}

type CreateUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateItemResponse struct {
	ItemID   string `json:"item_id"`
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

type AllItemsResponse struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	ItemName string `json:"item_name"`
	Quantity int64  `json:"quantity"`
	AddedAt  string `json:"added_at"`
	ExpiryAt string `json:"expiry_at"`
}
