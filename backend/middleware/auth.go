package middleware

import (
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your-secret-key")

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Требуется авторизация"})
            c.Abort()
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(401, gin.H{"error": "Неверный формат токена"})
            c.Abort()
            return
        }

        tokenString := parts[1]

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "Невалидный токен"})
            c.Abort()
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            userID := claims["user_id"].(string)
            username := claims["username"].(string)
            
            c.Set("userID", userID)
            c.Set("username", username)
        } else {
            c.JSON(401, gin.H{"error": "Невалидный токен"})
            c.Abort()
            return
        }

        c.Next()
    }
}