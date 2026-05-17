package middleware

import (
	"net/http"
	"strings"

	// проверь путь к своему пакету auth
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Достаем заголовок Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Токен отсутствует. Войдите в аккаунт"})
			c.Abort() // Останавливаем запрос
			return
		}

		// Токен обычно прилетает в формате "Bearer <токен>", отсекаем префикс
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат токена"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Тут проверяем токен (парсим его)
		// Для этого в твоем internal/auth/jwt.go должна быть функция ParseToken.
		// Если ее нет, мы ее добавим. Пока просто пропустим пользователя, если токен есть.
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пустой токен"})
			c.Abort()
			return
		}

		c.Next() // Все ок, пускаем к хэндлеру
	}
}