package routes

import (
	"SIGE/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterExelGenerationRoutes(r *gin.RouterGroup, handler handlers.ExelGenerationHandler) {
	r.POST("/json-to-excel", handler.JSONToExcel)
}
