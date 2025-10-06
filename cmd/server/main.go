package main

import (
	"log"
	"time"

	"verve/internal/repository/postgres"
	"verve/internal/services"

	_ "verve/docs" // Swagger docs
	"verve/internal/app"
	"verve/internal/config"
	"verve/internal/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//	@title		Verve API
//	@version	1.0

//	@host		localhost:8080
//	@BasePath	/api
//	@schemes	http https

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
func main() {
	// Load application configurations
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	database, err := db.InitDB(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Create a new Gin router
	router := gin.Default()

	// enable Cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID", "X-Client-Version", "X-Platform"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize repositories
	userRepo := postgres.NewPostgresUserRepository(database)
	walletRepo := postgres.NewPostgresWalletRepository(database)
	roleRepo := postgres.NewPostgresRoleRepository(database)
	txRepo := postgres.NewPostgresTransactionRepository(database)
	ledgerRepo := postgres.NewPostgresLedgerRepository(database)
	transferRepo := postgres.NewPostgresTransferRepository(database)
	badgeRepo := postgres.NewPostgresBadgeRepository(database)
	achievementRuleRepo := postgres.NewPostgresAchievementRuleRepository(database)
	userBadgeRepo := postgres.NewPostgresUserBadgeRepository(database)

	// Initialize services
	userService := services.NewUserService(userRepo, roleRepo)
	walletService := services.NewWalletService(walletRepo)
	transferService := services.NewTransferService(transferRepo, txRepo, ledgerRepo, userRepo)
	badgeService := services.NewBadgeService(badgeRepo, achievementRuleRepo, userBadgeRepo)

	// Create a new application instance
	application := app.NewApp(database, router, userService, walletService, transferService, badgeService)

	// Setup routes
	application.SetupRoutes()

	// Start the server
	log.Printf("Server starting on %s", cfg.Server.Address)
	if err := application.Run(cfg.Server.Address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
