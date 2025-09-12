package handlers

import (
    "net/http"

    "tutuplapak/internal/models"
    "tutuplapak/internal/services"

    "github.com/gin-gonic/gin"
)

type FileHandler struct {
    minioService *services.MinIOService
}

func NewFileHandler(minioService *services.MinIOService) *FileHandler {
    return &FileHandler{
        minioService: minioService,
    }
}

func (h *FileHandler) UploadFile(c *gin.Context) {
    // Public endpoint - no authentication required
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Error: "File is required",
        })
        return
    }

    // Open the file for validation and processing
    src, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{
            Error: "Failed to read file",
        })
        return
    }
    defer src.Close()

    // Pass 0 for userID since this is a public endpoint
    response, err := h.minioService.UploadFile(file, 0)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{
            Error: err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, response)
}