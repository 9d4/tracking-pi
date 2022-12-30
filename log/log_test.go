package log

import (
	"github.com/9d4/tracking-pi/place"
	"testing"
)

func TestPlaceAccuracy_Calculate(t *testing.T) {
	type fields struct {
		Place     *place.Place
		ToCompare *place.Coordinate
		Distance  float64
		InRange   *bool
	}
	tests := []struct {
		name    string
		fields  fields
		inRange bool
	}{
		{
			name: "-6.825262531887634, 111.10228693260008  vs  -6.825157809602922, 111.1016882218862",
			fields: fields{
				Place: &place.Place{
					Name:       "Loc 1",
					Coordinate: &place.Coordinate{111.10228693260008, -6.825262531887634},
					Wide:       70, // max allowed distance in metres
				},
				ToCompare: &place.Coordinate{111.1016882218862, -6.825157809602922},
			},
			inRange: true,
		},
		{
			name: "-6.825262531887634, 111.10228693260008  vs  -6.825157809602922, 111.1016882218862",
			fields: fields{
				Place: &place.Place{
					Name:       "Loc 1",
					Coordinate: &place.Coordinate{111.10228693260008, -6.825262531887634},
					Wide:       70, // max allowed distance in metres
				},
				ToCompare: &place.Coordinate{111.10168822, -6.82515780}, // usually we got .8f from browser
			},
			inRange: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PlaceAccuracy{
				Place:     tt.fields.Place,
				ToCompare: tt.fields.ToCompare,
			}
			p.Calculate()

			if tt.inRange != *p.InRange {
				t.Errorf("%s should be in range. Max allowed: %f, Got distance: %f", tt.name, tt.fields.Place.Wide, p.Distance)
			}
		})
	}
}
