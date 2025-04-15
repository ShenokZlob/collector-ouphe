package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Collection struct {
	ID        string        `bson:"-" json:"id"`
	ObjectID  bson.ObjectID `bson:"_id" json:"-"`
	UserID    string        `bson:"user_id" json:"user_id"`
	Name      string        `bson:"name" json:"name"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}
