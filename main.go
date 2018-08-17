package main

import (
	"./controller"
	"net/http"
	"html/template"
	"log"
)

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/query", controller.Location)
	
	log.Fatal(http.ListenAndServe(":8040", nil))
}

func Home(w http.ResponseWriter, r *http.Request) {
	var templates = template.Must(template.New("geolocate").ParseFiles("templates/form.gohtml"))

	err := templates.ExecuteTemplate(w, "form.gohtml", nil)

	if err != nil {
		 panic(err)
	}
}