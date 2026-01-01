package dto

// LoginResponse represents the response from a login operation
type LoginResponse struct {
	AccessToken  string       `json:"-"` // Not included in JSON, set as cookie
	RefreshToken string       `json:"-"` // Not included in JSON, set as cookie
	User         UserResponse `json:"user"`
	Message      string       `json:"message"`
}

// RefreshResponse represents the response from a token refresh operation
type RefreshResponse struct {
	AccessToken string `json:"-"` // Not included in JSON, set as cookie
	Message     string `json:"message"`
}

// LogoutResponse represents the response from a logout operation
type LogoutResponse struct {
	Message string `json:"message"`
}
