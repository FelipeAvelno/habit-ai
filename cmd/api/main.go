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

	// Rotas protegidas por middleware JWT
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/habit", handler.CreateHabit)
		protected.GET("/habits", handler.GetHabits)
		protected.PUT("/habit/:id", handler.UpdateHabit)
		protected.DELETE("/habit/:id", handler.DeleteHabit)
	}

	r.Run(":8080")
}