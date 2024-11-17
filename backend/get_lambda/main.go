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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"googlemaps.github.io/maps"
)

var mapsClient *maps.Client
var s3Client *s3.S3
var dbConn *pgxpool.Pool

type TaskRet struct {
	Id              int     `json:"id"`
	Title           string  `json:"title"`
	LocationName    string  `json:"location_name"`
	LocationAddress string  `json:"location_address"`
	Description     string  `json:"description"`
	Lat             float64 `json:"lat"`
	Lng             float64 `json:"lng"`
	Uploaded        int64   `json:"uploaded"`
	Start           int64   `json:"start"`
	Stop            int64   `json:"stop"`
	InitialImgId    int     `json:"initial_img_id"`
	Likes           int     `json:"likes"`
}

type ImgRet struct {
	Id       int    `json:"id"`
	TaskID   int    `json:"task_id"`
	Uploaded int64  `json:"uploaded"`
	Caption  string `json:"caption"`
	URL      string `json:"url"`
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

func getLatLngParameters(request events.APIGatewayProxyRequest) (float64, float64, *events.APIGatewayProxyResponse) {
	lat_str, lat_exists := request.QueryStringParameters["lat"]
	lng_str, lng_exists := request.QueryStringParameters["lng"]
	if !lat_exists || !lng_exists {
		return 0, 0, &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing required parameters: lat, lng",
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}
	}

	lat, lat_err := strconv.ParseFloat(lat_str, 64)
	lng, lng_err := strconv.ParseFloat(lng_str, 64)
	if lat_err != nil || lng_err != nil {
		return 0, 0, &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid required parameters: lat, lng",
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}
	}

	return lat, lng, nil
}

type RowScanner interface {
	Scan(dest ...interface{}) error
}

func parseTask(row RowScanner) (TaskRet, *events.APIGatewayProxyResponse) {
	var id int
	var title string
	var location_name string
	var location_address string
	var description string
	var lat float64
	var lng float64
	var uploaded time.Time
	var start time.Time
	var stop time.Time
	var initial_img_id int
	var likes int
	var distance float64
	err := row.Scan(&id, &title, &location_name, &location_address, &description, &lat, &lng, &uploaded, &start, &stop, &initial_img_id, &likes, &distance)
	if err != nil {
		return TaskRet{}, &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Database error: %v", err),
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}
	}

	return TaskRet{
		Id:              id,
		Title:           title,
		LocationName:    location_name,
		LocationAddress: location_address,
		Description:     description,
		Lat:             lat,
		Lng:             lng,
		Uploaded:        uploaded.Unix(),
		Start:           start.Unix(),
		Stop:            stop.Unix(),
		InitialImgId:    initial_img_id,
		Likes:           likes,
	}, nil
}

