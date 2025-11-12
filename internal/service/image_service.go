package service

import (
	"SIGE/pkg/crypto"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type ImageService interface {
	EncryptImage(file *multipart.FileHeader, userID int) (string, string, error)
	DecryptImage(authHeader string, fileID string) ([]byte, error)
}

type imageEncryptService struct{}

func NewImageService() ImageService {
	return &imageEncryptService{}
}

func (s *imageEncryptService) EncryptImage(file *multipart.FileHeader, userID int) (string, string, error) {
	src, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return "", "", err
	}

	// Генерируем AES-ключ
	key, err := crypto.GenerateAESKey()
	if err != nil {
		return "", "", err
	}

	// Шифруем данные
	ciphertext, iv, err := crypto.EncryptAES(data, key)
	if err != nil {
		return "", "", err
	}

	// Генерируем уникальное имя файла
	fileID := uuid.New().String()
	savePath := filepath.Join("uploads", "encrypted", fileID+".bin")

	// Создаём директорию, если её нет
	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
		return "", "", err
	}

	// Сохраняем файл (IV + зашифрованные данные)
	finalData := append(iv, ciphertext...)
	if err := os.WriteFile(savePath, finalData, 0644); err != nil {
		return "", "", err
	}

	// Создаём JWT, в который зашиваем ключ
	token, err := crypto.GenerateJWT(userID, crypto.EncodeKeyToBase64(key))
	if err != nil {
		return "", "", err
	}

	return fileID, token, nil
}

func (s *imageEncryptService) DecryptImage(authHeader string, fileID string) ([]byte, error) {
	token := authHeader[len("Bearer "):]
	claims, err := crypto.ParseJWT(token)
	if err != nil {
		return nil, err
	}

	key, _ := crypto.DecodeKeyFromBase64(claims.Key)
	filePath := filepath.Join("uploads", "encrypted", fileID+".bin")

	encryptedData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	decrypted, err := crypto.DecryptAES(encryptedData, key)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}
