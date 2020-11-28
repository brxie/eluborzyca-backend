package db

import (
	"context"
	"fmt"
	"time"

	"github.com/brxie/ebazarek-backend/config"
	"github.com/brxie/ebazarek-backend/utils/ilog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// DB is a database connection object
var DB *mongo.Database

func Connect() error {
	config := config.DBconfig()
	dbAddr := config["DB_ADDR"]
	dbPort := config["DB_PORT"]
	dbName := config["DB_NAME"]

	uri := fmt.Sprintf("mongodb://%s:%s", dbAddr, dbPort)
	ilog.Debug("Connecting database uri: " + uri)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	DB = client.Database(dbName)
	return nil
}

func initializeModel() error {
	// sessions collection
	collection := DB.Collection("sessions")
	ttl, err := config.SessionTTL()
	if err != nil {
		return err
	}
	tokenTTL := int32(ttl)
	mod := mongo.IndexModel{
		Keys: bson.M{
			"created": 1, // index in ascending order
		}, Options: &options.IndexOptions{ExpireAfterSeconds: &tokenTTL},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = collection.Indexes().CreateOne(ctx, mod)

	return err

}

func init() {
	if err := Connect(); err != nil {
		ilog.Panic("Can't connect to database: " + err.Error())
	}

	if err := initializeModel(); err != nil {
		ilog.Panic("Can't initialize database model: " + err.Error())
	}

	ilog.Debug("Sucesfully initialized database connection.")
}
