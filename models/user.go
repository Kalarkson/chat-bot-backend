package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// User представляет упрощенную модель пользователя
type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`    // Уникальный ID MongoDB
    Username string             `bson:"username" json:"username"`   // Имя пользователя (уникальное)
    Password string             `bson:"password" json:"-"`          // Хэш пароля (скрыт в JSON)
    Role     string             `bson:"role" json:"role"`           // Роль пользователя
}

// LoginRequest модель для запроса входа
type LoginRequest struct {
    Username string `json:"username" binding:"required"`  // Обязательное поле
    Password string `json:"password" binding:"required"`  // Обязательное поле
}

// RegisterRequest модель для запроса регистрации
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"` // Валидация длины
    Password string `json:"password" binding:"required,min=6"`        // Минимум 6 символов
}