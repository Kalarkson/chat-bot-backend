package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port      string
    Cors      string
    JWTSecret string
    MongoURI  string
    Database  string
}

var AppConfig *Config

func LoadConfig() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    AppConfig = &Config{
        Cors:      getEnv("CORS", "http://localhost:4003"),
        Port:      getEnv("PORT", "4002"),
        JWTSecret: getEnv("JWT_SECRET", "fallback-secret-key-change-in-production"),
        MongoURI:  getEnv("MONGODB_URI", "mongodb://Kalarkson:qwerty-chat-bot@chat-bot-mongo:27017/"),
        Database:  getEnv("DATABASE_NAME", "chat_bot"),
    }

    log.Println("Configuration loaded successfully")
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}