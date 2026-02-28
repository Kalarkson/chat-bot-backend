package database

import (
    "context"
    "time"

    "chat-bot-backend/models"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var ChatsCollection *mongo.Collection

func InitCollections() {
    UsersCollection = Database.Collection("users")
    ChatsCollection = Database.Collection("chats")
}

func CreateChat(userID string, firstMessage string) (*models.Chat, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    userObjectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return nil, err
    }

    title := firstMessage
    if len(firstMessage) > 30 {
        title = firstMessage[:30] + "..."
    }

    now := time.Now()

    chat := &models.Chat{
        ID:        primitive.NewObjectID(),
        UserID:    userObjectID,
        Title:     title,
        CreatedAt: now,
        UpdatedAt: now,
        IsPinned:  false,
        Messages: []models.Message{
            {
                Role:      "user",
                Content:   firstMessage,
                Timestamp: now,
            },
        },
    }

    _, err = ChatsCollection.InsertOne(ctx, chat)
    if err != nil {
        return nil, err
    }

    return chat, nil
}

func AddMessage(chatID string, role string, content string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    chatObjectID, err := primitive.ObjectIDFromHex(chatID)
    if err != nil {
        return err
    }

    now := time.Now()

    newMessage := models.Message{
        Role:      role,
        Content:   content,
        Timestamp: now,
    }

    _, err = ChatsCollection.UpdateOne(
        ctx,
        bson.M{"_id": chatObjectID},
        bson.M{
            "$push": bson.M{"messages": newMessage},
            "$set":  bson.M{
                "updated_at":      now,
                "last_message_at": now,
            },
        },
    )

    return err
}

func GetUserChats(userID string) ([]models.ChatResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    userObjectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return nil, err
    }

    findOptions := options.Find().
        SetSort(bson.D{
            {Key: "is_pinned", Value: -1},
            {Key: "updated_at", Value: -1},
        }).
        SetProjection(bson.M{"messages": 0})

    cursor, err := ChatsCollection.Find(ctx, bson.M{"user_id": userObjectID}, findOptions)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var chats []models.Chat
    if err = cursor.All(ctx, &chats); err != nil {
        return nil, err
    }

    var response []models.ChatResponse
    for _, chat := range chats {
        response = append(response, models.ChatResponse{
            ID:            chat.ID.Hex(),
            Title:         chat.Title,
            CreatedAt:     chat.CreatedAt,
            UpdatedAt:     chat.UpdatedAt,
            IsPinned:      chat.IsPinned,
            LastMessageAt: chat.LastMessageAt,
            MessageCount:  len(chat.Messages),
        })
    }

    return response, nil
}

func GetChatByID(chatID string) (*models.ChatDetailResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    chatObjectID, err := primitive.ObjectIDFromHex(chatID)
    if err != nil {
        return nil, err
    }

    var chat models.Chat
    err = ChatsCollection.FindOne(ctx, bson.M{"_id": chatObjectID}).Decode(&chat)
    if err != nil {
        return nil, err
    }

    response := &models.ChatDetailResponse{
        ID:            chat.ID.Hex(),
        Title:         chat.Title,
        CreatedAt:     chat.CreatedAt,
        UpdatedAt:     chat.UpdatedAt,
        IsPinned:      chat.IsPinned,
        LastMessageAt: chat.LastMessageAt,
        Messages:      chat.Messages,
    }

    return response, nil
}

func DeleteChat(chatID string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    chatObjectID, err := primitive.ObjectIDFromHex(chatID)
    if err != nil {
        return err
    }

    _, err = ChatsCollection.DeleteOne(ctx, bson.M{"_id": chatObjectID})
    return err
}

func TogglePinChat(chatID string, pin bool) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    chatObjectID, err := primitive.ObjectIDFromHex(chatID)
    if err != nil {
        return err
    }

    _, err = ChatsCollection.UpdateOne(
        ctx,
        bson.M{"_id": chatObjectID},
        bson.M{
            "$set": bson.M{
                "is_pinned":  pin,
                "updated_at": time.Now(),
            },
        },
    )

    return err
}

func UpdateChat(chatID string, updateData bson.M) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    chatObjectID, err := primitive.ObjectIDFromHex(chatID)
    if err != nil {
        return err
    }

    updateData["updated_at"] = time.Now()

    _, err = ChatsCollection.UpdateOne(
        ctx,
        bson.M{"_id": chatObjectID},
        bson.M{"$set": updateData},
    )

    return err
}

func GetPinnedChats(userID string) ([]models.ChatResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    userObjectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return nil, err
    }

    findOptions := options.Find().
        SetSort(bson.M{"updated_at": -1}).
        SetProjection(bson.M{"messages": 0})

    cursor, err := ChatsCollection.Find(ctx, bson.M{
        "user_id":   userObjectID,
        "is_pinned": true,
    }, findOptions)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var chats []models.Chat
    if err = cursor.All(ctx, &chats); err != nil {
        return nil, err
    }

    var response []models.ChatResponse
    for _, chat := range chats {
        response = append(response, models.ChatResponse{
            ID:            chat.ID.Hex(),
            Title:         chat.Title,
            CreatedAt:     chat.CreatedAt,
            UpdatedAt:     chat.UpdatedAt,
            IsPinned:      chat.IsPinned,
            LastMessageAt: chat.LastMessageAt,
            MessageCount:  len(chat.Messages),
        })
    }

    return response, nil
}