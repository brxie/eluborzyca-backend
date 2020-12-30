package model

import (
	"context"
	"time"

	"github.com/brxie/eluborzyca-backend/db"
	"go.mongodb.org/mongo-driver/bson"
)

func GetItems(query *Item) ([]Item, error) {
	var (
		err   error
		items []Item
		doc   *bson.M
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.DB.Collection(ItemsCollectionName)
	if doc, err = toBSON(query); err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, doc)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var item Item
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil

}
