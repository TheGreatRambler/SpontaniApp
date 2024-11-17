package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
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

func findClosestWaypoint(lat, lng float64) (string, string, error) {
	// Define the request for Nearby Search
	req := &maps.NearbySearchRequest{
		Location: &maps.LatLng{
			Lat: lat,
			Lng: lng,
		},
		Radius: 5000, // Search within 5km radius
	}

	// Execute the Nearby Search request
	resp, err := mapsClient.NearbySearch(context.Background(), req)
	if err != nil {
		return "", "", fmt.Errorf("failed to perform nearby search: %w", err)
	}

	// Check if there are any results
	if len(resp.Results) == 0 {
		return "", "", fmt.Errorf("no waypoints found near the specified location")
	}

	// Extract the name and address of the closest waypoint
	closest := resp.Results[0]

	return closest.Name, closest.FormattedAddress, nil
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

type TaskPost struct {
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Lat          float64 `json:"lat"`
	Lng          float64 `json:"lng"`
	Start        int64   `json:"start"`
	Stop         int64   `json:"stop"`
	InitialImgId int     `json:"initial_img_id"`
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
	case "create_task":
		var task_id int

		body := request.Body

		var request TaskPost
		err := json.Unmarshal([]byte(body), &request)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Invalid JSON body",
			}, nil
		}

		// Geolocate name and address
		location_name, location_address, err := findClosestWaypoint(request.Lat, request.Lng)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       fmt.Sprintf("Could not find closest endpoint: %v", err),
			}, nil
		}

		err = dbConn.QueryRow(context.Background(), `
			INSERT INTO task (title, location_name, location_address,
			description, lat, lng,uploaded, start, stop,
			initial_img_id, likes)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 0)
			RETURNING id
		`,
			request.Title,
			location_name,
			location_address,
			request.Description,
			request.Lat,
			request.Lng,
			time.Now().Unix(),
			request.Start,
			request.Stop,
			request.InitialImgId,
		).Scan(&task_id)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Database error: %v", err),
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       fmt.Sprintf("\"id\":%d}", task_id),
		}, nil
	case "upload_image":
		if request.Body == "" {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Missing required parameter: body",
			}, nil
		}
		if !request.IsBase64Encoded {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Request must be base64 encoded",
			}, nil
		}

		taskId, exists := request.QueryStringParameters["task_id"]
		if !exists {
			taskId = ""
		}
		caption, exists := request.QueryStringParameters["caption"]
		if !exists {
			caption = ""
		}

		var img_id int

		var err error

		if taskId == "" {
			err = dbConn.QueryRow(context.Background(), `
				INSERT INTO img (uploaded, caption)
				VALUES ($1, $2)
				RETURNING id
			`,
				time.Now().Unix(),
				caption,
			).Scan(&img_id)
		} else {
			err = dbConn.QueryRow(context.Background(), `
				INSERT INTO img (task_id, uploaded, caption)
				VALUES ($1, $2, $3)
				RETURNING id
			`,
				taskId,
				time.Now().Unix(),
				caption,
			).Scan(&img_id)
		}

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Failed to insert image: %v", err),
			}, nil
		}

		// Upload image to S3
		_, err = s3Client.PutObject(&s3.PutObjectInput{
			Bucket:      aws.String("spontaniapp-imgs"),
			Key:         aws.String(fmt.Sprintf("%d", img_id)),
			Body:        bytes.NewReader([]byte(request.Body)),
			ContentType: aws.String("image/jpeg"),
		})

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Failed to upload image: %v", err),
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string("{\"id\": " + fmt.Sprint(img_id) + "}"),
		}, nil

	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Incorrect parameter value: request_type",
		}, nil
	}

}

func main() {
	lambda.Start(handler)
}
