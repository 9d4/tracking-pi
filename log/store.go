package log

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var defaultStore *Store

type Store struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func (s *Store) Create(log *Log) (*mongo.InsertOneResult, error) {
	doc := bson.D{
		{"volunteer_code", log.VolunteerCode},
		{"photo", log.Photo},
		{"coordinate", log.Coordinate},
	}

	return s.coll.InsertOne(context.TODO(), doc)
}

func NewStore(db *mongo.Database) *Store {
	return &Store{
		db:   db,
		coll: db.Collection("logs"),
	}
}

func SetStore(store *Store) {
	defaultStore = store
}

func GetStore() *Store {
	return defaultStore
}
