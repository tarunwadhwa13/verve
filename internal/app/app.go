package app

import (
	"database/sql"

	"verve/internal/api"
	"verve/internal/services"

	"github.com/gin-gonic/gin"
)

type App struct {
	db              *sql.DB
	router          *gin.Engine
	userService     *services.UserService
	walletService   *services.WalletService
	transferService *services.TransferService
	badgeService    *services.BadgeService
}

func NewApp(db *sql.DB, router *gin.Engine, userService *services.UserService, walletService *services.WalletService, transferService *services.TransferService, badgeService *services.BadgeService) *App {
	return &App{
		db:              db,
		router:          router,
		userService:     userService,
		walletService:   walletService,
		transferService: transferService,
		badgeService:    badgeService,
	}
}

func (a *App) SetupRoutes() {
	api.RegisterRoutes(a.router, a.db)
	api.RegisterUserRoutes(a.router, a.userService)
	api.RegisterWalletRoutes(a.router, a.walletService)
	api.RegisterTransferRoutes(a.router, a.transferService)
	api.RegisterBadgeRoutes(a.router, a.badgeService)
}

func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}
