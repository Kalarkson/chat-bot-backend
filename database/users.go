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

func CreateUser(user *models.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    existingUser, _ := GetUserByUsername(user.Username)
    if existingUser != nil {
        return &mongo.WriteError{Code: 11000, Message: "User already exists"}
    }

    user.CreatedAt = time.Now()
    user.Role = "user"

    _, err := UsersCollection.InsertOne(ctx, user)
    return err
}

func GetUserByUsername(username string) (*models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var user models.User
    err := UsersCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}

func GetUserByID(id string) (*models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    var user models.User
    err = UsersCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}