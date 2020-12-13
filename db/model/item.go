package model

import (
	"context"
	"time"

	"github.com/brxie/ebazarek-backend/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name          string             `bson:"name,omitempty" json:"name,omitempty"`
	Owner         string             `bson:"owner,omitempty" json:"owner,omitempty"`
	Price         uint64             `bson:"price,omitempty" json:"price,omitempty"`
	Unit          string             `bson:"unit,omitempty" json:"unit,omitempty"`
	Availability  int                `bson:"availability,omitempty" json:"availability,omitempty"`
	FirstLastName string             `bson:"firstLastName,omitempty" json:"firstLastName,omitempty"`
	Village       string             `bson:"village,omitempty" json:"village,omitempty"`
	HomeNumber    string             `bson:"homeNumber,omitempty" json:"homeNumber,omitempty"`
	Phone         string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Category      string             `bson:"category,omitempty" json:"category,omitempty"`
	Description   string             `bson:"description,omitempty" json:"description,omitempty"`
	Popular       bool               `bson:"popular,omitempty" json:"popular,omitempty"`
	Active        bool               `bson:"active,omitempty" json:"active,omitempty"`
	Images        []Image            `bson:"images,omitempty" json:"images,omitempty"`
	Created       time.Time          `bson:"created,omitempty" json:"created,omitempty"`
	Deleted       time.Time          `bson:"deleted,omitempty" json:"deleted,omitempty"`
}

type ItemUpdate struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name          string             `bson:"name,omitempty" json:"name,omitempty"`
	Price         uint64             `bson:"price,omitempty" json:"price,omitempty"`
	Unit          string             `bson:"unit,omitempty" json:"unit,omitempty"`
	Availability  int                `bson:"availability,omitempty" json:"availability,omitempty"`
	FirstLastName string             `bson:"firstLastName,omitempty" json:"firstLastName,omitempty"`
	Village       string             `bson:"village,omitempty" json:"village,omitempty"`
	HomeNumber    string             `bson:"homeNumber,omitempty" json:"homeNumber,omitempty"`
	Phone         string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Category      string             `bson:"category,omitempty" json:"category,omitempty"`
	Description   string             `bson:"description,omitempty" json:"description,omitempty"`
}

type ItemActivate struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Active bool
}

const ItemsCollectionName = "items"
const DeletedItemsCollectionName = "items_deleted"

func GetItem(query *Item) (*Item, error) {
	var (
		err  error
		item Item
		doc  *bson.M
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.DB.Collection(ItemsCollectionName)
	if doc, err = toBSON(query); err != nil {
		return nil, err
	}
	if err := collection.FindOne(ctx, doc).Decode(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

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

func UpdateItem(filter *Item, update *ItemUpdate) error {
	var (
		err                  error
		filterDoc, updateDoc *bson.M
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.DB.Collection(ItemsCollectionName)

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

func ActivateItem(filter *Item, update *ItemActivate) error {
	var (
		err                  error
		filterDoc, updateDoc *bson.M
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := db.DB.Collection(ItemsCollectionName)

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

func DeleteItem(query *Item) error {

	var (
		err  error
		item Item
		doc  *bson.M
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.DB.Collection(ItemsCollectionName)
	collectionDeleted := db.DB.Collection(DeletedItemsCollectionName)
	if doc, err = toBSON(query); err != nil {
		return err
	}
	if err := collection.FindOne(ctx, doc).Decode(&item); err != nil {
		return err
	}

	item.Deleted = time.Now()
	if _, err = collectionDeleted.InsertOne(ctx, item); err != nil {
		return err
	}

	if _, err := collection.DeleteOne(ctx, &doc); err != nil {
		return err
	}

	return nil
}
