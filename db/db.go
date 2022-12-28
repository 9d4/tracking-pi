package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var defaultClient *mongo.Client
var defaultDB *mongo.Database

func Open(uri string) (*mongo.Client, error) {
	return mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
}

func SetClient(client *mongo.Client) {
	defaultClient = client
}

func Client() *mongo.Client {
	return defaultClient
}

func SetDB(dbname string) {
	defaultDB = defaultClient.Database(dbname)
}

func DB() *mongo.Database {
	return defaultDB
}
