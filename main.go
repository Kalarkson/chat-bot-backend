package main

import (
    "log"

    "chat-bot-backend/config"
    "chat-bot-backend/database"
    "chat-bot-backend/handlers"
    "chat-bot-backend/middleware"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {
    config.LoadConfig()

    if err := database.ConnectDB(); err != nil {
        log.Fatal("❌ Ошибка подключения к БД:", err)
    }
    defer database.DisconnectDB()

    r := gin.Default()
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    }))

    public := r.Group("/api/auth")
    {
        public.POST("/register", handlers.Register)
        public.POST("/login", handlers.Login)
    }

    chats := r.Group("/api/chats")
    {
        chats.POST("/", handlers.CreateChatHandler)
        chats.POST("/message", handlers.AddMessageHandler)
        chats.GET("/user/:userID", handlers.GetUserChatsHandler)
        chats.GET("/pinned/:userID", handlers.GetPinnedChatsHandler) 
        chats.GET("/:chatID", handlers.GetChatHandler)
        chats.PUT("/:chatID", handlers.UpdateChatHandler)           
        chats.PUT("/:chatID/pin", handlers.TogglePinChatHandler)    
        chats.DELETE("/:chatID", handlers.DeleteChatHandler)
    }

    protected := r.Group("/api")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/profile", handlers.GetProfile)
    }

    r.OPTIONS("/*path", func(c *gin.Context) {
        c.Status(200)
    })

    port := ":" + config.AppConfig.Port
    log.Printf("🚀 Сервер запущен на порту %s", port)
    
    if err := r.Run(port); err != nil {
        log.Fatal("❌ Ошибка запуска сервера:", err)
    }
}