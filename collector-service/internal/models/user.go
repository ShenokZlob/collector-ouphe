package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID          string               `bson:"-" json:"id"`
	ObjectID    bson.ObjectID        `bson:"_id,omitempty" json:"-"`
	TelegramID  int64                `bson:"telegram_id" json:"telegram_id"`
	FirstName   string               `bson:"first_name" json:"first_name"`
	LastName    string               `bson:"last_name,omitempty" json:"last_name,omitempty"`
	Username    string               `bson:"username,omitempty" json:"username,omitempty"`
	Collections []*UserCollectionRef `bson:"collections,omitempty" json:"collections,omitempty"`
	CreatedAt   time.Time            `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at,omitempty" json:"updated_at"`
}

type UserCollectionRef struct {
	ObjectID bson.ObjectID `bson:"_id,omitempty" json:"-"`
	ID       string        `bson:"-" json:"id"`
	Name     string        `bson:"name" json:"name"`
}

func (u *User) PrepareForResponse() {
	// fmt.Printf("Type of ObjectID: %T, value: %#v\n", u.ObjectID, u.ObjectID)
	u.ID = u.ObjectID.Hex()

	for i := range u.Collections {
		// fmt.Printf("Collection #%d ID: %T -> %#v\n", i, u.Collections[i].ObjectID, u.Collections[i].ObjectID)
		u.Collections[i].ID = u.Collections[i].ObjectID.Hex()
	}
}
