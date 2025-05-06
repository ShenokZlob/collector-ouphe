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
	log    logger.Logger
	server *http.Server
}

func InitServer(config *viper.Viper, log logger.Logger, db *mongo.Client) *App {
	host := config.GetString("server_http.host")

	// Init repository
	rep := repositories.NewRepository(db)

	// Init services
	servAuth := services.NewAuthService(rep, log)
	servCollections := services.NewCollectionsService(rep, log)
	servCards := services.NewCardsService(rep, log)

	// Init controllers
	ctrlAuth := controllers.NewAuthController(servAuth, log)
	ctrlCollections := controllers.NewCollectionsController(servCollections, log)
	ctrlCards := controllers.NewCardsController(servCards, log)

	router := gin.Default()

	// Middleware
	mid := middleware.NewJWTMiddleware(os.Getenv("JWT_SECRET"), log)
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
		log:    log,
		server: server,
	}
}

func (a *App) Run() {
	a.log.Info("Running server", logger.String("host", a.host))
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.log.Error("ListenAndServe", logger.Error(err))
	}
}

func (a *App) Stop(ctx context.Context) {
	a.log.Info("Stopping server", logger.String("host", a.host))
	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Error("Server Shutdown Failed", logger.Error(err))
	}
}
