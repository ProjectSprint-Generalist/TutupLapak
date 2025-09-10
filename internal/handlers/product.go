package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"tutuplapak/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	db *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var product models.ProductInput
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid product input",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate file ID belongs to the user
	var fileUpload models.FileUpload
	if err := h.db.Where("id = ? and user_id = ?", product.FileID, userIDUint).First(&fileUpload).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "fileId is not valid / exists",
				Code:    http.StatusBadRequest,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Server Error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Ensure SKU is unique per user
	var existing models.Product
	if err := h.db.Where("user_id = ? AND sku = ?", userIDUint, product.SKU).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Success: false,
			Error:   "sku already exists",
			Code:    http.StatusConflict,
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Server Error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	sku := strings.TrimSpace(product.SKU)

	p := models.Product{
		UserID:           userIDUint,
		Name:             product.Name,
		Category:         product.Category,
		Qty:              product.Qty,
		Price:            product.Price,
		SKU:              sku,
		FileID:           product.FileID,
		FileURI:          fileUpload.FileURI,
		// FileThumbnailURI: "", // let Go generate zero value
	}

	if err := h.db.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Server Error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	resp := models.ProductOutput{
		ProductID:        strconv.FormatUint(uint64(p.ID), 10),
		Name:             p.Name,
		Category:         string(p.Category),
		Quantity:         p.Qty,
		Price:            p.Price,
		SKU:              p.SKU,
		FileID:           p.FileID,
		FileURI:          p.FileURI,
		FileThumbnailURI: p.FileThumbnailURI,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
	}

	c.JSON(http.StatusCreated, resp)
}
