package volunteer

import (
	"github.com/9d4/tracking-pi/industry"
	"github.com/9d4/tracking-pi/person"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Volunteer struct {
	ID primitive.ObjectID `bson:"_id" json:"id"`

	*person.Person `bson:",inline"`
	Code           string `bson:"code" json:"code"`
	IndustryCode   string `bson:"industry_code" json:"industry_code"`

	ModelPath string
	Industry  *industry.Industry `bson:"industry,omitempty" json:"industry,omitempty"`

	Photo string `json:"photo"`
}
