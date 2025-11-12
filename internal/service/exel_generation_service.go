package service

import (
	"SIGE/internal/models"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExcelGenerationService interface {
	JSONToExcel(req models.JSONToExcelRequest) (string, error)
}
type excelGenerationService struct{}

func NewExcelGenerationService() ExcelGenerationService {
	return &excelGenerationService{}
}

func (s *excelGenerationService) JSONToExcel(req models.JSONToExcelRequest) (string, error) {

	if len(req.ID) != len(req.Path) || len(req.ID) != len(req.CreatedAt) {
		return "", fmt.Errorf("arrays must be of the same length")
	}

	f := excelize.NewFile()
	sheet := f.GetSheetName(0)

	// Заголовки
	f.SetCellValue(sheet, "A1", "ID")
	f.SetCellValue(sheet, "B1", "Path")
	f.SetCellValue(sheet, "C1", "Created At")

	// Заполняем строки
	for i := 0; i < len(req.ID); i++ {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", i+2), req.ID[i])
		f.SetCellValue(sheet, fmt.Sprintf("B%d", i+2), req.Path[i])
		f.SetCellValue(sheet, fmt.Sprintf("C%d", i+2), req.CreatedAt[i])
	}

	// Генерируем безопасный путь для сохранения
	months := map[string]string{
		"January":   "Января",
		"February":  "Февраля",
		"March":     "Марта",
		"April":     "Апреля",
		"May":       "Мая",
		"June":      "Июня",
		"July":      "Июля",
		"August":    "Августа",
		"September": "Сентября",
		"October":   "Октября",
		"November":  "Ноября",
		"December":  "Декабря",
	}

	now := time.Now()
	engMonth := now.Format("January")
	rusMonth := months[engMonth]
	monthYear := fmt.Sprintf("%s %d", rusMonth, now.Year())

	dir := filepath.Join("uploads", monthYear)
	dir = filepath.Clean(dir) // нормализация пути

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Println("Ошибка создания директории:", err)
		return "", err
	}

	fileName := fmt.Sprintf("excel_dynamic_%d.xlsx", now.Unix())
	rawPath := filepath.Join(dir, fileName)
	savePath := filepath.Clean(rawPath)

	// Сохраняем файл
	if err := f.SaveAs(savePath); err != nil {
		log.Println("Ошибка сохранения Excel:", err)
		return "", err
	}

	log.Println("Excel успешно создан:", savePath)
	return savePath, nil
}
