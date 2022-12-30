package log

import (
	"context"
	"github.com/9d4/tracking-pi/industry"
	"github.com/9d4/tracking-pi/place"
	"github.com/9d4/tracking-pi/volunteer"
	"github.com/jftuga/geodist"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	ilog "log"
)

type Log struct {
	VolunteerCode string `bson:"volunteer_code" json:"volunteer_code"`

	// base64 encoded
	Photo string `bson:"photo" json:"photo"`

	// current user coordinate
	Coordinate *place.Coordinate   `bson:"coordinate" json:"coordinate"`
	Timestamp  *primitive.DateTime `bson:"timestamp" json:"timestamp"`
}

// PlaceAccuracy used to compare two coordinates
type PlaceAccuracy struct {
	*place.Place `bson:",inline" json:",inline"`

	// This is coordinate from volunteer's log
	ToCompare *place.Coordinate `bson:"to_compare" json:"to_compare"`

	// Result of comparison
	Distance float64 `bson:"radius" json:"radius"`
	InRange  *bool   `bson:"in_range" json:"in_range"`
}

func (p *PlaceAccuracy) Calculate() {
	if p.ToCompare == nil || p.Place == nil {
		return
	}

	center := geodist.Coord{
		Lat: p.Coordinate.Latitude,
		Lon: p.Coordinate.Longitude,
	}
	target := geodist.Coord{
		Lat: p.ToCompare.Latitude,
		Lon: p.ToCompare.Longitude,
	}

	_, km := geodist.HaversineDistance(center, target)
	p.Distance = km * 1000 // save in metres

	p.InRange = new(bool)
	*p.InRange = true

	if p.Distance > p.Wide {
		*p.InRange = false
	}
}

func NewPlaceAccuracy(center *place.Place, targetCoord *place.Coordinate) *PlaceAccuracy {
	return &PlaceAccuracy{
		Place:     center,
		ToCompare: targetCoord,
		InRange:   new(bool),
	}
}

// LogResult represents complete log, log that has been processed.
type LogResult struct {
	*Log `bson:",inline" json:",inline"`

	Volunteer *volunteer.Volunteer `bson:"volunteer" json:"volunteer"`
	Industry  *industry.Industry   `bson:"industry" json:"industry"`

	// This would be true, if the volunteer's face model match with Log.Photo
	FaceMatch *bool `bson:"face_match" json:"face_match"`

	Places []*PlaceAccuracy `bson:"places" json:"places"`
}

// Run this in new go routine
func ProcessLogResult(logID primitive.ObjectID) {
	store := GetStore()

	// find log
	var log Log
	filter := bson.M{"_id": logID}
	err := store.coll.FindOne(context.Background(), filter).Decode(&log)
	if err != nil {
		ilog.Println(err)
		return
	}

	var logRes LogResult
	logRes.Log = &log

	// find volunteer with industry
	var (
		vols []volunteer.Volunteer
		vol  volunteer.Volunteer
	)
	pipeline := []bson.M{
		{"$match": bson.M{"code": log.VolunteerCode}},
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

	cursor, err := volunteer.GetCollection().Aggregate(context.Background(), pipeline)
	if err != nil {
		ilog.Println(err)
		return
	}
	if err = cursor.All(context.TODO(), &vols); err != nil {
		ilog.Println(err)
		return
	}
	if len(vols) < 1 {
		ilog.Println("volunteer not found")
		return
	}

	vol = vols[0]
	logRes.Volunteer = &vol
	logRes.Industry = logRes.Volunteer.Industry

	for _, industryPlace := range logRes.Industry.Places {
		// copy to avoid reference in array
		ip := industryPlace

		pa := NewPlaceAccuracy(&ip, logRes.Coordinate)
		pa.Calculate()
		logRes.Places = append(logRes.Places, pa)
	}

	// TODO: face matching

	filterD := bson.D{{"_id", logID}}
	updateD := bson.D{
		{
			"$set",
			bson.D{
				{"volunteer", logRes.Volunteer},
				{"industry", logRes.Industry},
				{"places", logRes.Places},
			},
		},
	}

	_, err = store.coll.UpdateOne(context.Background(), filterD, updateD)
	if err != nil {
		ilog.Println(err)
		return
	}
}
