package handlers

import (
	"net/http"

	"github.com/Erzhan/weekend-warrior-backend/internal/db"
	"github.com/Erzhan/weekend-warrior-backend/internal/models"
	"github.com/gin-gonic/gin"
)

// CreateActivity godoc
// @Summary Создать новую активность
// @Description Принимает JSON и сохраняет активность в базу данных
// @Tags activities
// @Accept  json
// @Produce  json
// @Param activity body models.Activity true "Данные активности"
// @Success 201 {object} models.Activity
// @Router /activities [post]
func GetUserHandler(c *gin.Context) {
	user := models.User{
		Name:  "Erzhan",
		Email: "erzhan@example.com",
	}
	c.JSON(http.StatusOK, user)
}


func GetUserByIDHandler(c *gin.Context) {
    // 1. Достаем ID из параметров пути
    id := c.Param("id") 

    var user models.User

    // 2. Ищем юзера в базе по ID
    // .First() автоматически добавит LIMIT 1 и найдет по первичному ключу
    if err := db.DB.First(&user, id).Error; err != nil {
        // Если запись не найдена
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Пользователь не найден, u know?",
        })
        return
    }

    // 3. Если нашли — возвращаем данные
    c.JSON(http.StatusOK, gin.H{
        "user_data": user,
        "message":   "Данные успешно получены",
    })
}