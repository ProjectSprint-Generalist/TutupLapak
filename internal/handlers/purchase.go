package handlers

import (
	"net/http"
	"strconv"
	"time"
	"tutuplapak/internal/models"
	"tutuplapak/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PurchaseHandler struct {
	db *gorm.DB
}

func NewPurchaseHandler(db *gorm.DB) *PurchaseHandler {
	return &PurchaseHandler{db: db}
}

func (h *PurchaseHandler) PurchaseProducts(c *gin.Context) {
	var req models.PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	if len(req.PurchasedItems) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "No items to purchase",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if req.SenderContactType == models.ContactTypePhone {
		if validationErr := utils.PhoneValidation(req.SenderContactDetail); validationErr != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   validationErr.Error,
				Code:    http.StatusBadRequest,
			})
			return
		}
	} else if req.SenderContactType == models.ContactTypeEmail {
		if err := utils.EmailValidation(req.SenderContactDetail); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   err.Error(),
				Code:    http.StatusBadRequest,
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "Sender contact type must be 'phone' or 'email'",
			Code:    http.StatusBadRequest,
		})
		return
	}

	var productIDs []uint
	productQuantityMap := make(map[uint]uint)

	for _, item := range req.PurchasedItems {
		if item.Quantity < 2 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "Minimum quantity is 2",
				Code:    http.StatusBadRequest,
			})
			return
		}

		productID, err := strconv.ParseUint(item.ProductID, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "Invalid product ID",
				Code:    http.StatusBadRequest,
			})
			return
		}

		productIDs = append(productIDs, uint(productID))
		productQuantityMap[uint(productID)] = item.Quantity
	}

	// Batch fetch all products in one query
	var products []models.Product
	if err := h.db.Where("id IN ?", productIDs).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Database error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	if len(products) != len(productIDs) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success: false,
			Error:   "One or more products not found",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Create product map and validate inventory
	productMap := make(map[uint]models.Product)
	for _, product := range products {
		productMap[product.ID] = product

		requestedQty := productQuantityMap[product.ID]
		if product.Qty < requestedQty {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Success: false,
				Error:   "Insufficient product quantity for product ID " + strconv.FormatUint(uint64(product.ID), 10),
				Code:    http.StatusBadRequest,
			})
			return
		}
	}

	// NOTE: use transaction for multiple related operations (this is to handle racing condition actually)
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to start transaction",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Decrease product stock based on requested quantities
	for _, product := range products {
		requestedQty := productQuantityMap[product.ID]
		if err := tx.Model(&product).Update("qty", product.Qty-requestedQty).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   "Failed to update product quantity",
				Code:    http.StatusInternalServerError,
			})
			return
		}
	}

	sellerIDs := make(map[uint]bool)
	for _, product := range products {
		sellerIDs[product.UserID] = true
	}

	var sellerIDList []uint
	for sellerID := range sellerIDs {
		sellerIDList = append(sellerIDList, sellerID)
	}

	// Batch fetch required sellers
	var sellers []models.User
	if err := h.db.Where("id IN ?", sellerIDList).Find(&sellers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to fetch seller information",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	sellerMap := make(map[uint]models.User)
	for _, seller := range sellers {
		sellerMap[seller.ID] = seller
	}

	// Build response data
	var purchasedItems []models.PurchasedItemResponse
	var totalPrice uint
	sellerPaymentMap := make(map[uint]models.SellerPaymentInfo)
	var purchaseItemsToCreate []models.PurchaseItem

	for _, item := range req.PurchasedItems {
		productID, _ := strconv.ParseUint(item.ProductID, 10, 32)
		product := productMap[uint(productID)]

		itemTotalPrice := product.Price * item.Quantity
		totalPrice += itemTotalPrice

		purchasedItem := models.PurchasedItemResponse{
			ProductID:        item.ProductID,
			Name:             product.Name,
			Category:         string(product.Category),
			Qty:              item.Quantity,
			Price:            product.Price,
			SKU:              product.SKU,
			FileID:           strconv.FormatUint(uint64(product.FileID), 10),
			FileURI:          product.FileURI,
			FileThumbnailURI: product.FileThumbnailURI,
			CreatedAt:        product.CreatedAt.Format(time.RFC3339),
			UpdatedAt:        product.UpdatedAt.Format(time.RFC3339),
		}
		purchasedItems = append(purchasedItems, purchasedItem)

		purchaseItemsToCreate = append(purchaseItemsToCreate, models.PurchaseItem{
			ProductID: uint(productID),
			Quantity:  item.Quantity,
			Price:     product.Price,
		})

		if seller, exists := sellerMap[product.UserID]; exists {
			if _, exists := sellerPaymentMap[product.UserID]; !exists {
				sellerPaymentMap[product.UserID] = models.SellerPaymentInfo{
					BankAccountName:   seller.BankAccountName,
					BankAccountHolder: seller.BankAccountHolder,
					BankAccountNumber: seller.BankAccountNumber,
					TotalPrice:        itemTotalPrice,
				}
			} else {
				paymentInfo := sellerPaymentMap[product.UserID]
				paymentInfo.TotalPrice += itemTotalPrice
				sellerPaymentMap[product.UserID] = paymentInfo
			}
		}
	}

	var paymentDetails []models.SellerPaymentInfo
	for _, payment := range sellerPaymentMap {
		paymentDetails = append(paymentDetails, payment)
	}

	purchase := models.Purchase{
		SenderName:          req.SenderName,
		SenderContactType:   req.SenderContactType,
		SenderContactDetail: req.SenderContactDetail,
		TotalPrice:          totalPrice,
	}

	if err := tx.Create(&purchase).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to create purchase",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	purchaseID := purchase.ID

	// prepare bulk create purchase items
	for i := range purchaseItemsToCreate {
		purchaseItemsToCreate[i].PurchaseID = purchase.ID
	}

	if err := tx.Create(&purchaseItemsToCreate).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to create purchase items",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success: false,
			Error:   "Failed to complete purchase transaction",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	response := models.PurchaseResponse{
		PurchaseID:     purchaseID,
		PurchasedItems: purchasedItems,
		TotalPrice:     totalPrice,
		PaymentDetails: paymentDetails,
	}

	c.JSON(http.StatusCreated, response)
}

