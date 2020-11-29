package model

import (
	"context"
	"time"

	"github.com/brxie/ebazarek-backend/db"
	"go.mongodb.org/mongo-driver/bson"
)

// Session represents user session
type Session struct {
	Token   string    `bson:"token,omitempty"   json:"token,omitempty"`
	Email   string    `bson:"email,omitempty"   json:"email,omitempty"`
	Created time.Time `bson:"created,omitempty" json:"created,omitempty"`
}

const SessionsCollectionName = "sessions"

func GetSession(query *Session) (*Session, error) {
	var (
		err     error
		session Session
		doc     *bson.M
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.DB.Collection(SessionsCollectionName)
	if doc, err = toBSON(query); err != nil {
		return nil, err
	}
	if err := collection.FindOne(ctx, doc).Decode(&session); err != nil {
		return nil, err
	}
	return &session, nil
}

func InsertSession(session *Session) error {
	var (
		err error
		doc *bson.M
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.DB.Collection(SessionsCollectionName)
	if doc, err = toBSON(session); err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, doc)
	return err
}

func DestroySession(session *Session) error {
	var (
		err error
		doc *bson.M
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.DB.Collection(SessionsCollectionName)
	if doc, err = toBSON(session); err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, doc)
	return err
}
