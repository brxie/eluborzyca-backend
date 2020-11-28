package model

import (
	"context"
	"time"

	"github.com/brxie/ebazarek-backend/db"
	"go.mongodb.org/mongo-driver/bson"
)

// Session represents user session
type Session struct {
	Token   string    `bson:"token,omitempty"`
	User    string    `bson:"user,omitempty"`
	Created time.Time `bson:"created,omitempty"`
}

func GetSession(query *Session) (*Session, error) {
	var (
		err     error
		session Session
		doc     *bson.M
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	colection := db.DB.Collection("sessions")
	if doc, err = toBSON(query); err != nil {
		return nil, err
	}
	if err := colection.FindOne(ctx, doc).Decode(&session); err != nil {
		return nil, err
	}
	return &session, nil
}
