package main

import (
	"SIGE/internal/handlers"
	"SIGE/pkg/config"
	"SIGE/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	database.InitDB()

	c := gin.Default()

	c.POST("/upload", handlers.UploadImage)
	c.GET("/image/:file_id", handlers.GetDecryptedImage)
	c.Run(":8080")
}
