package dto

import "github.com/yuki5155/go-google-auth/internal/domain/user"

// UserResponse represents user information in API responses
type UserResponse struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// FromDomain converts a domain User to a UserResponse DTO
func FromDomain(u *user.User) UserResponse {
	return UserResponse{
		ID:      u.ID().Value(),
		Email:   u.Email().Value(),
		Name:    u.Profile().Name(),
		Picture: u.Profile().Picture(),
	}
}

// NewUserResponse creates a UserResponse from individual fields
func NewUserResponse(id, email, name, picture string) UserResponse {
	return UserResponse{
		ID:      id,
		Email:   email,
		Name:    name,
		Picture: picture,
	}
}
