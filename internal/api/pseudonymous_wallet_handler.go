package api

import (
	"net/http"
	"verve/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterPseudonymousWalletRoutes(router *gin.Engine, service *services.PseudonymousWalletService) {
	router.POST("/api/pseudonymous_wallets", CreatePseudonymousWalletHandler(service))
}

func CreatePseudonymousWalletHandler(service *services.PseudonymousWalletService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, pub, priv, err := service.CreatePseudonymousWallet()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, CreatePseudonymousWalletResponse{
			WalletID:   id,
			PublicKey:  pub,
			PrivateKey: priv,
		})
	}
}
