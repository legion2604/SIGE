package service

import (
	"SIGE/pkg/crypto"
	"fmt"
	"io"
	"log"
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
	// Открываем файл
	src, err := file.Open()
	if err != nil {
		log.Println("Ошибка открытия файла:", err)
		return "", "", err
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		log.Println("Ошибка чтения файла:", err)
		return "", "", err
	}

	// Генерируем AES-ключ
	key, err := crypto.GenerateAESKey()
	if err != nil {
		log.Println("Ошибка генерации AES ключа:", err)
		return "", "", err
	}

	// Шифруем данные
	ciphertext, iv, err := crypto.EncryptAES(data, key)
	if err != nil {
		log.Println("Ошибка шифрования:", err)
		return "", "", err
	}

	// Генерируем уникальное имя файла
	fileID := uuid.New().String()
	rawPath := filepath.Join("uploads", "encrypted", fileID+".bin")

	// Нормализуем путь
	savePath := filepath.Clean(rawPath)
	if !filepath.HasPrefix(savePath, filepath.Join("uploads", "encrypted")) {
		err := fmt.Errorf("недопустимый путь к файлу")
		log.Println(err)
		return "", "", err
	}

	// Создаём директорию, если её нет
	if err := os.MkdirAll(filepath.Dir(savePath), os.ModePerm); err != nil {
		log.Println("Ошибка создания директории:", err)
		return "", "", err
	}

	// Сохраняем файл (IV + зашифрованные данные)
	finalData := append(iv, ciphertext...)
	if err := os.WriteFile(savePath, finalData, 0644); err != nil {
		log.Println("Ошибка записи файла:", err)
		return "", "", err
	}

	// Создаём JWT, в который зашиваем ключ
	token, err := crypto.GenerateJWT(userID, crypto.EncodeKeyToBase64(key))
	if err != nil {
		log.Println("Ошибка генерации JWT:", err)
		return "", "", err
	}

	log.Println("Файл успешно зашифрован и сохранён:", savePath)
	return fileID, token, nil
}

func (s *imageEncryptService) DecryptImage(authHeader string, fileID string) ([]byte, error) {
	// Убираем префикс Bearer
	if len(authHeader) < len("Bearer ")+1 {
		err := fmt.Errorf("некорректный Authorization header")
		log.Println(err)
		return nil, err
	}
	token := authHeader[len("Bearer "):]

	// Парсим JWT
	claims, err := crypto.ParseJWT(token)
	if err != nil {
		log.Println("Ошибка парсинга JWT:", err)
		return nil, err
	}

	// Декодируем ключ
	key, err := crypto.DecodeKeyFromBase64(claims.Key)
	if err != nil {
		log.Println("Ошибка декодирования ключа из JWT:", err)
		return nil, err
	}

	// Формируем путь к файлу
	rawPath := filepath.Join("uploads", "encrypted", fileID+".bin")
	filePath := filepath.Clean(rawPath)

	// Проверяем, что файл находится внутри uploads/encrypted
	if !filepath.HasPrefix(filePath, filepath.Join("uploads", "encrypted")) {
		err := fmt.Errorf("недопустимый путь к файлу")
		log.Println(err)
		return nil, err
	}

	// Читаем зашифрованный файл
	encryptedData, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Ошибка чтения файла:", err)
		return nil, err
	}

	// Расшифровываем
	decrypted, err := crypto.DecryptAES(encryptedData, key)
	if err != nil {
		log.Println("Ошибка расшифровки файла:", err)
		return nil, err
	}

	log.Println("Файл успешно расшифрован:", filePath)
	return decrypted, nil
}
