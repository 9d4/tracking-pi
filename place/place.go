package place

type Coordinate struct {
	Longitude float64 `bson:"longitude" json:"longitude"`
	Latitude  float64 `bson:"latitude" json:"latitude"`
}

type Place struct {
	Name        string `bson:"name" json:"name"`
	*Coordinate `bson:",inline" json:",inline"`

	// Wide represents the wide of the area
	// This is used to calculate whether x,y is in *Coordinate or not.
	Wide float64 `bson:"wide" json:"wide"`
}
