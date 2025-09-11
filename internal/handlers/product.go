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
		UserID:   userIDUint,
		Name:     product.Name,
		Category: product.Category,
		Qty:      product.Qty,
		Price:    product.Price,
		SKU:      sku,
		FileID:   product.FileID,
		FileURI:  fileUpload.FileURI,
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

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	// Ambil user_id dari context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "User not authenticated",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Invalid user ID",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	// Ambil productId dari URL
	productIdStr := c.Param("productId")
	productIdUint, err := strconv.ParseUint(productIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid productId",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Bind & validasi request
	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Validation error: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Cari produk milik user
	var product models.Product
	if err := h.db.Where("id = ? AND user_id = ?", productIdUint, userIDUint).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success: false,
				Error:   "productId not found",
				Code:    http.StatusNotFound,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Server error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Cek SKU conflict (per user)
	var conflict models.Product
	if err := h.db.Where("user_id = ? AND sku = ? AND id <> ?", userIDUint, req.SKU, productIdUint).First(&conflict).Error; err == nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Success: false,
			Error:   "sku already exists",
			Code:    http.StatusConflict,
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Server error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Validasi fileId: apakah fileId milik user
	fileIdUint := req.FileID

	var fileUpload models.FileUpload
	if err := h.db.Where("id = ? AND user_id = ?", fileIdUint, userIDUint).First(&fileUpload).Error; err != nil {
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
			Error:   "Server error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Update product
	product.Name = req.Name
	product.Category = models.ProductCategory(req.Category)
	product.Qty = req.Qty
	product.Price = req.Price
	product.SKU = strings.TrimSpace(req.SKU)
	product.FileID = uint(fileIdUint)
	product.FileURI = fileUpload.FileURI
	// FileThumbnailURI bisa diisi kalau ada service thumbnail

	if err := h.db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Server error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Response sesuai kontrak
	resp := models.ProductResponse{
		ProductID:        strconv.FormatUint(uint64(product.ID), 10),
		Name:             product.Name,
		Category:         string(product.Category),
		Qty:              product.Qty,
		Price:            product.Price,
		SKU:              product.SKU,
		FileID:           product.FileID,
		FileURI:          product.FileURI,
		FileThumbnailURI: product.FileThumbnailURI,
		CreatedAt:        product.CreatedAt,
		UpdatedAt:        product.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}


// DeleteProduct DELETE /v1/product/productId
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Success: false,
			Error:   "Expired / invalid / missing request token",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	productIdStr := c.Param("productId")
	productId, err := strconv.ParseUint(productIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Invalid product ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var product models.Product
	if err := h.db.Where("id = ? AND user_id = ?", productId, userID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Success: false,
				Error:   "ProductId is not found",
				Code:    http.StatusNotFound,
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

	if err := h.db.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Server Error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Product deleted successfully",
	})
}
