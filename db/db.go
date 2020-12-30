package db

import (
	"context"
	"fmt"
	"time"

	"github.com/brxie/eluborzyca-backend/config"
	"github.com/brxie/eluborzyca-backend/utils/ilog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// DB is a database connection object
var DB *mongo.Database

func Connect() error {
	dbAddr := config.Viper.GetString("DB_ADDR")
	dbPort := config.Viper.GetInt("DB_PORT")
	dbName := config.Viper.GetString("DB_NAME")

	uri := fmt.Sprintf("mongodb://%s:%d", dbAddr, dbPort)
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// sessions collection
	collection := DB.Collection("sessions")

	tokenTTL := config.Viper.GetInt32("SESSION_TOKEN_TTL")
	model := mongo.IndexModel{
		Keys: bson.M{
			"created": 1, // index in ascending order
		}, Options: &options.IndexOptions{ExpireAfterSeconds: &tokenTTL},
	}
	_, err := collection.Indexes().CreateOne(ctx, model)
	if err != nil {
		return err
	}

	// users collection
	collection = DB.Collection("users")
	uniqie := true
	model = mongo.IndexModel{
		Keys: bson.M{
			"email": 1, // index in ascending order
		}, Options: &options.IndexOptions{Unique: &uniqie},
	}
	_, err = collection.Indexes().CreateOne(ctx, model)
	if err != nil {
		return err
	}

	// units collection
	collection = DB.Collection("units")
	uniqie = true
	model = mongo.IndexModel{
		Keys: bson.M{
			"name": 1, // index in ascending order
		}, Options: &options.IndexOptions{Unique: &uniqie},
	}
	_, err = collection.Indexes().CreateOne(ctx, model)
	if err != nil {
		return err
	}

	// villages collection
	collection = DB.Collection("villages")
	uniqie = true
	model = mongo.IndexModel{
		Keys: bson.M{
			"name": 1, // index in ascending order
		}, Options: &options.IndexOptions{Unique: &uniqie},
	}
	_, err = collection.Indexes().CreateOne(ctx, model)
	if err != nil {
		return err
	}

	// categories collection
	collection = DB.Collection("categories")
	uniqie = true
	model = mongo.IndexModel{
		Keys: bson.M{
			"name": 1, // index in ascending order
		}, Options: &options.IndexOptions{Unique: &uniqie},
	}
	_, err = collection.Indexes().CreateOne(ctx, model)
	if err != nil {
		return err
	}

	// verify_tokens collection
	collection = DB.Collection("verify_tokens")

	tokenTTL = 3600 * 24 // 24h
	model = mongo.IndexModel{
		Keys: bson.M{
			"created": 1, // index in ascending order
		}, Options: &options.IndexOptions{ExpireAfterSeconds: &tokenTTL},
	}
	_, err = collection.Indexes().CreateOne(ctx, model)
	if err != nil {
		return err
	}

	return nil

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
