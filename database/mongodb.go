package database

import (
    "context"
    "log"
    "time"

    "chat-bot-backend/config"  // Импортируем наш конфиг

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

// Глобальные переменные для подключения к БД
var (
    Client   *mongo.Client
    Database *mongo.Database
)

// ConnectDB устанавливает подключение к MongoDB
func ConnectDB() error {
    // Создаем контекст с таймаутом 10 секунд
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel() // Гарантируем освобождение ресурсов

    // Подключаемся к MongoDB
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.MongoURI))
    if err != nil {
        return err
    }

    // Проверяем подключение пингом
    if err = client.Ping(ctx, readpref.Primary()); err != nil {
        return err
    }

    // Сохраняем подключение в глобальных переменных
    Client = client
    Database = client.Database(config.AppConfig.Database)

    // Инициализируем коллекции
    InitCollections()

    log.Println("✅ Успешное подключение к MongoDB!")
    return nil
}

// DisconnectDB закрывает подключение к БД
func DisconnectDB() {
    if Client != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        Client.Disconnect(ctx)
        log.Println("✅ Отключение от MongoDB")
    }
}