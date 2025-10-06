package api

import (
	"net/http"
	"strconv"
	"verve/internal/api/middleware"
	"verve/internal/auth"
	"verve/internal/services"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes sets up user-related routes
// @Summary Register user routes
// @Description Register all user-related routes including authentication and user management
// @Tags users
func RegisterUserRoutes(router *gin.Engine, userService *services.UserService) {
	userRoutes := router.Group("/api/user")
	{
		userRoutes.POST("/register", middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"), CreateUserHandler(userService))
		userRoutes.GET("/connected", middleware.AuthMiddleware(), GetAllUsersHandler(userService))
		userRoutes.GET("/:id", middleware.AuthMiddleware(), GetUserHandler(userService))
		userRoutes.POST("/:id/pin", middleware.AuthMiddleware(), SetPinHandler(userService))
		userRoutes.PUT("/:id", middleware.AuthMiddleware(), UpdateUserHandler(userService))
	}
	router.POST("/api/auth/login", LoginHandler(userService))
}

// LoginHandler authenticates users
// @Summary User login
// @Description Authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Router /api/auth/login [post]

// UpdateUserHandler handles user profile updates
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Param id path integer true "User ID"
// @Param user body UpdateUserRequest true "User details to update"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden - Can only update own profile"
// @Security ApiKeyAuth
// @Router /api/user/{id} [put]
func UpdateUserHandler(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		if userID != id {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own profile"})
			return
		}

		var req struct {
			DisplayName            *string `json:"display_name"`
			ProfilePhotoURL        *string `json:"profile_photo_url"`
			PinRequiredForTransfer *bool   `json:"pin_required_for_transfer"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userService.UpdateUser(userID, req.DisplayName, req.ProfilePhotoURL, req.PinRequiredForTransfer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// CreateUserHandler handles user creation (admin only)
// @Summary Create new user
// @Description Create a new user with specified roles (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User details"
// @Success 201 {object} models.User
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Forbidden - Admin only"
// @Security ApiKeyAuth
// @Router /api/user/register [post]
func CreateUserHandler(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userService.CreateUser(req.Username, req.Password, req.Pin, req.Roles)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}
func SetPinHandler(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		var req struct {
			Pin string `json:"pin"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := userService.SetPin(userID, req.Pin); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set PIN"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "PIN set successfully"})
	}
}

// GetUserHandler retrieves user details
// @Summary Get user details
// @Description Get details of a specific user
// @Tags users
// @Produce json
// @Param id path integer true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse "Invalid user ID"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "User not found"
// @Security ApiKeyAuth
// @Router /api/user/{id} [get]
func GetUserHandler(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		user, err := userService.GetUserByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// GetAllUsersHandler retrieves all users
// @Summary List all users
// @Description Get a list of all users in the system
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "No users found"
// @Security ApiKeyAuth
// @Router /api/user/connected [get]
func GetAllUsersHandler(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := userService.GetAllUsers()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// LoginHandler authenticates users
// @Summary User login
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Router /api/auth/login [post]
func LoginHandler(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, roles, err := userService.Authenticate(req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		token, err := auth.GenerateJWT(userID, roles)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"username": req.Username,
			},
			"refresh_token": "",
		})
	}
}
