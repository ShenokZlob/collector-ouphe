package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ObjectID         primitive.ObjectID  `bson:"_id,omitempty" json:"-"`
	ID               string              `bson:"-" json:"id"`
	TelegramID       int64               `bson:"telegram_id" json:"telegram_id"`
	Name             string              `bson:"name" json:"name"`
	TelegramNickname string              `bson:"telegram_nickname,omitempty" json:"telegram_nickname,omitempty"`
	Collections      []UserCollectionRef `bson:"collections,omitempty" json:"collections,omitempty"`
	CreatedAt        time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time           `bson:"updated_at" json:"updated_at"`
}

type UserCollectionRef struct {
	ObjectID primitive.ObjectID `bson:"_id" json:"-"`
	ID       string             `bson:"-" json:"id"`
	Name     string             `bson:"name" json:"name"`
}

func (u *User) PrepareForResponse() {
	u.ID = u.ObjectID.Hex()

	for i := range u.Collections {
		u.Collections[i].ID = u.Collections[i].ObjectID.Hex()
	}
}