func buildTaskJSON(rows pgx.Rows) (string, *events.APIGatewayProxyResponse) {
	tasks := []TaskRet{}
	for rows.Next() {
		task, res := parseTask(rows)
		if res != nil {
			return "", res
		}

		tasks = append(tasks, task)
	}

	tasks_json, err := json.Marshal(tasks)
	if err != nil {
		return "", &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("JSON marshalling error: %v", err),
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}
	}

	return string(tasks_json), nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	request_type, exists := request.QueryStringParameters["request_type"]
	if !exists {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing required parameter: request_type",
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}, nil
	}

	switch request_type {
	case "get_google_maps_key":
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       os.Getenv("GOOGLE_MAPS_KEY"),
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}, nil
	case "get_nearby_recent_tasks":
		lat, lng, res := getLatLngParameters(request)
		if res != nil {
			return *res, nil
		}

		rows, err := dbConn.Query(context.Background(), `
			SELECT id, title, location_name, location_address,
				description, lat, lng, uploaded,
				start, stop, initial_img_id, likes,
				point($1, $2) <@>  (point(lat, lng)::point) as distance
				FROM task WHERE start > $3 AND stop < $3
				ORDER BY distance ASC
		`, lat, lng, time.Now())
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Database error: %v", err),
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
			}, nil
		}
		defer rows.Close()

		tasks_json, res := buildTaskJSON(rows)
		if res != nil {
			return *res, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       tasks_json,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	case "get_task":
		id, id_exists := request.QueryStringParameters["id"]
		if !id_exists {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Missing required parameters: id",
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
			}, nil
		}

		row := dbConn.QueryRow(context.Background(), `
			SELECT id, title, location_name, location_address,
				description, lat, lng, uploaded,
				start, stop, initial_img_id, likes
				FROM task WHERE id = $1
		`, id)

		task, res := parseTask(row)
		if res != nil {
			return *res, nil
		}

		task_json, err := json.Marshal(task)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("JSON marshalling error: %v", err),
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(task_json),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	case "get_recent_tasks":
		rows, err := dbConn.Query(context.Background(), `
			SELECT id, title, location_name, location_address,
				description, lat, lng, uploaded,
				start, stop, initial_img_id, likes
				FROM task WHERE start > $1 AND stop < $1 ORDER BY start ASC
		`, time.Now())
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Database error: %v", err),
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
			}, nil
		}
		defer rows.Close()

		tasks_json, res := buildTaskJSON(rows)
		if res != nil {
			return *res, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       tasks_json,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	case "get_popular_tasks":
		rows, err := dbConn.Query(context.Background(), `
			SELECT id, title, location_name, location_address,
				description, lat, lng, uploaded,
				start, stop, initial_img_id, likes
				FROM task WHERE start > $1 AND stop < $1 ORDER BY likes DESC
		`, time.Now())
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Database error: %v", err),
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
			}, nil
		}
		defer rows.Close()

		tasks_json, res := buildTaskJSON(rows)
		if res != nil {
			return *res, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       tasks_json,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	case "get_active_tasks":
		rows, err := dbConn.Query(context.Background(), `
			SELECT id, title, location_name, location_address,
				description, lat, lng, uploaded,
				start, stop, initial_img_id, likes
				FROM task WHERE start > $1 AND stop < $1 ORDER BY num_submissions DESC
		`, time.Now())
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Database error: %v", err),
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
			}, nil
		}
		defer rows.Close()

		tasks_json, res := buildTaskJSON(rows)
		if res != nil {
			return *res, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       tasks_json,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	case "get_recently_uploaded_tasks":
		rows, err := dbConn.Query(context.Background(), `
			SELECT id, title, location_name, location_address,
				description, lat, lng, uploaded,
				start, stop, initial_img_id, likes
				FROM task WHERE start > $1 AND stop < $1 ORDER BY uploaded ASC
		`, time.Now())
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Database error: %v", err),
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
			}, nil
		}
		defer rows.Close()

		tasks_json, res := buildTaskJSON(rows)
		if res != nil {
			return *res, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       tasks_json,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	case "get_images":
		task_id_str, task_id_exists := request.QueryStringParameters["task_id"]
		if !task_id_exists {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Missing required parameters: task_id",
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
			}, nil
		}

		task_id, err := strconv.Atoi(task_id_str)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "Invalid parameters: task_id",
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
			}, nil
		}

		rows, err := dbConn.Query(context.Background(), `
			SELECT id, uploaded, caption
				FROM img WHERE task_id = $1
		`, task_id)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("Database error: %v", err),
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
			}, nil
		}
		defer rows.Close()

		imgs := []ImgRet{}
		for rows.Next() {
			var id int
			var uploaded time.Time
			var caption string
			err := rows.Scan(&id, &uploaded, &caption)
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 500,
					Body:       fmt.Sprintf("Database error: %v", err),
					Headers: map[string]string{
						"Content-Type": "text/plain",
					},
				}, nil
			}

			presigned_req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String("spontaniapp-imgs"),
				Key:    aws.String(fmt.Sprintf("%d", id)),
			})
			presigned_url, err := presigned_req.Presign(7 * 24 * time.Hour)
			if err != nil {
				panic(fmt.Errorf("error in raw video presigned URL: %v", err))
			}

			imgs = append(imgs, ImgRet{
				Id:       id,
				TaskID:   task_id,
				Uploaded: uploaded.Unix(),
				Caption:  caption,
				URL:      presigned_url,
			})
		}

		imgs_json, err := json.Marshal(imgs)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       fmt.Sprintf("JSON marshalling error: %v", err),
				Headers: map[string]string{
					"Content-Type": "text/plain",
				},
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(imgs_json),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Incorrect parameter value: request_type",
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
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
