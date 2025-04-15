package repositories

import (
	"context"
	"net/http"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	client *mongo.Client
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) CreateUser(user *models.User) *models.ResponseErr {
	collection := r.client.Database("collector_ouphe_db").Collection("users")
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return &models.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ObjectID = id
	}

	return nil
}

func (r *Repository) GetUserByTelegramID(telegramId int64) (*models.User, *models.ResponseErr) {
	collection := r.client.Database("collector_ouphe_db").Collection("users")
	filter := bson.D{
		{Key: "telegramId", Value: telegramId},
	}

	var user *models.User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, &models.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return user, nil
}
