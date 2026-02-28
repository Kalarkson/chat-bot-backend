package handlers

import (
    "net/http"
    "time"

    "chat-bot-backend/database"
    "chat-bot-backend/models"
    "chat-bot-backend/utils"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "go.mongodb.org/mongo-driver/mongo"
)

var jwtSecret = []byte("your-secret-key")

func Register(c *gin.Context) {
    var req models.RegisterRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные: " + err.Error()})
        return
    }

    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
        return
    }

    user := &models.User{
        Username: req.Username,
        Password: hashedPassword,
    }

    if err := database.CreateUser(user); err != nil {
        if mongoErr, ok := err.(*mongo.WriteError); ok && mongoErr.Code == 11000 {
            c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким именем уже существует"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания пользователя"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":  user.ID.Hex(),
        "username": user.Username,
        "exp":      time.Now().Add(24 * time.Hour).Unix(),
    })

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания токена"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Пользователь успешно создан",
        "token":   tokenString,
        "user": gin.H{
            "id":       user.ID.Hex(),
            "username": user.Username,
        },
    })
}

func Login(c *gin.Context) {
    var req models.LoginRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
        return
    }

    user, err := database.GetUserByUsername(req.Username)
    if err != nil || user == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
        return
    }

    if !utils.CheckPasswordHash(req.Password, user.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверное имя пользователя или пароль"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":  user.ID.Hex(),
        "username": user.Username,
        "exp":      time.Now().Add(24 * time.Hour).Unix(),
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
            "id":       user.ID.Hex(),
            "username": user.Username,
        },
    })
}

func GetProfile(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
        return
    }

    user, err := database.GetUserByID(userID.(string))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "user": gin.H{
            "id":       user.ID.Hex(),
            "username": user.Username,
        },
    })
}