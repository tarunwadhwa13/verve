package api

import (
	"net/http"
	"strconv"
	"verve/internal/api/middleware"
	"verve/internal/services"

	"github.com/gin-gonic/gin"
)

// RegisterWalletRoutes sets up the wallet-related routes
// @Summary Register wallet routes
// @Description Register all wallet-related routes under user namespace
// @Tags wallets
func RegisterWalletRoutes(router *gin.Engine, walletService *services.WalletService) {
	userRoutes := router.Group("/api/user/:id")
	userRoutes.Use(middleware.AuthMiddleware())
	{
		walletRoutes := userRoutes.Group("/wallets")
		walletRoutes.POST("", CreateWalletHandler(walletService))
		walletRoutes.GET("", GetUserWalletsHandler(walletService))
		walletRoutes.GET("/:wallet_id", GetWalletHandler(walletService))
	}
}

// CreateWalletHandler creates a new wallet for a user
// @Summary Create wallet
// @Description Create a new wallet for the authenticated user
// @Tags wallets
// @Accept json
// @Produce json
// @Param id path integer true "User ID"
// @Param wallet body CreateWalletRequest true "Wallet details"
// @Success 201 {object} models.Wallet
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden - Can only create wallets for yourself"
// @Security ApiKeyAuth
// @Router /user/{id}/wallets [post]
func CreateWalletHandler(walletService *services.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Ensure user can only create wallets for themselves
		requestingUserID := c.GetInt("userID")
		if requestingUserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only create wallets for yourself"})
			return
		}

		var req CreateWalletRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		wallet, err := walletService.CreateWallet(userID, req.Currency)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create wallet"})
			return
		}
		c.JSON(http.StatusCreated, wallet)
	}
}

// GetUserWalletsHandler retrieves all wallets for a user
// @Summary Get user wallets
// @Description Get all wallets belonging to the authenticated user
// @Tags wallets
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {array} models.Wallet
// @Failure 400 {object} ErrorResponse "Invalid user ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden - Can only view your own wallets"
// @Security ApiKeyAuth
// @Router /user/{id}/wallets [get]
func GetUserWalletsHandler(walletService *services.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Ensure user can only view their own wallets
		requestingUserID := c.GetInt("userID")
		if requestingUserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own wallets"})
			return
		}

		wallets, err := walletService.GetWalletsByUserID(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch wallets"})
			return
		}
		c.JSON(http.StatusOK, wallets)
	}
}

// GetWalletHandler retrieves a specific wallet
// @Summary Get wallet details
// @Description Get details of a specific wallet belonging to the authenticated user
// @Tags wallets
// @Produce json
// @Param user_id path integer true "User ID"
// @Param wallet_id path integer true "Wallet ID"
// @Success 200 {object} models.Wallet
// @Failure 400 {object} ErrorResponse "Invalid wallet ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden - Can only view your own wallet"
// @Failure 404 {object} ErrorResponse "Wallet not found"
// @Security ApiKeyAuth
// @Router /user/{user_id}/wallets/{wallet_id} [get]
func GetWalletHandler(walletService *services.WalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		walletID, err := strconv.ParseInt(c.Param("wallet_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet ID"})
			return
		}

		// Ensure user can only view their own wallet
		requestingUserID := c.GetInt("userID")
		if requestingUserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own wallet"})
			return
		}

		wallet, err := walletService.GetWalletByID(walletID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Wallet not found"})
			return
		}

		// Double check the wallet belongs to the user
		if wallet.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "This wallet does not belong to you"})
			return
		}

		c.JSON(http.StatusOK, wallet)
	}
}
