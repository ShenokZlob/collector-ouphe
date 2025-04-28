package app

import (
	"context"
	"net/http"
	"os"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/controllers"
	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/middleware"
	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/repositories"
	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type App struct {
	host   string
	logger *zap.Logger
	server *http.Server
}

func InitServer(config *viper.Viper, logger *zap.Logger, db *mongo.Client) *App {
	host := config.GetString("server_http.host")

	// Init repository
	rep := repositories.NewRepository(db)

	// Init services
	servAuth := services.NewAuthService(rep)
	servCollections := services.NewCollectionsService(rep)
	servCards := services.NewCardsService(rep)

	// Init controllers
	ctrlAuth := controllers.NewAuthController(servAuth)
	ctrlCollections := controllers.NewCollectionsController(servCollections)
	ctrlCards := controllers.NewCardsController(servCards)

	router := gin.Default()
	authMiddleware := middleware.JWTMiddleware(os.Getenv("JWT_SECRET"))

	// public routes
	public := router.Group("/")
	{
		public.POST("/register", ctrlAuth.Register)
		public.GET("/user/telegram/:telegram_id", ctrlAuth.Who)
		public.POST("/login", ctrlAuth.Login)
	}

	// protected routes
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
	a.logger.Info("Running server", zap.String("host", a.host))
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.logger.Fatal("Something wrong...", zap.Error(err))
	}
}

func (a *App) Stop(ctx context.Context) {
	a.logger.Info("Stopping server", zap.String("host", a.host))
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Fatal("Server Shutdown:", zap.Error(err))
	}
}
