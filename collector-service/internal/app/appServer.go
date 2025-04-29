package app

import (
	"context"
	"net/http"
	"os"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/controllers"
	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/middleware"
	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/repositories"
	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/services"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	host   string
	logger logger.Logger
	server *http.Server
}

func InitServer(config *viper.Viper, logger logger.Logger, db *mongo.Client) *App {
	host := config.GetString("server_http.host")

	// Init repository
	rep := repositories.NewRepository(db)

	// Init services
	servAuth := services.NewAuthService(rep, logger)
	servCollections := services.NewCollectionsService(rep, logger)
	servCards := services.NewCardsService(rep, logger)

	// Init controllers
	ctrlAuth := controllers.NewAuthController(servAuth, logger)
	ctrlCollections := controllers.NewCollectionsController(servCollections, logger)
	ctrlCards := controllers.NewCardsController(servCards, logger)

	router := gin.Default()

	// Middleware
	mid := middleware.NewJWTMiddleware(os.Getenv("JWT_SECRET"), logger)
	authMiddleware := mid.Authorization()

	// Public routes
	public := router.Group("/")
	{
		public.POST("/register", ctrlAuth.Register)
		public.GET("/user/telegram/:telegram_id", ctrlAuth.Who)
		public.POST("/login", ctrlAuth.Login)
	}

	// Protected routes
	authorized := router.Group("/", authMiddleware)
	{
		authorized.GET("/collections", ctrlCollections.AllUsersCollections)
		authorized.POST("/collections", ctrlCollections.CreateCollection)
		authorized.PATCH("/collections/:id", ctrlCollections.RenameCollection)
		authorized.DELETE("/collections/:id", ctrlCollections.DeleteCollection)

		authorized.GET("/collections/:id/cards", ctrlCards.ListCardsInCollection)
		authorized.POST("/collections/:id/cards", ctrlCards.AddCardToCollection)
		authorized.PATCH("/collections/:id/cards/:id", ctrlCards.SetCardCountInCollection)
		authorized.DELETE("/collections/:id/cards/:id", ctrlCards.DeleteCardFromCollection)
	}

	server := &http.Server{
		Addr:    host,
		Handler: router.Handler(),
	}

	return &App{
		host:   host,
		logger: logger,
		server: server,
	}
}

func (a *App) Run() {
	a.logger.Info("Running server", logger.Field{Key: "host", String: a.host})
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.logger.Error("ListenAndServe", logger.Field{Key: "error", String: err.Error()})
	}
}

func (a *App) Stop(ctx context.Context) {
	a.logger.Info("Stopping server", logger.Field{Key: "host", String: a.host})
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("Server Shutdown Failed", logger.Field{Key: "error", String: err.Error()})
	}
}
