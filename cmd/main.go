package main

import (
	"SIGE/internal/handlers"
	"SIGE/internal/routes"
	"SIGE/internal/service"
	"SIGE/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	c := gin.Default()

	imageService := service.NewImageService()
	imageHandler := handlers.NewImageHandler(imageService)

	api := c.Group("api")
	{
		routes.RegisterImageRoutes(api, imageHandler)
	}

	c.Run(config.GetEnv("PORT"))
}
