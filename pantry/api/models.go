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
	UserID   string
	ItemName string
	Quantity int64
	ExpiryAt string
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
	UserID   string `json:"user_id"`
	ItemName string `json:"item_name"`
	Quantity int64  `json:"quantity"`
	ExpiryAt string `json:"expiry_at"`
}

type AddItemResponse struct {
	UserID   string `json:"user_id"`
	ItemName string `json:"item_name"`
	Quantity int64  `json:"quantity"`
	ExpiryAt string `json:"expiry_at"`
}
