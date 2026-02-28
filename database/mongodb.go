package database

import (
    "context"
    "log"
    "time"

    "chat-bot-backend/config"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
    Client   *mongo.Client
    Database *mongo.Database
)

func ConnectDB() error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.MongoURI))
    if err != nil {
        return err
    }

    if err = client.Ping(ctx, readpref.Primary()); err != nil {
        return err
    }

    Client = client
    Database = client.Database(config.AppConfig.Database)

    InitCollections()

    log.Println("Connected to MongoDB!")
    return nil
}

func DisconnectDB() {
    if Client != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        Client.Disconnect(ctx)
        log.Println("Disconnected from MongoDB")
    }
}