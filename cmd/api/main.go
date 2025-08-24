package main

import (
	"habit-ai/internal/handler"
	"habit-ai/internal/models"
	"habit-ai/pkg"
	"github.com/gin-gonic/gin"
)

func main() {
	pkg.ConnectDB()

	models.Migrate()

	r := gin.Default()

	r.POST("/register", handler.RegisterUser)
	r.POST("/login", handler.LoginUser)

	r.Run(":8080")
}