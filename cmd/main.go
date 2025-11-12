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

	exelGenerationService := service.NewExcelGenerationService()
	exelGenerationHandler := handlers.NewExelGenerationHandler(exelGenerationService)

	api := c.Group("/")
	{
		routes.RegisterImageRoutes(api, imageHandler)
		routes.RegisterExelGenerationRoutes(api, exelGenerationHandler)
	}

	c.Run(":8080")
}
