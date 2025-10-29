package handlers

import (
    "net/http"
    "time"

    "chat-bot-backend/database"
    "chat-bot-backend/models"
    "chat-bot-backend/utils"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "go.mongodb.org/mongo-driver/mongo" // 👈 ДОБАВЬ ЭТОТ ИМПОРТ
)

// JWT секретный ключ (будет из конфига)
var jwtSecret = []byte("your-secret-key")

// Register обработчик регистрации
func Register(c *gin.Context) {
    var req models.RegisterRequest
    
    // Валидируем входящие данные
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
        return
    }

    // Хэшируем пароль
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
        return
    }

    // Создаем пользователя
    user := &models.User{
        Username: req.Username,
        Password: hashedPassword,
    }

    // Сохраняем в БД
    if err := database.CreateUser(user); err != nil {
        // Обрабатываем ошибку дубликата username
        if mongoErr, ok := err.(*mongo.WriteError); ok && mongoErr.Code == 11000 {
            c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким именем уже существует"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания пользователя"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Пользователь успешно создан",
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
        },
    })
}

// Login обработчик входа
func Login(c *gin.Context) {
    var req models.LoginRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
        return
    }

    // Ищем пользователя в БД
    user, err := database.GetUserByUsername(req.Username)
    if err != nil || user == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
        return
    }

    // Проверяем пароль
    if !utils.CheckPasswordHash(req.Password, user.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
        return
    }

    // Создаем JWT токен
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":  user.ID.Hex(),
        "username": user.Username,
        "exp":      time.Now().Add(24 * time.Hour).Unix(), // Токен на 24 часа
    })

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания токена"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Успешный вход",
        "token":   tokenString,
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
        },
    })
}

// GetProfile получение профиля пользователя
func GetProfile(c *gin.Context) {
    // userID из middleware (пока заглушка)
    userID := c.GetString("userID")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
        return
    }

    user, err := database.GetUserByID(userID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
        },
    })
}