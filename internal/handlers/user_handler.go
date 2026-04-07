package handlers

import (
	"net/http"

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
	// Достаем ID из URL (все, что после /user/)
	id := c.Param("id") 

	user := models.User{
		Name:  "User Name",
		Email: "user@example.com",
	}

	c.JSON(http.StatusOK, gin.H{
		"requested_id": id,
		"user_data":    user,
		"message":      "Данные успешно получены",
	})
}