package models

// Message represents a message in the database
type Message struct {
	MessageId       int    `json:"messageId"`
	MessageForUserId int    `json:"messageForUserId"`
	Comment         string `json:"comment"`
	IsActive        bool   `json:"isActive"`
	IsRead          bool   `json:"isRead"`
	CreatedBy       int    `json:"createdBy"`
	CreatedAt       string `json:"createdAt"`
}

// ChatParticipant represents a user you've chatted with
type ChatParticipant struct {
	UserId            int     `json:"userId"`
	UserName          string  `json:"userName"`
	ProfilePic        *string `json:"profilePic"`
	RoleId            int     `json:"roleId"`
	LastMessageText   string  `json:"lastMessageText"`
	LastMessageTime   string  `json:"lastMessageTime"`
	LastMessageSentBy int     `json:"lastMessageSentBy"`
	UnreadCount       int     `json:"unreadCount"`
}

// ChatMessage represents a message with sender details
type ChatMessage struct {
	MessageId       int     `json:"messageId"`
	MessageForUserId int     `json:"messageForUserId"`
	Comment         string  `json:"comment"`
	IsActive        bool    `json:"isActive"`
	IsRead          bool    `json:"isRead"`
	CreatedBy       int     `json:"createdBy"`
	CreatedByName   string  `json:"createdByName"`
	CreatedByPic    *string `json:"createdByPic"`
	CreatedAt       string  `json:"createdAt"`
	IsSentByMe      bool    `json:"isSentByMe"`
}

// CreateMessageRequest represents the request to create a message
type CreateMessageRequest struct {
	MessageForUserId int    `json:"messageForUserId"`
	Comment         string `json:"comment"`
}

// UpdateMessageRequest represents the request to update a message
type UpdateMessageRequest struct {
	MessageId int    `json:"messageId"`
	Comment   string `json:"comment"`
}

// DeleteMessageRequest represents the request to delete a message
type DeleteMessageRequest struct {
	MessageId int `json:"messageId"`
}

// AvailableChatUser represents a user available for starting a chat
type AvailableChatUser struct {
	UserId     int     `json:"userId"`
	UserName   string  `json:"userName"`
	Mobile     string  `json:"mobile"`
	ProfilePic *string `json:"profilePic"`
	RoleId     int     `json:"roleId"`
}
