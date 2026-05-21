package middleware

import (
	"net/http"
	"strings"

	// проверь путь к своему пакету auth
	"github.com/Erzhan/weekend-warrior-backend/internal/auth"
	"github.com/gin-gonic/gin"
)


func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Достаем заголовок Authorization
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен отсутствует. Войдите в аккаунт"})
            c.Abort() 
            return
        }

        // Отсекаем префикс "Bearer "
        parts := strings.SplitN(authHeader, " ", 2)
        if !(len(parts) == 2 && parts[0] == "Bearer") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат токена"})
            c.Abort()
            return
        }

        tokenString := parts[1]
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Пустой токен"})
            c.Abort()
            return
        }

        // 🔥 НАЧАЛО МАГИИ: Парсим токен и достаемClaims (данные внутри токена)
        // Замени auth.ParseToken на название твоей функции парсинга (например, utils.ParseToken)
        claims, err := auth.ParseToken(tokenString) 
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный или просроченный токен"})
            c.Abort()
            return
        }

        // 🔥 САМАЯ ВАЖНАЯ СТРОЧКА: Записываем ID пользователя в контекст Gin!
        // Проверь, как называется поле ID в твоей структуре Claims (обычно UserID или ID)
        c.Set("user_id", claims.UserID) 

        c.Next() // Теперь хэндлер CreateActivity сможет прочитать этот "userID"!
    }
}