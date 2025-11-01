package models

import (
    "time"
    
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Username  string             `bson:"username" json:"username"`
    Password  string             `bson:"password" json:"-"`
    Role      string             `bson:"role" json:"role"`
    CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type UserResponse struct {
    ID        string    `json:"id"`
    Username  string    `json:"username"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Password string `json:"password" binding:"required,min=6"`
}