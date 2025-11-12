package routes

import (
	"SIGE/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterImageRoutes(r *gin.RouterGroup, handler handlers.ImageHandler) {
	r.POST("/upload", handler.UploadImage)
	r.GET("/image/:file_id", handler.GetDecryptedImage)
}
