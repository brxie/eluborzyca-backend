package model

import (
	"context"
	"time"

	"github.com/brxie/eluborzyca-backend/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string             `bson:"name,omitempty" json:"name,omitempty"`
	Icon string             `bson:"icon,omitempty" json:"icon,omitempty"`
}

const CategoriesCollectionName = "categories"

func GetCategory(query *Category) (*Category, error) {
	var (
		err      error
		category Category
		doc      *bson.M
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.DB.Collection(CategoriesCollectionName)
	if doc, err = toBSON(query); err != nil {
		return nil, err
	}
	if err := collection.FindOne(ctx, doc).Decode(&category); err != nil {
		return nil, err
	}
	return &category, nil
}
