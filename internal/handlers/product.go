package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"tutuplapak/internal/models"
)

type ProductHandler struct {
	db *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	// Ambil productId sesuai kontrak
	productId := c.Param("productId")

	// Bind + validasi request
	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation error", "details": err.Error()})
		return
	}

	// Cari produk
	var product models.Product
	if err := h.db.First(&product, "id = ?", productId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "productId is not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	// Cek SKU conflict (per account)
	var conflict models.Product
	if err := h.db.Where("sku = ? AND id <> ?", req.Sku, productId).First(&conflict).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "sku already exists"})
		return
	}

	// Validasi fileId (misal cek storage service)
	if !isValidFileId(req.FileId) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fileId is not valid / exists"})
		return
	}

	// Update product
	product.Name = req.Name
	product.Category = req.Category
	product.Qty = req.Qty
	product.Price = req.Price
	product.Sku = req.Sku
	product.FileID = req.FileId
	// FileURI & FileThumbnailURI bisa diambil dari file service

	if err := h.db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	resp := models.ProductResponse{
		ProductID:        strconv.Itoa(product.ID),
		Name:             product.Name,
		Category:         product.Category,
		Qty:              product.Qty,
		Price:            product.Price,
		Sku:              product.Sku,
		FileID:           product.FileID,
		FileURI:          product.FileURI,
		FileThumbnailURI: product.FileThumbnailURI,
		CreatedAt:        product.CreatedAt,
		UpdatedAt:        product.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// dummy validator
func isValidFileId(fileId string) bool {
	// implementasi cek file service / DB
	return fileId != ""
}
