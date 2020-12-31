package model

import (
	"context"
	"time"

	"github.com/brxie/eluborzyca-backend/db"
	"go.mongodb.org/mongo-driver/bson"
)

func GetCategories(query *Category) ([]Category, error) {
	var (
		err error
		doc *bson.M
	)
	categories := make([]Category, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.DB.Collection(CategoriesCollectionName)
	if doc, err = toBSON(query); err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, doc)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var category Category
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
