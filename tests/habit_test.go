package tests

import (
	"bytes"
	"encoding/json"
	"habit-ai/internal/handler"
	"habit-ai/internal/models"
	"habit-ai/pkg"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

type HabitInput struct {
	Name          string `json:"name"`
	Category      string `json:"category"`
	PreferredHour string `json:"preferred_hour"`
	Frequency     int    `json:"frequency"`
}

var testJWT string
var testUserID string
var router *gin.Engine

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Erro ao carregar .env")
	}

	pkg.ConnectDB()
	models.Migrate()

	pkg.DB.Where("email = ?", "teste@teste.com").Delete(&models.User{})

	user := models.User{
		FullName: "Usuário Teste",
		Email:    "teste@teste.com",
	}
	user.SetPassword("123456")
	pkg.DB.Create(&user)
	testUserID = user.ID

	token, err := handler.GenerateTestJWT(testUserID)
	if err != nil {
		log.Fatal("Erro ao gerar JWT de teste:", err)
	}
	testJWT = token

	router = gin.Default()
	protected := router.Group("/")
	protected.Use(func(c *gin.Context) {
		c.Set("user_id", testUserID)
		c.Next()
	})
	{
		protected.POST("/habit", handler.CreateHabit)
		protected.GET("/habits", handler.GetHabits)
		protected.PUT("/habit/:id", handler.UpdateHabit)
		protected.DELETE("/habit/:id", handler.DeleteHabit)
		protected.POST("/habit/:id/complete", handler.CompleteHabit)
		protected.GET("/habit/:id/history", handler.GetHabitHistory)
		protected.GET("/habits/filter", handler.GetFilteredHabits)
	}

	os.Exit(m.Run())
}

func performRequest(method, path string, body any) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(jsonBody))
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set("Authorization", "Bearer "+testJWT)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestHabitCRUD(t *testing.T) {
	assert := assert.New(t)

	habitInput := HabitInput{
		Name:          "Beber Água",
		Category:      "Saúde",
		PreferredHour: "08:00",
		Frequency:     3,
	}
	resp := performRequest("POST", "/habit", habitInput)
	assert.Equal(201, resp.Code)

	var createdHabit models.Habit
	err := pkg.DB.Where("user_id = ?", testUserID).Order("created_at desc").First(&createdHabit).Error
	assert.NoError(err)

	resp = performRequest("GET", "/habits", nil)
	assert.Equal(200, resp.Code)

	updatedInput := HabitInput{
		Name:          "Beber Água Atualizado",
		Category:      "Saúde Atualizado",
		PreferredHour: "09:00",
		Frequency:     5,
	}
	updatePath := "/habit/" + createdHabit.ID
	resp = performRequest("PUT", updatePath, updatedInput)
	assert.Equal(200, resp.Code)

	var updatedHabit models.Habit
	err = pkg.DB.First(&updatedHabit, "id = ?", createdHabit.ID).Error
	assert.NoError(err)
	assert.Equal("Beber Água Atualizado", updatedHabit.Nome)

	completePath := "/habit/" + createdHabit.ID + "/complete"
	resp = performRequest("POST", completePath, nil)
	assert.Equal(200, resp.Code)

	var logs []models.HabitLog
	err = pkg.DB.Where("habit_id = ?", createdHabit.ID).Find(&logs).Error
	assert.NoError(err)
	assert.Len(logs, 1)

	historyPath := "/habit/" + createdHabit.ID + "/history"
	resp = performRequest("GET", historyPath, nil)
	assert.Equal(200, resp.Code)

	resp = performRequest("GET", "/habits/filter?category=Saúde&frequency=3", nil)
	assert.Equal(200, resp.Code)

	deletePath := "/habit/" + createdHabit.ID
	resp = performRequest("DELETE", deletePath, nil)
	assert.Equal(200, resp.Code)

	err = pkg.DB.First(&models.Habit{}, "id = ?", createdHabit.ID).Error
	assert.Error(err)
}