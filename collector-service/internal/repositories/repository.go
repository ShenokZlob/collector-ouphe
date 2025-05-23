package repositories

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	client *mongo.Client
}

const (
	database               = "collector_ouphe_db"
	users_collection       = "users"
	collections_collection = "collections"
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

	user.PrepareForResponse()
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

	user.PrepareForResponse()
	return &user, nil
}

// Search names in Users collection
func (r Repository) UsersCollections(userId string) ([]*models.UserCollectionRef, *models.ResponseErr) {
	objectID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, &models.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid user ID format",
		}
	}

	collection := r.client.Database(database).Collection(users_collection)
	filter := bson.D{{Key: "_id", Value: objectID}}

	var user models.User
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &models.ResponseErr{
				Status:  http.StatusNotFound,
				Message: "User not found",
			}
		}
		return nil, &models.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Find user error: %v", err),
		}
	}

	user.PrepareForResponse()
	return user.Collections, nil
}

func (r Repository) CreateCollection(collection *models.Collection) (*models.Collection, *models.ResponseErr) {
	collectionRef := r.client.Database(database).Collection(collections_collection)
	result, err := collectionRef.InsertOne(context.TODO(), collection)
	if err != nil {
		return nil, &models.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if id, ok := result.InsertedID.(bson.ObjectID); ok {
		collection.ObjectID = id
	}

	collection.PrepareForResponse()
	return collection, nil
}

func (r Repository) RenameCollection(collection *models.Collection) (*models.Collection, *models.ResponseErr) {
	collectionRef := r.client.Database(database).Collection(collections_collection)
	objectId, err := bson.ObjectIDFromHex(collection.ID)
	if err != nil {
		return nil, &models.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection ID format",
		}
	}

	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: collection.Name}}}}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated models.Collection
	err = collectionRef.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &models.ResponseErr{
				Status:  http.StatusNotFound,
				Message: "Collection not found",
			}
		}
		return nil, &models.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Error updating collection: %v", err),
		}
	}

	return &updated, nil
}

func (r Repository) DeleteCollection(collection *models.Collection) *models.ResponseErr {
	collectionRef := r.client.Database(database).Collection(collections_collection)
	objectId, err := bson.ObjectIDFromHex(collection.ID)
	if err != nil {
		return &models.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection ID format",
		}
	}

	filter := bson.D{{Key: "_id", Value: objectId}}
	_, err = collectionRef.DeleteOne(context.TODO(), filter)
	if err != nil {
		return &models.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Delete collection error: %v", err),
		}
	}

	return nil
}

// Search in Collections collection
func (r Repository) GetCollection(collectionId string) (*models.Collection, *models.ResponseErr) {
	objectId, err := bson.ObjectIDFromHex(collectionId)
	if err != nil {
		return nil, &models.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection ID format",
		}
	}
	collection := r.client.Database(database).Collection(collections_collection)
	filter := bson.D{{Key: "_id", Value: objectId}}

	var col models.Collection
	err = collection.FindOne(context.TODO(), filter).Decode(&col)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &models.ResponseErr{
				Status:  http.StatusNotFound,
				Message: "Collection not found",
			}
		}
		return nil, &models.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Find collection error: %v", err),
		}
	}

	col.PrepareForResponse()
	return &col, nil
}

func (r Repository) AddCardToCollection(collectionId string, card *models.Card) *models.ResponseErr {
	objectId, err := bson.ObjectIDFromHex(collectionId)
	if err != nil {
		return &models.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection ID format",
		}
	}

	collection := r.client.Database(database).Collection(collections_collection)

	// Try to update the card count first
	filter := bson.M{
		"_id":               objectId,
		"cards.scryfall_id": card.ScryfallID,
	}
	update := bson.M{
		"$inc": bson.M{"cards.$.count": card.Count},
		"$set": bson.M{"updated_at": time.Now()},
	}
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return &models.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Error updating card count: %v", err),
		}
	}

	if result.MatchedCount > 0 {
		return nil
	}

	// If the card doesn't exist, add it to the collection
	filter = bson.M{"_id": objectId}
	push := bson.M{
		"$push": bson.M{"cards": card},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	_, err = collection.UpdateOne(context.TODO(), filter, push)
	if err != nil {
		return &models.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Error adding new card: %v", err),
		}
	}

	return nil
}

func (r Repository) SetCardCountInCollection(collectionId string, card *models.Card) *models.ResponseErr {
	objectId, err := bson.ObjectIDFromHex(collectionId)
	if err != nil {
		return &models.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection ID format",
		}
	}

	collection := r.client.Database(database).Collection(collections_collection)
	filter := bson.D{
		{Key: "_id", Value: objectId},
		{Key: "cards.scryfall_id", Value: card.ScryfallID},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "cards.$.count", Value: card.Count},
			{Key: "updated_at", Value: time.Now()},
		}},
	}
	//opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	// var updatedColection models.Collection
	res := collection.FindOneAndUpdate(context.TODO(), filter, update)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return &models.ResponseErr{
				Status:  http.StatusNotFound,
				Message: "Collection not found",
			}
		}
		return &models.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Update collection error: %v", err),
		}
	}

	return nil
}

func (r Repository) DeleteCardFromCollection(collectionId string, card *models.Card) *models.ResponseErr {
	objectId, err := bson.ObjectIDFromHex(collectionId)
	if err != nil {
		return &models.ResponseErr{
			Status:  http.StatusBadRequest,
			Message: "Invalid collection ID format",
		}
	}

	collection := r.client.Database(database).Collection(collections_collection)
	filter := bson.D{
		{Key: "_id", Value: objectId},
		{Key: "cards.scryfall_id", Value: card.ScryfallID},
	}
	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "cards", Value: bson.D{
				{Key: "scryfall_id", Value: card.ScryfallID},
			}},
		}},
		{Key: "$set", Value: bson.D{
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	res := collection.FindOneAndUpdate(context.TODO(), filter, update)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return &models.ResponseErr{
				Status:  http.StatusNotFound,
				Message: "Collection not found",
			}
		}
		return &models.ResponseErr{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("Update collection error: %v", err),
		}
	}

	return nil
}
