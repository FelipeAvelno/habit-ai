package main

import (
	"habit-ai/internal/handler"
	"habit-ai/internal/models"
	"habit-ai/pkg"
	"habit-ai/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	pkg.ConnectDB()
	models.Migrate()

	r := gin.Default()

	// Rotas públicas
	r.POST("/register", handler.RegisterUser)
	r.POST("/login", handler.LoginUser)
	r.POST("/refresh-token", handler.RefreshToken)

	// Rotas protegidas por middleware JWT
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// CRUD de hábitos
		protected.POST("/habit", handler.CreateHabit)
		protected.GET("/habits", handler.GetHabits)
		protected.PUT("/habit/:id", handler.UpdateHabit)
		protected.DELETE("/habit/:id", handler.DeleteHabit)

		// Endpoints extras
		protected.POST("/habit/:id/complete", handler.CompleteHabit)
		protected.GET("/habit/:id/history", handler.GetHabitHistory)
		protected.GET("/habits/filter", handler.GetFilteredHabits) // filtro via query params
	}

	r.Run(":8080")
}