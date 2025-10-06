package api

import (
	"net/http"
	"strconv"
	"verve/internal/api/middleware"
	"verve/internal/models"
	"verve/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterBadgeRoutes sets up the badge-related routes
// @Summary Register badge routes
// @Description Register all badge-related routes with authentication and authorization
// @Tags badges
func RegisterBadgeRoutes(router *gin.Engine, badgeService *services.BadgeService) {
	// Enable Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Admin routes under /badges
	adminBadgeRoutes := router.Group("/api/badges")
	adminBadgeRoutes.Use(middleware.AuthMiddleware())
	{
		adminBadgeRoutes.POST("", middleware.RoleMiddleware("admin"), CreateBadgeHandler(badgeService))
		adminBadgeRoutes.PUT("/:id", middleware.RoleMiddleware("admin"), UpdateBadgeHandler(badgeService))
		adminBadgeRoutes.POST("/:id/award/:user_id", middleware.RoleMiddleware("admin"), AwardBadgeHandler(badgeService))
	}

	// Public badge listing and details
	publicBadgeRoutes := router.Group("/api/badges")
	publicBadgeRoutes.Use(middleware.AuthMiddleware())
	{
		publicBadgeRoutes.GET("", ListBadgesHandler(badgeService))
		publicBadgeRoutes.GET("/:id", GetBadgeHandler(badgeService))
		publicBadgeRoutes.GET("/:id/holders", GetBadgeHoldersHandler(badgeService))
	}

	// User-specific badge routes under /user/{user_id}/badges
	userRoutes := router.Group("/api/user/:id")
	userRoutes.Use(middleware.AuthMiddleware())
	{
		badgeRoutes := userRoutes.Group("/badges")
		badgeRoutes.GET("", GetUserBadgesHandler(badgeService))
	}
}

// CreateBadgeHandler creates a new badge
// @Summary Create a new badge
// @Description Create a new badge with optional achievement rules
// @Tags badges
// @Accept json
// @Produce json
// @Param badge body CreateBadgeRequest true "Badge details"
// @Success 201 {object} models.Badge "Created badge"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Security ApiKeyAuth
// @Router /badges [post]
func CreateBadgeHandler(badgeService *services.BadgeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name        string                   `json:"name" binding:"required"`
			Description string                   `json:"description" binding:"required"`
			IconURL     string                   `json:"icon_url"`
			Points      int                      `json:"points"`
			Rules       []models.AchievementRule `json:"rules"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate each rule if provided
		for _, rule := range req.Rules {
			if err := badgeService.ValidateAchievementRule(&rule); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule: " + err.Error()})
				return
			}
		}

		createdBy := c.GetInt("userID")
		badge, err := badgeService.CreateBadge(
			req.Name,
			req.Description,
			req.IconURL,
			req.Points,
			createdBy,
			req.Rules,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create badge"})
			return
		}

		c.JSON(http.StatusCreated, badge)
	}
}

// UpdateBadgeHandler updates an existing badge
// @Summary Update a badge
// @Description Update an existing badge's details
// @Tags badges
// @Accept json
// @Produce json
// @Param id path integer true "Badge ID"
// @Param badge body UpdateBadgeRequest true "Badge details to update"
// @Success 200 {object} models.Badge
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Badge not found"
// @Security ApiKeyAuth
// @Router /badges/{id} [put]
func UpdateBadgeHandler(badgeService *services.BadgeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid badge ID"})
			return
		}

		var req struct {
			Name        *string `json:"name"`
			Description *string `json:"description"`
			IconURL     *string `json:"icon_url"`
			Points      *int    `json:"points"`
			IsActive    *bool   `json:"is_active"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		badge, err := badgeService.UpdateBadge(
			id,
			req.Name,
			req.Description,
			req.IconURL,
			req.Points,
			req.IsActive,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update badge"})
			return
		}

		c.JSON(http.StatusOK, badge)
	}
}

// GetBadgeHandler retrieves a badge by ID
// @Summary Get a badge
// @Description Get a badge's details and achievement rules
// @Tags badges
// @Produce json
// @Param id path integer true "Badge ID"
// @Success 200 {object} BadgeResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Badge not found"
// @Security ApiKeyAuth
// @Router /badges/{id} [get]
func GetBadgeHandler(badgeService *services.BadgeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid badge ID"})
			return
		}

		badge, rules, err := badgeService.GetBadge(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Badge not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"badge": badge,
			"rules": rules,
		})
	}
}

// ListBadgesHandler retrieves all badges
// @Summary List all badges
// @Description Get a list of all badges (admins see all, users see only active)
// @Tags badges
// @Produce json
// @Success 200 {array} models.Badge
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Security ApiKeyAuth
// @Router /badges [get]
func ListBadgesHandler(badgeService *services.BadgeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only admins can see inactive badges
		includeInactive := false
		roles := c.GetStringSlice("roles")
		for _, role := range roles {
			if role == "admin" {
				includeInactive = true
				break
			}
		}

		badges, err := badgeService.ListBadges(includeInactive)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch badges"})
			return
		}

		c.JSON(http.StatusOK, badges)
	}
}

// AwardBadgeHandler awards a badge to a user
// @Summary Award badge to user
// @Description Award a badge to a specific user
// @Tags badges
// @Produce json
// @Param id path integer true "Badge ID"
// @Param user_id path integer true "User ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Badge or user not found"
// @Security ApiKeyAuth
// @Router /badges/{id}/award/{user_id} [post]
func AwardBadgeHandler(badgeService *services.BadgeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		badgeID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid badge ID"})
			return
		}

		userID, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		awardedBy := c.GetInt("userID")
		if err := badgeService.AwardBadge(userID, badgeID, &awardedBy); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Badge awarded successfully"})
	}
}

// GetUserBadgesHandler retrieves all badges awarded to a user
// @Summary Get user's badges
// @Description Get all badges that have been awarded to the authenticated user
// @Tags badges
// @Produce json
// @Param user_id path integer true "User ID"
// @Success 200 {array} models.UserBadge
// @Failure 400 {object} ErrorResponse "Invalid user ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden - Can only view your own badges"
// @Failure 404 {object} ErrorResponse "User not found"
// @Security ApiKeyAuth
// @Router /user/{id}/badges [get]
func GetUserBadgesHandler(badgeService *services.BadgeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Ensure user can only view their own badges
		requestingUserID := c.GetInt("userID")
		if requestingUserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own badges"})
			return
		}

		badges, err := badgeService.GetUserBadges(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user badges"})
			return
		}

		c.JSON(http.StatusOK, badges)
	}
}

// GetBadgeHoldersHandler retrieves all users who have been awarded a specific badge
// @Summary Get badge holders
// @Description Get all users who have been awarded a specific badge
// @Tags badges
// @Produce json
// @Param id path integer true "Badge ID"
// @Success 200 {array} models.UserBadge
// @Failure 400 {object} ErrorResponse "Invalid badge ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Badge not found"
// @Security ApiKeyAuth
// @Router /badges/{id}/holders [get]
func GetBadgeHoldersHandler(badgeService *services.BadgeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		badgeID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid badge ID"})
			return
		}

		// Check if badge exists
		badge, _, err := badgeService.GetBadge(badgeID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Badge not found"})
			return
		}

		holders, err := badgeService.GetBadgeHolders(badgeID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch badge holders"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"badge":   badge,
			"holders": holders,
		})
	}
}
