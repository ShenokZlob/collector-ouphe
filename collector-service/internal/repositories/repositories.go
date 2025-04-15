package repositories

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	client *mongo.Client
}

const (
	database         = "collector_ouphe_db"
	users_collection = "users"
)

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (r Repository) CreateUser(user *models.User) (*models.User, *models.ResponseErr) {
	collection := r.client.Database(database).Collection(users_collection)
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, &models.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if id, ok := result.InsertedID.(bson.ObjectID); ok {
		user.ObjectID = id
	}

	return user, nil
}

func (r Repository) FindUserByTelegramID(telegramId int64) (*models.User, *models.ResponseErr) {
	collection := r.client.Database(database).Collection(users_collection)
	filter := bson.D{{Key: "telegram_id", Value: telegramId}}

	var user models.User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &models.ResponseErr{
				Status:  http.StatusNotFound,
				Message: "User not found",
			}
		}
		fmt.Println(user)
		return nil, &models.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Find user error: %v", err),
		}
	}
	return &user, nil
}
