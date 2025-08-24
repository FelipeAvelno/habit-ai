package handler

import (
	"bytes"
	"habit-ai/internal/models"
	"habit-ai/pkg"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HabitInput struct {
	Name          string `json:"name" binding:"required"`
	Category      string `json:"category"`
	PreferredHour string `json:"preferred_hour"`
	Frequency     int    `json:"frequency"`
}

func getUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não autenticado"})
		return "", false
	}
	return userID.(string), true
}

func CreateHabit(c *gin.Context) {
	var input HabitInput

	body, _ := io.ReadAll(c.Request.Body)
	log.Println("Body recebido:", string(body))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	habit := models.Habit{
		Nome:            input.Name,
		Categoria:       input.Category,
		HorarioPreferido: input.PreferredHour,
		Frequencia:      input.Frequency,
		UserID:          userID,
		CreatedAt:       time.Now(),
	}

	if err := pkg.DB.Create(&habit).Error; err != nil {
		log.Println("Erro ao criar hábito:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Hábito criado com sucesso"})
}

func GetHabits(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	var habits []models.Habit
	if err := pkg.DB.Where("user_id = ?", userID).Find(&habits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, habits)
}

func UpdateHabit(c *gin.Context) {
	var input HabitInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	habitID := c.Param("id")
	var habit models.Habit
	if err := pkg.DB.Where("id = ? AND user_id = ?", habitID, userID).First(&habit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hábito não encontrado"})
		return
	}

	habit.Nome = input.Name
	habit.Categoria = input.Category
	habit.HorarioPreferido = input.PreferredHour
	habit.Frequencia = input.Frequency

	if err := pkg.DB.Save(&habit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hábito atualizado com sucesso"})
}

func DeleteHabit(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	habitID := c.Param("id")
	if err := pkg.DB.Where("id = ? AND user_id = ?", habitID, userID).Delete(&models.Habit{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hábito deletado com sucesso"})
}