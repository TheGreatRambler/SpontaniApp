package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"googlemaps.github.io/maps"
	"honnef.co/go/spew"
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

	s3Client = s3.New(session.Must(session.NewSession()))

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
			description, lat, lng, uploaded, start, stop,
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
			time.Now(),
			time.Unix(request.Start, 0),
			time.Unix(request.Stop, 0),
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

		err = dbConn.QueryRow(context.Background(), `
			INSERT INTO img (task_id, uploaded, caption)
			VALUES ($1, $2, $3)
			RETURNING id
		`,
			taskId,
			time.Now(),
			caption,
		).Scan(&img_id)

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Failed to insert image: %v", err),
			}, nil
		}

		spew.Dump(request)

		// Upload image to S3
		_, err = s3Client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String("spontaniapp-imgs"),
			Key:    aws.String(fmt.Sprintf("%d", img_id)),
			Body:   bytes.NewReader([]byte(request.Body)),
		})

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Failed to upload image: %v", err),
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       fmt.Sprintf("{\"id\":%d}", img_id),
		}, nil

	case "update_image":
		img_id_str, exists := request.QueryStringParameters["id"]
		if !exists {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Missing required parameter: id",
			}, nil
		}

		task_id_str, exists := request.QueryStringParameters["task_id"]
		if !exists {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Missing required parameter: body",
			}, nil
		}

		img_id, img_id_err := strconv.Atoi(img_id_str)
		task_id, task_id_err := strconv.Atoi(task_id_str)
		if img_id_err != nil || task_id_err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Invalid required parameters: id, task_id",
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
			}, nil
		}

		_, err := dbConn.Exec(context.Background(), `
			UPDATE img SET task_id = $1 WHERE id = $2
		`, task_id, img_id)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       fmt.Sprintf("Database error: %v", err),
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil

	case "like":

		task_id_str, exists := request.QueryStringParameters["task_id"]

		if !exists {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Missing required parameter: task_id",
			}, nil
		}

		task_id, task_id_err := strconv.Atoi(task_id_str)
		if task_id_err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Invalid required parameter: task_id",
			}, nil
		}

		var likes int

		err := dbConn.QueryRow(context.Background(), `
				UPDATE task
				SET likes = likes + 1
				WHERE id = $1
				RETURNING likes
			`,
			task_id,
		).Scan(&likes)

		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Database error: %v", err),
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       fmt.Sprintf("{\"likes\":%d}", likes),
		}, nil

	case "get_presigned_url":
		img_id_str, exists := request.QueryStringParameters["id"]
		if !exists {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Missing required parameter: id",
			}, nil
		}

		img_id, img_id_err := strconv.Atoi(img_id_str)
		if img_id_err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Invalid required parameter: id",
			}, nil
		}

		req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String("spontaniapp-imgs"),
			Key:    aws.String(fmt.Sprintf("%d", img_id)),
		})

		url, err := req.Presign(15 * time.Minute)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Failed to generate presigned URL: %v", err),
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       fmt.Sprintf("{\"url\":\"%s\"}", url),
		}, nil

	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Incorrect parameter value: request_type",
		}, nil
	}

}

func corsHandlerWrapper(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response, err := handler(request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Internal server error: %v", err),
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}, nil
	}

	response.Headers = map[string]string{}
	response.Headers["Access-Control-Allow-Origin"] = "*"
	response.Headers["Access-Control-Allow-Headers"] = "*"
	response.Headers["Access-Control-Allow-Methods"] = "*"
	response.Headers["Access-Control-Allow-Credentials"] = "true"

	return response, nil
}

func main() {
	lambda.Start(corsHandlerWrapper)
}
