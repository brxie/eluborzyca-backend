package model

import (
	"context"
	"time"

	"github.com/brxie/ebazarek-backend/db"
	"go.mongodb.org/mongo-driver/bson"
)

type Item struct {
	Name          string    `bson:"name,omitempty" json:"name,omitempty"`
	Owner         string    `bson:"owner,omitempty" json:"owner,omitempty"`
	Price         uint64    `bson:"price,omitempty" json:"price,omitempty"`
	Unit          string    `bson:"unit,omitempty" json:"unit,omitempty"`
	Availability  int       `bson:"availability,omitempty" json:"availability,omitempty"`
	FirstLastName string    `bson:"firstLastName,omitempty" json:"firstLastName,omitempty"`
	Village       string    `bson:"village,omitempty" json:"village,omitempty"`
	HomeNumber    string    `bson:"homeNumber,omitempty" json:"homeNumber,omitempty"`
	Phone         string    `bson:"phone,omitempty" json:"phone,omitempty"`
	Category      string    `bson:"category,omitempty" json:"category,omitempty"`
	Description   string    `bson:"description,omitempty" json:"description,omitempty"`
	Popular       bool      `bson:"popular,omitempty" json:"popular,omitempty"`
	Active        bool      `bson:"active,omitempty" json:"active,omitempty"`
	Images        []Image   `bson:"images,omitempty" json:"images,omitempty"`
	Created       time.Time `bson:"created,omitempty" json:"created,omitempty"`
}

const ItemsCollectionName = "items"

func InsertItem(item *Item) error {
	var (
		err error
		doc *bson.M
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.DB.Collection(ItemsCollectionName)
	if doc, err = toBSON(item); err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, doc)
	return err
}
