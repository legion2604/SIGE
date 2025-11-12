package handlers

import (
	"SIGE/internal/models"
	"SIGE/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExelGenerationHandler interface {
	JSONToExcel(cxt *gin.Context)
}

type exelGenerationHandler struct {
	service service.ExcelGenerationService
}

func NewExelGenerationHandler(s service.ExcelGenerationService) ExelGenerationHandler {
	return &exelGenerationHandler{service: s}
}

func (h *exelGenerationHandler) JSONToExcel(cxt *gin.Context) {
	var req models.JSONToExcelRequest
	if err := cxt.ShouldBindJSON(&req); err != nil {
		cxt.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}
	res, err := h.service.JSONToExcel(req)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}
	cxt.JSON(http.StatusOK, gin.H{"status": "success", "path": res})
}
