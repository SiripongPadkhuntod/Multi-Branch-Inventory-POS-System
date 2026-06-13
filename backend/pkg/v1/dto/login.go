package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Code         string `json:"code"`
	Description  string `json:"description"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         any    `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}
