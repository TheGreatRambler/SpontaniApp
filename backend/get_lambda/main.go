package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"googlemaps.github.io/maps"
)

// Response structure for returning the nearby places
type Response struct {
	StatusCode int                       `json:"statusCode"`
	Body       []maps.PlacesSearchResult `json:"body"`
}

var mapsClient *maps.Client
var s3Client *s3.S3

func init() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file, proceeding without it...")
	}

	// Initialize Google Maps client
	mapsClient, err = maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_KEY")))
	if err != nil {
		panic(fmt.Sprintf("Failed to create Google Maps client: %v", err))
	}

	s3Client = s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("S3_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	})))
}

func handler(request events.APIGatewayProxyRequest) (Response, error) {
	// Extract query parameters
	latStr, latSet := request.QueryStringParameters["lat"]
	lngStr, lngSet := request.QueryStringParameters["lng"]

	if !latSet || !lngSet {
		return Response{
			StatusCode: 400,
			Body:       nil,
		}, fmt.Errorf("location not provided")
	}

	lat, err1 := strconv.ParseFloat(latStr, 64)
	lng, err2 := strconv.ParseFloat(lngStr, 64)
	if err1 != nil || err2 != nil {
		return Response{
			StatusCode: 400,
			Body:       nil,
		}, fmt.Errorf("invalid location coordinates")
	}

	// Perform Nearby Search
	response, err := mapsClient.NearbySearch(context.Background(), &maps.NearbySearchRequest{
		Location: &maps.LatLng{
			Lat: lat,
			Lng: lng,
		},
		Radius: uint(10), // Adjust radius as needed
	})
	if err != nil {
		return Response{
			StatusCode: 500,
			Body:       nil,
		}, fmt.Errorf("failed to perform Nearby Search: %v", err)
	}

	// Return results as JSON
	return Response{
		StatusCode: 200,
		Body:       response.Results,
	}, nil
}

func main() {
	lambda.Start(handler)
}
