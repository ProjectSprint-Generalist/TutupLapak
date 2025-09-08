package middleware

import (
	"net/http"

	"tutuplapak/internal/models"

	"github.com/gin-gonic/gin"
)

// Recovery returns a gin.HandlerFunc for recovering from panics
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   "Internal server error: " + err,
				Code:    http.StatusInternalServerError,
			})
		}
		c.Abort()
	})
}