// ProcessPurchase handles POST /v1/purchase/:purchaseId
// It accepts a list of fileIds (payment proof image ids) and stores them
// associated to the purchase. It validates that:
// - caller is authenticated
// - purchase exists
// - provided fileIds exist and are owned by the caller
// - number of fileIds equals the number of distinct sellers in the purchase

/*
Flow:
- Accepts array of fileIds
-
*/
func (h *PurchaseHandler) ProcessPurchase(c *gin.Context) {
	// Auth
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "User not authenticated"})
		return
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	purchaseID := c.Param("purchaseId")
	if purchaseID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "purchaseId is required", Code: http.StatusBadRequest})
		return
	}

	var req models.ProcessPurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "Invalid request body", Code: http.StatusBadRequest})
		return
	}

	if len(req.FileIDs) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "fileIds is required", Code: http.StatusBadRequest})
		return
	}

	var purchase models.Purchase
	if err := h.db.Where("id = ?", purchaseID).First(&purchase).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Success: false, Error: "Purchase not found", Code: http.StatusNotFound})
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: "Database error", Code: http.StatusInternalServerError})
		return
	}

	var items []models.PurchaseItem
	if err := h.db.Where("purchase_id = ?", purchase.ID).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: "Failed to load purchase items", Code: http.StatusInternalServerError})
		return
	}

	if len(items) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "Purchase has no items", Code: http.StatusBadRequest})
		return
	}

	// Build product -> seller map and distinct sellers
	productIDs := make([]uint, 0, len(items))
	for _, it := range items {
		productIDs = append(productIDs, it.ProductID)
	}

	var products []models.Product
	if err := h.db.Where("id IN ?", productIDs).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: "Failed to load products", Code: http.StatusInternalServerError})
		return
	}

	sellerSet := map[uint]struct{}{}

	for _, p := range products {
		sellerSet[p.UserID] = struct{}{}
	}

	expectedProofs := len(sellerSet)
	if expectedProofs == 0 {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: "No sellers found for purchase items", Code: http.StatusInternalServerError})
		return
	}

	if len(req.FileIDs) != expectedProofs {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "fileIds count does not match required proofs", Code: http.StatusBadRequest})
		return
	}

	// Validate fileIds: must exist and belong to caller
	var fileIDsUint []uint
	for _, s := range req.FileIDs {
		id64, err := strconv.ParseUint(s, 10, 32)
		if err != nil || id64 == 0 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "fileIds must be valid numeric ids", Code: http.StatusBadRequest})
			return
		}
		fileIDsUint = append(fileIDsUint, uint(id64))
	}

	var files []models.FileUpload
	if err := h.db.Where("id IN ? AND user_id = ?", fileIDsUint, userIDUint).Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: "Failed to validate fileIds", Code: http.StatusInternalServerError})
		return
	}

	if len(files) != len(fileIDsUint) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Success: false, Error: "One or more file IDs are invalid, do not exist, or are not owned by the user", Code: http.StatusBadRequest})
		return
	}

	// Persist proofs in a transaction
	tx := h.db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: "Failed to start transaction", Code: http.StatusInternalServerError})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, fid := range fileIDsUint {
		proof := models.PurchasePaymentProof{PurchaseID: purchase.ID, FileUploadID: fid}
		if err := tx.Create(&proof).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: "Failed to save payment proofs", Code: http.StatusInternalServerError})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Success: false, Error: "Failed to commit transaction", Code: http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true})
}
