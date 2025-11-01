package handlers

import (
    "net/http"

    "chat-bot-backend/database"
    "chat-bot-backend/models"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
)

func CreateChatHandler(c *gin.Context) {
    var req models.CreateChatRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные: " + err.Error()})
        return
    }

    chat, err := database.CreateChat(req.UserID, req.Message)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания чата"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Чат создан",
        "chat": gin.H{
            "id":              chat.ID.Hex(),
            "title":           chat.Title,
            "created_at":      chat.CreatedAt,
            "updated_at":      chat.UpdatedAt,
            "is_pinned":       chat.IsPinned,
            "last_message_at": chat.LastMessageAt,
            "message_count":   1,
        },
    })
}

func AddMessageHandler(c *gin.Context) {
    var req models.AddMessageRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
        return
    }

    if req.Role != "user" && req.Role != "assistant" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Роль должна быть 'user' или 'assistant'"})
        return
    }

    err := database.AddMessage(req.ChatID, req.Role, req.Content)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления сообщения"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Сообщение добавлено",
    })
}

func GetUserChatsHandler(c *gin.Context) {
    userID := c.Param("userID")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID пользователя обязателен"})
        return
    }

    chats, err := database.GetUserChats(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения чатов"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "chats": chats,
    })
}

func GetChatHandler(c *gin.Context) {
    chatID := c.Param("chatID")
    if chatID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID чата обязателен"})
        return
    }

    chat, err := database.GetChatByID(chatID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Чат не найден"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "chat": chat,
    })
}

func DeleteChatHandler(c *gin.Context) {
    chatID := c.Param("chatID")
    if chatID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID чата обязателен"})
        return
    }

    err := database.DeleteChat(chatID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления чата"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Чат удален",
    })
}

func TogglePinChatHandler(c *gin.Context) {
    chatID := c.Param("chatID")
    if chatID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID чата обязателен"})
        return
    }

    var req models.UpdateChatRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
        return
    }

    if req.IsPinned == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Поле is_pinned обязательно"})
        return
    }

    err := database.TogglePinChat(chatID, *req.IsPinned)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления чата"})
        return
    }

    action := "закреплен"
    if !*req.IsPinned {
        action = "откреплен"
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Чат " + action,
    })
}

func UpdateChatHandler(c *gin.Context) {
    chatID := c.Param("chatID")
    if chatID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID чата обязателен"})
        return
    }

    var req map[string]interface{}
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
        return
    }

    updateData := bson.M{}
    
    if isPinned, exists := req["is_pinned"]; exists {
        if pinned, ok := isPinned.(bool); ok {
            updateData["is_pinned"] = pinned
        }
    }

    if title, exists := req["title"]; exists {
        if titleStr, ok := title.(string); ok && titleStr != "" {
            if len(titleStr) > 30 {
                updateData["title"] = titleStr[:30] + "..."
            } else {
                updateData["title"] = titleStr
            }
        }
    }

    if len(updateData) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Нет данных для обновления"})
        return
    }

    err := database.UpdateChat(chatID, updateData)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления чата"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Чат обновлен",
    })
}

func GetPinnedChatsHandler(c *gin.Context) {
    userID := c.Param("userID")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID пользователя обязателен"})
        return
    }

    chats, err := database.GetPinnedChats(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения чатов"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "chats": chats,
    })
}