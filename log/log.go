package log

import (
	"github.com/9d4/tracking-pi/place"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Log struct {
	VolunteerCode string `bson:"volunteer_code" json:"volunteer_code"`

	// base64 encoded
	Photo string `bson:"photo" json:"photo"`

	*place.Coordinate `bson:",inline"`
	*primitive.Timestamp
}
