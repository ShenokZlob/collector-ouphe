package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID               string              `bson:"-" json:"id"`
	ObjectID         bson.ObjectID       `bson:"_id" json:"-"`
	TelegramID       int64               `bson:"telegram_id" json:"telegram_id"`
	Name             string              `bson:"name" json:"name"`
	TelegramNickname string              `bson:"telegram_nickname,omitempty" json:"telegram_nickname,omitempty"`
	Collections      []UserCollectionRef `bson:"collections,omitempty" json:"collections,omitempty"`
	CreatedAt        time.Time           `bson:"-" json:"created_at"`
	UpdatedAt        time.Time           `bson:"-" json:"updated_at"`
}

type UserCollectionRef struct {
	ObjectID bson.ObjectID `bson:"_id" json:"-"`
	ID       string        `bson:"-" json:"id"`
	Name     string        `bson:"name" json:"name"`
}

func (u *User) PrepareForResponse() {
	fmt.Printf("Type of ObjectID: %T, value: %#v\n", u.ObjectID, u.ObjectID)

	u.ID = u.ObjectID.Hex()

	// for i := range u.Collections {
	// 	fmt.Printf("Collection #%d ID: %T -> %#v\n", i, u.Collections[i].ObjectID, u.Collections[i].ObjectID)
	// 	u.Collections[i].ID = u.Collections[i].ObjectID.Hex()
	// }
}
