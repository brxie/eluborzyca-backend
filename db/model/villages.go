package model

import (
	"context"
	"time"

	"github.com/brxie/eluborzyca-backend/db"
	"go.mongodb.org/mongo-driver/bson"
)

func GetVillages(query *Village) ([]Village, error) {
	var (
		err      error
		villages []Village
		doc      *bson.M
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.DB.Collection(VillagesCollectionName)
	if doc, err = toBSON(query); err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, doc)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var village Village
		if err := cursor.Decode(&village); err != nil {
			return nil, err
		}
		villages = append(villages, village)
	}

	return villages, nil
}
