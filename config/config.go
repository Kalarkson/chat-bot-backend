package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port      string
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
        Port:      getEnv("PORT", "8080"),
        JWTSecret: getEnv("JWT_SECRET", "fallback-secret-key-change-in-production"),
        MongoURI:  getEnv("MONGODB_URI", "mongodb://localhost:27017"),
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