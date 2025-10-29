package database

import (
    "context"
    "time"

    "chat-bot-backend/models"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

var UsersCollection *mongo.Collection

// InitCollections инициализирует коллекции БД
func InitCollections() {
    UsersCollection = Database.Collection("users")
}

// CreateUser создает нового пользователя
func CreateUser(user *models.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Проверяем, существует ли пользователь с таким username
    existingUser, _ := GetUserByUsername(user.Username)
    if existingUser != nil {
        return &mongo.WriteError{Code: 11000, Message: "User already exists"}
    }

    // Устанавливаем роль по умолчанию
    user.Role = "user"

    // Вставляем документ в коллекцию
    _, err := UsersCollection.InsertOne(ctx, user)
    return err
}

// GetUserByUsername находит пользователя по имени
func GetUserByUsername(username string) (*models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user models.User
    // Ищем пользователя по username
    err := UsersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil // Пользователь не найден
        }
        return nil, err // Другая ошибка БД
    }
    return &user, nil
}

// GetUserByID находит пользователя по ID
func GetUserByID(id string) (*models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Конвертируем string ID в ObjectID MongoDB
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err // Невалидный ID
    }

    var user models.User
    err = UsersCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}