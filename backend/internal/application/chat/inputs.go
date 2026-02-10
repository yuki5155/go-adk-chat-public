package chat

// CreateThreadCommand represents the command to create a new thread
type CreateThreadCommand struct {
	UserID string
	Title  string
	Model  string
}

// SendMessageCommand represents the command to send a message
type SendMessageCommand struct {
	UserID   string
	ThreadID string
	Content  string
}

// DeleteThreadCommand represents the command to delete a thread
type DeleteThreadCommand struct {
	UserID   string
	ThreadID string
}

// GetThreadQuery represents the query to get a thread with messages
type GetThreadQuery struct {
	UserID   string
	ThreadID string
	Limit    int
}

// ListThreadsQuery represents the query to list threads
type ListThreadsQuery struct {
	UserID  string
	Limit   int
	LastKey string
}
