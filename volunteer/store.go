package volunteer

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var defaultStore *Store

type Store struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func (s *Store) Create(vol *Volunteer) (*mongo.InsertOneResult, error) {
	doc := bson.D{
		{"name", vol.Name},
		{"code", vol.Code},
		{"industry_code", vol.IndustryCode},
	}

	return s.coll.InsertOne(context.TODO(), doc)
}

func (s *Store) GetAll() ([]Volunteer, error) {
	var volunteers []Volunteer

	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         "industries",
			"localField":   "industry_code",
			"foreignField": "code",
			"as":           "industries",
		}},
		{"$unwind": "$industries"},
		{"$group": bson.M{
			"_id":           "$_id",
			"name":          bson.M{"$first": "$name"},
			"industry_code": bson.M{"$first": "$industry_code"},
			"industry":      bson.M{"$first": "$industries"},
		}},
		{
			"$limit": 1,
		},
	}

	opts := options.Aggregate().SetAllowDiskUse(true)

	cursor, err := s.coll.Aggregate(context.TODO(), pipeline, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &volunteers); err != nil {
		return nil, err
	}

	return volunteers, nil
}

func (s *Store) Delete(ids []string) (*mongo.DeleteResult, error) {
	objectIDs := []primitive.ObjectID{}

	for _, id := range ids {
		if oid, err := primitive.ObjectIDFromHex(id); err == nil {
			objectIDs = append(objectIDs, oid)
		}
	}

	filter := bson.M{
		"_id": bson.M{
			"$in": objectIDs,
		},
	}

	return s.coll.DeleteMany(context.TODO(), filter)
}

func NewStore(db *mongo.Database) *Store {
	return &Store{
		db:   db,
		coll: db.Collection("volunteers"),
	}
}

func SetStore(store *Store) {
	defaultStore = store
}

func GetStore() *Store {
	return defaultStore
}

func GetCollection() *mongo.Collection {
	return defaultStore.coll
}
