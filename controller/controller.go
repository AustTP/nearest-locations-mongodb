package controller

import (
	"../models"
	"context"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	
	"googlemaps.github.io/maps"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const kmtomiles = float64(0.621371192)
const earthRadius = float64(6371)

func Location(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	
	youraddress := r.Form.Get("address")
	yourlongitude := r.Form.Get("longitude")
	yourlatitude := r.Form.Get("latitude")

	session, err := mgo.Dial("localhost")
  
	if err != nil {
		log.Fatal("could not connect to db: ", err)
		panic(err)
	}
	
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()
	
	var long float64
	var lat float64
	
	if youraddress != "" {
		var address = Google(youraddress)
		var numbers = strings.Split(address, ",")
		lo, la := numbers[1], numbers[0]
		long, err = strconv.ParseFloat(lo, 64)
		lat, err = strconv.ParseFloat(la, 64)
	} else {	
		long, err = strconv.ParseFloat(yourlongitude, 64)
		lat, err = strconv.ParseFloat(yourlatitude, 64)
	}

	scope := 3000

	var results []models.ShopLocation 

	c := session.DB("test").C("shops")
	err = c.Find(bson.M{
		"location": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
				"$maxDistance": scope,
			},
		},
	}).All(&results)

	if err != nil {
		panic(err)
	}
	
	digits := []models.Output{}
	var digit models.Output
	
	for _, widget := range results {
		digit.Name = widget.Name
	
		digit.Longitude = widget.Location.Coordinates[0]
		digit.Latitude = widget.Location.Coordinates[1]
		
		digit.Hours = widget.Hours
		
		var distance = Haversine(long, lat, digit.Longitude, digit.Latitude)
		digit.Distance = distance * kmtomiles
		digit.Address = widget.Address
		digit.City = widget.City
		digit.Phone = widget.Phone
		digit.Present = int(time.Weekday(5) - 1)
		
		digits = append(digits, digit)
	}
	
	MyPageVariables := models.PageVariables{
		Response: digits,
		Latitude: lat,
		Longitude: long,
		Address: youraddress,
    }

    t, err := template.ParseFiles("templates/form.gohtml")
	
    if err != nil { // if there is an error
		log.Print("template parsing error: ", err)
    }

    err = t.Execute(w, MyPageVariables)
	
    if err != nil { // if there is an error
		log.Print("template executing error: ", err)
    }
}

func Haversine(lonFrom float64, latFrom float64, lonTo float64, latTo float64) (distance float64) {
	var deltaLat = (latTo - latFrom) * (math.Pi / 180)
	var deltaLon = (lonTo - lonFrom) * (math.Pi / 180)
	
	var a = math.Sin(deltaLat / 2) * math.Sin(deltaLat / 2) + 
		math.Cos(latFrom * (math.Pi / 180)) * math.Cos(latTo * (math.Pi / 180)) *
		math.Sin(deltaLon / 2) * math.Sin(deltaLon / 2)
	var c = 2 * math.Atan2(math.Sqrt(a),math.Sqrt(1-a))
	
	distance = earthRadius * c
	
	return
}

func Google(location string) string {
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyBNgo0by1kOgD5m_ThExSCT2ize0Izy0nY"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	
	gr := &maps.GeocodingRequest{
		Address: location,
	}
	gresp, err := c.Geocode(context.Background(), gr)

	coords := gresp[0].Geometry.Location.String()

	return fmt.Sprintf("%s", coords)
}