package handlers

import (
	"SIGE/pkg/crypto"
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	userID, _ := strconv.Atoi(c.DefaultPostForm("user_id", "0"))

	// Генерируем AES-ключ
	key, _ := crypto.GenerateAESKey()

	// Читаем файл
	src, _ := file.Open()
	defer src.Close()
	data, _ := io.ReadAll(src)

	// Шифруем
	ciphertext, iv, _ := crypto.EncryptAES(data, key)
	finalData := append(iv, ciphertext...) // IV + зашифрованные данные
	os.WriteFile("uploads/encrypted/file.bin", finalData, 0644)

	// Генерируем уникальное имя
	fileID := uuid.New().String()
	savePath := filepath.Join("uploads", "encrypted", fileID+".bin")

	os.MkdirAll(filepath.Dir(savePath), os.ModePerm)
	os.WriteFile(savePath, ciphertext, 0644)

	// Создаём JWT
	token, _ := crypto.GenerateJWT(userID, crypto.EncodeKeyToBase64(key))

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "File encrypted and saved",
		"file_id": fileID,
		"token":   token,
	})
}

func GetDecryptedImage(c *gin.Context) {
	fileID := c.Param("file_id")
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	token := authHeader[len("Bearer "):]
	claims, err := crypto.ParseJWT(token)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid token"})
		return
	}

	key, _ := crypto.DecodeKeyFromBase64(claims.Key)
	filePath := filepath.Join("uploads", "encrypted", fileID+".bin")

	encryptedData, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	decrypted, err := crypto.DecryptAES(encryptedData, key)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "decryption failed"})
		return
	}

	c.DataFromReader(http.StatusOK, int64(len(decrypted)), "image/jpeg", io.NopCloser(bytes.NewReader(decrypted)), nil)
}
