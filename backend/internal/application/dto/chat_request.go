package dto

// CreateThreadRequest represents a request to create a new chat thread
type CreateThreadRequest struct {
	Title string `json:"title"`
	Model string `json:"model"`
}

// SendMessageRequest represents a request to send a message to a thread
type SendMessageRequest struct {
	Content string `json:"content" binding:"required"`
}
