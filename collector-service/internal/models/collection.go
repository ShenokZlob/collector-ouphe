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
	Cards     []*Card       `bson:"cards,omitempty" json:"cards,omitempty"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

type Card struct {
	ScryfallID string    `bson:"scryfall_id" json:"scryfall_id"`
	Name       string    `bson:"name" json:"name"`
	CardUrl    string    `bson:"card_url" json:"card_url"`
	Count      int       `bson:"count" json:"count"`
	AddedAt    time.Time `bson:"added_at" json:"added_at"`
}

func (c *Collection) PrepareForResponse() {
	c.ID = c.ObjectID.Hex()
}
