package api

type UserAdd struct {
	Name     string
	Email    string
	Password string
}

type CreateUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
