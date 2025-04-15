package servers

import (
	"context"
	"net/http"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/controllers"
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

	reps := repositories.NewRepository(db)
	servs := services.NewService(reps)

	ctrlAuth := controllers.NewAuthController(servs)
	ctrlCollections := controllers.NewCollectionsController(servs)
	ctrlCards := controllers.NewCardsController(servs)

	router := gin.Default()

	router.POST("/register", ctrlAuth.Register)
	router.GET("/user/telegram/:telegram_id", ctrlAuth.Who)

	router.GET("/collections", ctrlCollections.AllCollections)
	router.POST("/collections", ctrlCollections.CreateCollection)
	router.PATCH("/collections/:id", ctrlCollections.RenameCollection)
	router.DELETE("/collections/:id", ctrlCollections.DeleteCollection)

	router.GET("/collections/:id/cards", ctrlCards.AllCardsByCollection)
	router.POST("/collections/:id/cards", ctrlCards.AddCardToCollection)
	router.PATCH("/collections/:id/cards/:id", ctrlCards.SetCardCount)
	router.DELETE("/collections/:id/cards/:id", ctrlCards.DeleteCard)

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
