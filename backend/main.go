package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"googlemaps.github.io/maps"
)

func main() {
	// Design: Return constructed HTML that solidjs embeds. SolidJS is capable of rendering its own components if neccesary
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Create a client
	maps_client, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_KEY")))
	if err != nil {
		panic(err)
	}

	const port = 8080

	r := mux.NewRouter()

	r.HandleFunc("/close_place", func(w http.ResponseWriter, r *http.Request) {
		response, err := maps_client.NearbySearch(context.Background(), &maps.NearbySearchRequest{
			Location: &maps.LatLng{
				Lat: 1,
				Lng: 2,
			},
			Radius: uint(10),
		})
		if err != nil {
			w.WriteHeader(500)
		}

		for _, result := range response.Results {
			fmt.Printf("Name: %s, Location: (%f, %f)\n",
				result.Name,
				result.Geometry.Location.Lat,
				result.Geometry.Location.Lng)
		}
	})

	fmt.Printf("Starting server at http://localhost:%v\n", port)

	go http.ListenAndServe(fmt.Sprintf(":%v", port), r)
}
