package model

import (
	"context"
	"time"

	"github.com/brxie/ebazarek-backend/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Image struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Src             string             `bson:"src,omitempty" json:"src,omitempty"`
	Thumbnail       string             `bson:"thumbnail,omitempty" json:"thumbnail,omitempty"`
	ThumbnailWidth  int                `bson:"thumbnailWidth,omitempty" json:"thumbnailWidth,omitempty"`
	ThumbnailHeight int                `bson:"thumbnailHeight,omitempty" json:"thumbnailHeight,omitempty"`
}

const ImagesCollectionName = "images"

func GetImage(query *Image) (*Image, error) {
	var (
		err   error
		image Image
		doc   *bson.M
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := db.DB.Collection(ImagesCollectionName)
	if doc, err = toBSON(query); err != nil {
		return nil, err
	}
	if err := collection.FindOne(ctx, doc).Decode(&image); err != nil {
		return nil, err
	}
	return &image, nil
}
