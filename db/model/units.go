package model

import (
	"context"
	"time"

	"github.com/brxie/eluborzyca-backend/db"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUnits(query *Unit) ([]Unit, error) {
	var (
		err   error
		units []Unit
		doc   *bson.M
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.DB.Collection(UnitsCollectionName)
	if doc, err = toBSON(query); err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, doc)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var unit Unit
		if err := cursor.Decode(&unit); err != nil {
			return nil, err
		}
		units = append(units, unit)
	}

	return units, nil
}
