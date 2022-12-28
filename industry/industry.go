package industry

import (
	"github.com/9d4/tracking-pi/person"
	"github.com/9d4/tracking-pi/place"
)

type Industry struct {
	Name     string        `bson:"name" json:"name"`
	Code     string        `bson:"code" json:"code"`
	Places   []place.Place `bson:"places" json:"places"`
	Advisers []Adviser     `bson:"advisers" json:"advisers"`
}

type Adviser struct {
	*person.Person
}
