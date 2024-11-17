package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"googlemaps.github.io/maps"
)

var mapsClient *maps.Client
var s3Client *s3.S3
var dbConn *pgxpool.Pool

type TaskRet struct {
	Id           int     `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Lat          float64 `json:"lat"`
	Lng          float64 `json:"lng"`
	Uploaded     int64   `json:"uploaded"`
	Start        int64   `json:"start"`
	Stop         int64   `json:"stop"`
	InitialImgId int     `json:"initial_img_id"`
	Likes        int     `json:"likes"`
}

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

	pgx_config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(fmt.Sprintf("Invalid databse URL: %v", os.Getenv("DATABASE_URL")))
	}
	pgx_config.MaxConns = 32
	pgx_config.ConnConfig.ConnectTimeout = 10 * time.Second
	dbConn, err = pgxpool.NewWithConfig(context.Background(), pgx_config)
	if err != nil {
		panic(fmt.Sprintf("Could not connect to database: %v", os.Getenv("DATABASE_URL")))
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request_type, exists := request.QueryStringParameters["request_type"]
	if !exists {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing required parameter: request_type",
		}, nil
	}

	switch request_type {
	case "get_google_maps_key":
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       os.Getenv("DATABASE_URL"),
		}, nil
	case "get_nearby_recent_tasks":
		lat_str, lat_exists := request.QueryStringParameters["lat"]
		lng_str, lng_exists := request.QueryStringParameters["lng"]
		if !lat_exists || !lng_exists {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Missing required parameters: lat, lng",
			}, nil
		}

		lat, lat_err := strconv.ParseFloat(lat_str, 64)
		lng, lng_err := strconv.ParseFloat(lng_str, 64)
		if lat_err != nil || lng_err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Invalid required parameters: lat, lng",
			}, nil
		}

		rows, err := dbConn.Query(context.Background(), `
			SELECT id, title, description, lat, lng, uploaded,
				start, stop, initial_img_id, likes,
				point($1, $2) <@>  (point(lat, lng)::point) as distance
				FROM task ORDER BY distance ASC
		`, lat, lng)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Database error: %v", err),
			}, nil
		}
		defer rows.Close()

		tasks := []TaskRet{}
		for rows.Next() {
			var id int
			var title string
			var description string
			var lat float64
			var lng float64
			var uploaded int64
			var start int64
			var stop int64
			var initial_img_id int
			var likes int
			var distance float64
			err := rows.Scan(&id, &title, &description, &lat, &lng, &uploaded, &start, &stop, &initial_img_id, &likes, &distance)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 500,
					Body:       fmt.Sprintf("Database error: %v", err),
				}, nil
			}
			tasks = append(tasks, TaskRet{
				Id:           id,
				Title:        title,
				Description:  description,
				Lat:          lat,
				Lng:          lng,
				Uploaded:     uploaded,
				Start:        start,
				Stop:         stop,
				InitialImgId: initial_img_id,
				Likes:        likes,
			})
		}

		tasks_json, err := json.Marshal(tasks)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("JSON marshalling error: %v", err),
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(tasks_json),
		}, nil
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Incorrect parameter value: request_type",
		}, nil
	}

	// // Perform Nearby Search
	// response, err := mapsClient.NearbySearch(context.Background(), &maps.NearbySearchRequest{
	// 	Location: &maps.LatLng{
	// 		Lat: lat,
	// 		Lng: lng,
	// 	},
	// 	Radius: uint(10), // Adjust radius as needed
	// })
	// if err != nil {
	// 	return Response{
	// 		StatusCode: 500,
	// 		Body:       nil,
	// 	}, fmt.Errorf("failed to perform Nearby Search: %v", err)
	// }

	// // Return results as JSON
	// return Response{
	// 	StatusCode: 200,
	// 	Body:       response.Results,
	// }, nil
}

func main() {
	lambda.Start(handler)
}
