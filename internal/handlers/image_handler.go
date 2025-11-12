package handlers

import (
	"SIGE/internal/service"
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ImageHandler interface {
	UploadImage(c *gin.Context)
	GetDecryptedImage(c *gin.Context)
}

type imageHandler struct {
	service service.ImageService
}

func NewImageHandler(s service.ImageService) ImageHandler {
	return &imageHandler{service: s}
}

func (h *imageHandler) UploadImage(cxt *gin.Context) {
	file, err := cxt.FormFile("file")
	if err != nil {
		cxt.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	userID, _ := strconv.Atoi(cxt.DefaultPostForm("user_id", "0"))

	fileID, token, err := h.service.EncryptImage(file, userID)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cxt.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "File encrypted and saved",
		"file_id": fileID,
		"token":   token,
	})
}

func (h *imageHandler) GetDecryptedImage(cxt *gin.Context) {
	fileID := cxt.Param("file_id")
	authHeader := cxt.GetHeader("Authorization")
	if authHeader == "" {
		cxt.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}
	decrypted, err := h.service.DecryptImage(authHeader, fileID)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cxt.DataFromReader(http.StatusOK, int64(len(decrypted)), "image/jpeg", io.NopCloser(bytes.NewReader(decrypted)), nil)
}
