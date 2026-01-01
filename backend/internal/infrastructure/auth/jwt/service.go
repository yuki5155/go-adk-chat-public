package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yuki5155/go-google-auth/internal/application/ports"
)

// Service handles JWT token operations and implements ports.TokenGenerator
type Service struct {
	secretKey          []byte
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

// tokenClaims represents the internal JWT claims structure
type tokenClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Picture   string `json:"picture"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// NewService creates a new JWT Service instance
func NewService(secretKey string) *Service {
	return &Service{
		secretKey:          []byte(secretKey),
		accessTokenExpiry:  15 * time.Minute,   // Access token expires in 15 minutes
		refreshTokenExpiry: 7 * 24 * time.Hour, // Refresh token expires in 7 days
	}
}

// GenerateTokenPair generates both access and refresh tokens
func (s *Service) GenerateTokenPair(user ports.UserInfo) (accessToken, refreshToken string, err error) {
	accessToken, err = s.generateToken(user, "access", s.accessTokenExpiry)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.generateToken(user, "refresh", s.refreshTokenExpiry)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// generateToken creates a JWT token with the specified claims
func (s *Service) generateToken(user ports.UserInfo, tokenType string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := tokenClaims{
		UserID:    user.UserID,
		Email:     user.Email,
		Name:      user.Name,
		Picture:   user.Picture,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "go-google-auth",
			Subject:   user.UserID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken validates a JWT token and returns the claims
func (s *Service) validateToken(tokenString string) (*tokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ports.ErrInvalidToken
		}
		return s.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ports.ErrExpiredToken
		}
		return nil, ports.ErrInvalidToken
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return nil, ports.ErrInvalidToken
	}

	return claims, nil
}

// ValidateAccessToken validates an access token and returns the claims
func (s *Service) ValidateAccessToken(tokenString string) (*ports.TokenClaims, error) {
	claims, err := s.validateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "access" {
		return nil, ports.ErrInvalidToken
	}

	// Convert to ports.TokenClaims
	return &ports.TokenClaims{
		UserID:  claims.UserID,
		Email:   claims.Email,
		Name:    claims.Name,
		Picture: claims.Picture,
	}, nil
}

// validateRefreshToken validates a refresh token
func (s *Service) validateRefreshToken(tokenString string) (*tokenClaims, error) {
	claims, err := s.validateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, ports.ErrInvalidToken
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token from a valid refresh token
func (s *Service) RefreshAccessToken(refreshTokenString string) (string, error) {
	claims, err := s.validateRefreshToken(refreshTokenString)
	if err != nil {
		return "", err
	}

	user := ports.UserInfo{
		UserID:  claims.UserID,
		Email:   claims.Email,
		Name:    claims.Name,
		Picture: claims.Picture,
	}

	return s.generateToken(user, "access", s.accessTokenExpiry)
}

// GetAccessTokenExpiry returns the access token expiry duration in seconds
func (s *Service) GetAccessTokenExpiry() int {
	return int(s.accessTokenExpiry.Seconds())
}

// GetRefreshTokenExpiry returns the refresh token expiry duration in seconds
func (s *Service) GetRefreshTokenExpiry() int {
	return int(s.refreshTokenExpiry.Seconds())
}
