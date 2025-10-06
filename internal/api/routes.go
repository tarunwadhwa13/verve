package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	api := router.Group("/api")
	{
		api.GET("/health", HealthCheckHandler(db))
	}
}

func HealthCheckHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "db not connected"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}
