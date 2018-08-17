package models

import (
	"gopkg.in/mgo.v2/bson"
)

type PageVariables struct {
	PageTitle        string
	Response		[]Output
	Distance		[]float64
	Longitude	float64
	Latitude	float64
	Address			string
}

type ShopLocation struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"shopid"`
	Name     string        `bson:"name" json:"name"`
	Location   GeoJson       `bson:"location" json:"location"`
	Distance	float64		`bson:"distance" json:"distance"`
	Address		string		`bson:"address" json:"address"`
	City	string			`bson:"city" json:"city"`
	Phone 	string 			`bson:"phone" json:"phone"`
	Hours	[]string		`bson:"hours" json:"hours"`
}

type GeoJson struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}	

type Output struct {
	Name     string        `bson:"name" json:"name"`
	Longitude	float64		`bson:"longitude" json:"longitude"`
	Latitude	float64		`bson:"latitude" json:"latitude"`
	Distance	float64		`bson:"distance" json:"distance"`
	Address		string		`bson:"address" json:"address"`
	City		string		`bson:"city" json:"city"`
	Phone		string		`bson:"phone" json:"phone"`
	Hours		[]string	`bson:"hours" json:"hours"`
	Present		int			`bson:"present" json:"present"`
}

