package api

import (
	"net/http"
	"strconv"
	"verve/internal/api/middleware"
	"verve/internal/services"

	"github.com/gin-gonic/gin"
) // RegisterTransferRoutes sets up the routes for transfers.
// @Summary Register transfer routes
// @Description Register all transfer-related routes
// @Tags transfers
func RegisterTransferRoutes(router *gin.Engine, transferService *services.TransferService) {
	transferRoutes := router.Group("/api/transfer")
	transferRoutes.Use(middleware.AuthMiddleware())
	{
		transferRoutes.POST("", InitiateTransferHandler(transferService))
		transferRoutes.GET("/:id", GetTransferStatusHandler(transferService))
	}
}

// InitiateTransferHandler handles the creation of a new transfer.
// @Summary Initiate a transfer
// @Description Create a new transfer between wallets
// @Tags transfers
// @Accept json
// @Produce json
// @Param transfer body TransferRequest true "Transfer details"
// @Success 200 {object} models.Transfer
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Insufficient funds or invalid PIN"
// @Security ApiKeyAuth
// @Router /transfer [post]
func InitiateTransferHandler(transferService *services.TransferService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req TransferRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID := c.GetInt("userID")

		transfer, err := transferService.InitiateTransfer(
			userID,
			req.SenderWalletID,
			req.ReceiverWalletID,
			req.Amount,
			req.IsAnonymous,
			req.Pin,
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, transfer)
	}
}

// GetTransferStatusHandler handles checking the status of a transfer.
// GetTransferStatusHandler retrieves the status of a transfer
// @Summary Get transfer status
// @Description Get the current status of a transfer
// @Tags transfers
// @Produce json
// @Param id path integer true "Transfer ID"
// @Success 200 {object} models.Transfer
// @Failure 400 {object} ErrorResponse "Invalid transfer ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Transfer not found"
// @Security ApiKeyAuth
// @Router /transfer/{id} [get]
func GetTransferStatusHandler(transferService *services.TransferService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transfer ID"})
			return
		}

		transfer, err := transferService.GetTransferStatus(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "transfer not found"})
			return
		}

		c.JSON(http.StatusOK, transfer)
	}
}
