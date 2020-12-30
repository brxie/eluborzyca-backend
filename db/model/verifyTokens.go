package model

import (
	"context"
	"time"

	"github.com/brxie/ebazarek-backend/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VerifyToken struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email   string             `bson:"email,omitempty" json:"email,omitempty"`
	Token   string             `bson:"token,omitempty" json:"token,omitempty"`
	Created time.Time          `bson:"created,omitempty" json:"created,omitempty"`
}

const VerifyTokensCollectionName = "verify_tokens"

func GetVerifyToken(query *VerifyToken) (*VerifyToken, error) {
	var (
		err         error
		verifyToken VerifyToken
		doc         *bson.M
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.DB.Collection(VerifyTokensCollectionName)
	if doc, err = toBSON(query); err != nil {
		return nil, err
	}
	if err := collection.FindOne(ctx, doc).Decode(&verifyToken); err != nil {
		return nil, err
	}
	return &verifyToken, nil
}

func InsertVerifyToken(verifyToken *VerifyToken) error {
	var (
		err error
		doc *bson.M
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.DB.Collection(VerifyTokensCollectionName)
	if doc, err = toBSON(verifyToken); err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, doc)
	return err
}

func UpdateVerifyToken(filter, update *VerifyToken) error {
	var (
		err                  error
		filterDoc, updateDoc *bson.M
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.DB.Collection(VerifyTokensCollectionName)

	if filterDoc, err = toBSON(filter); err != nil {
		return err
	}

	if updateDoc, err = toBSON(update); err != nil {
		return err
	}

	_, err = collection.UpdateOne(ctx, filterDoc, bson.M{
		"$set": updateDoc,
	})
	return err
}
