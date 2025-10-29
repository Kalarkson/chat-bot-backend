package utils

import (
    "golang.org/x/crypto/bcrypt"
)

// HashPassword создает bcrypt хэш из пароля
func HashPassword(password string) (string, error) {
    // bcrypt.DefaultCost = 10 (хороший баланс безопасности и производительности)
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// CheckPasswordHash проверяет пароль против хэша
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil // true если пароль верный
}