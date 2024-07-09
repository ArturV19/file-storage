package dto

type AuthenticateUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
