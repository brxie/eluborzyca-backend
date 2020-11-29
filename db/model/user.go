package model

import (
	"context"
	"time"

	"github.com/brxie/ebazarek-backend/db"
	"go.mongodb.org/mongo-driver/bson"
)

// User represents user user
type User struct {
	Email    string    `bson:"email,omitempty"`
	Password string    `bson:"password,omitempty"`
	Created  time.Time `bson:"created,omitempty"`
}

const UsersCollectionName = "users"

func GetUser(query *User) (*User, error) {
	var (
		err  error
		user User
		doc  *bson.M
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.DB.Collection(UsersCollectionName)
	if doc, err = toBSON(query); err != nil {
		return nil, err
	}
	if err := collection.FindOne(ctx, doc).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func InsertUser(user *User) error {
	var (
		err error
		doc *bson.M
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.DB.Collection(UsersCollectionName)
	if doc, err = toBSON(user); err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, doc)
	return err
}
