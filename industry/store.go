package industry

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

func (s *Store) Create(industry *Industry) (*mongo.InsertOneResult, error) {
	doc := bson.D{
		{"name", industry.Name},
		{"places", industry.Places},
		{"advisers", industry.Advisers},
	}

	return s.coll.InsertOne(context.TODO(), doc)
}

func (s *Store) GetAll() ([]Industry, error) {
	inds := []Industry{}

	cursor, err := s.coll.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &inds); err != nil {
		return nil, err
	}

	return inds, nil
}

func NewStore(db *mongo.Database) *Store {
	return &Store{
		db:   db,
		coll: db.Collection("industries"),
	}
}

func SetStore(store *Store) {
	defaultStore = store
}

func GetStore() *Store {
	return defaultStore
}
