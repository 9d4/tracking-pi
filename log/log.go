package log

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/9d4/tracking-pi/industry"
	"github.com/9d4/tracking-pi/place"
	"github.com/9d4/tracking-pi/volunteer"
	"github.com/gofiber/fiber/v2"
	"github.com/jftuga/geodist"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	ilog "log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

	photoMatch := make(chan bool, 1)
	go func() {
		ProcessPhoto(vol.ID, log.Photo, photoMatch)
	}()

	for _, industryPlace := range logRes.Industry.Places {
		// copy to avoid reference in array
		ip := industryPlace

		pa := NewPlaceAccuracy(&ip, logRes.Coordinate)
		pa.Calculate()
		logRes.Places = append(logRes.Places, pa)
	}

	// TODO: face matching
	logRes.FaceMatch = new(bool)
	*logRes.FaceMatch = <-photoMatch

	filterD := bson.D{{"_id", logID}}
	updateD := bson.D{
		{
			"$set",
			bson.D{
				{"volunteer", logRes.Volunteer},
				{"industry", logRes.Industry},
				{"places", logRes.Places},
				{"face_match", logRes.FaceMatch},
			},
		},
	}

	_, err = store.coll.UpdateOne(context.Background(), filterD, updateD)
	if err != nil {
		ilog.Println(err)
		return
	}
}

const ModelDir = "data/models/"

func ProcessPhoto(volunteerID primitive.ObjectID, targetPhotoB64 string, match chan<- bool) {
	defer close(match)

	frHost := os.Getenv("FACE_RECOGNITION_URI")
	if frHost == "" {
		ilog.Println("FACE_RECOGNITION_URI empty")
		match <- false
		return
	}

	frVerifyUrl, err := url.Parse(frHost)
	if err != nil {
		ilog.Println(err)
		match <- false
		return
	}
	frVerifyUrl = frVerifyUrl.JoinPath("verify")

	basePhoto := filepath.Join(ModelDir, volunteerID.Hex())
	basePhotoB64, err := os.ReadFile(basePhoto)
	if err != nil {
		ilog.Println(err)
		match <- false
		return
	}

	jsonBody := fiber.Map{
		"model_name": "Facenet",
		"img": []fiber.Map{
			fiber.Map{
				"img1": string(basePhotoB64),
				"img2": targetPhotoB64,
			},
		},
	}
	jsonBodyBytes := bytes.Buffer{}
	if err = json.NewEncoder(&jsonBodyBytes).Encode(jsonBody); err != nil {
		ilog.Println(err)
		match <- false
		return
	}

	response, err := http.Post(frVerifyUrl.String(), fiber.MIMEApplicationJSON, &jsonBodyBytes)
	if err != nil {
		ilog.Println(err)
		match <- false
		return
	}

	//{
	//	"pair_1": {
	//		"detector_backend": "opencv",
	//		"distance": 0.16681972852604476,
	//		"model": "Facenet",
	//		"similarity_metric": "cosine",
	//		"threshold": 0.4,
	//		"verified": true
	//	},
	//	"seconds": 0.421642541885376,
	//	"trx_id": "b05dbbfa-be3d-4fed-9f09-791b129a04be"
	//}

	var resData map[string]interface{}
	if err = json.NewDecoder(response.Body).Decode(&resData); err != nil {
		ilog.Println(err)
		match <- false
		return
	}

	if response.StatusCode != 200 {
		match <- false
		return
	}

	pair1, ok := resData["pair_1"].(map[string]interface{})
	if !ok {
		match <- false
		return
	}

	verified, ok := pair1["verified"].(bool)
	if !ok {
		match <- false
		return
	}

	match <- verified
}
