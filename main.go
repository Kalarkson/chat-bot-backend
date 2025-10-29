package main

import (
    "log"

    "chat-bot-backend/config"
    "chat-bot-backend/database"
    "chat-bot-backend/handlers"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors" // 👈 ДОБАВЬ ЭТОТ ИМПОРТ
)

func main() {
    // Загружаем конфигурацию
    config.LoadConfig()

    // Подключаемся к БД
    if err := database.ConnectDB(); err != nil {
        log.Fatal("❌ Ошибка подключения к БД:", err)
    }
    defer database.DisconnectDB()

    // Создаем Gin роутер
    r := gin.Default()

    // 👈 ДОБАВЬ CORS middleware ПЕРЕД маршрутами
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"}, // URL твоего фронтенда
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    }))

    // Группа публичных маршрутов
    public := r.Group("/api/auth")
    {
        public.POST("/register", handlers.Register)
        public.POST("/login", handlers.Login)
    }

    // Группа защищенных маршрутов
    protected := r.Group("/api")
    {
        protected.GET("/profile", handlers.GetProfile)
    }

    // 👈 ДОБАВЬ обработчик для OPTIONS запросов
    r.OPTIONS("/*path", func(c *gin.Context) {
        c.Status(200)
    })

    // Запускаем сервер
    port := ":" + config.AppConfig.Port
    log.Printf("🚀 Сервер запущен на порту %s", port)
    
    if err := r.Run(port); err != nil {
        log.Fatal("❌ Ошибка запуска сервера:", err)
    }
}