package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
    Role      string    `bson:"role" json:"role"`
    Content   string    `bson:"content" json:"content"`
    Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

type Chat struct {
    ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    UserID         primitive.ObjectID `bson:"user_id" json:"user_id"`
    Title          string             `bson:"title" json:"title"`
    CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
    IsPinned       bool               `bson:"is_pinned" json:"is_pinned"`
    LastMessageAt  *time.Time         `bson:"last_message_at,omitempty" json:"last_message_at"`
    Messages       []Message          `bson:"messages" json:"messages"`
}

type CreateChatRequest struct {
    UserID  string `json:"user_id" binding:"required"`
    Message string `json:"message" binding:"required"`
}

type AddMessageRequest struct {
    ChatID  string `json:"chat_id" binding:"required"`
    Role    string `json:"role" binding:"required"`
    Content string `json:"content" binding:"required"`
}

type UpdateChatRequest struct {
    IsPinned *bool `json:"is_pinned"`
}

type ChatResponse struct {
    ID            string    `json:"id"`
    Title         string    `json:"title"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    IsPinned      bool      `json:"is_pinned"`
    LastMessageAt *time.Time `json:"last_message_at,omitempty"`
    MessageCount  int       `json:"message_count"`
}

type ChatDetailResponse struct {
    ID            string    `json:"id"`
    Title         string    `json:"title"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    IsPinned      bool      `json:"is_pinned"`
    LastMessageAt *time.Time `json:"last_message_at,omitempty"`
    Messages      []Message `json:"messages"`
}