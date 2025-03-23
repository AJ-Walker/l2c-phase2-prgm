package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	Status     bool   `json:"status"`
	Data       any    `json:"data"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Print("Inside lambdaHandler func")
	log.Printf("Context: %v\n", ctx)
	log.Printf("Event: %v\n", event)

	log.Printf("Resource: %v\n", event.Resource)
	log.Printf("Query Params: %v\n", event.QueryStringParameters)

	switch event.Resource {
	case "/api/movies":
		// call movies func

		if year, ok := event.QueryStringParameters["year"]; ok {
			return getMoviesByYear(year)
		} else {
			return getMovies()
		}

	case "/api/movies/summary":
		// call movies summary func

		movieId, ok := event.QueryStringParameters["movieId"]

		if !ok {
			return response(http.StatusNotFound, false, "movieId query param missing", nil), nil
		}
		return getMoviesSummary(movieId)
	}
	return response(http.StatusInternalServerError, false, "Wrong path provided", nil), nil
}

func main() {
	log.Print("Inside main func")
	lambda.Start(HandleRequest)
}

func getMovies() (events.APIGatewayProxyResponse, error) {
	log.Print("Inside getMovies func")

	return response(http.StatusOK, true, "getMovies works", nil), nil
}

func getMoviesByYear(year string) (events.APIGatewayProxyResponse, error) {
	log.Print("Inside getMoviesByYear func")
	if year == "" {
		return response(http.StatusBadRequest, false, "year field missing", nil), nil
	}

	return response(http.StatusOK, true, fmt.Sprintf("Year: %v", year), nil), nil
}

func getMoviesSummary(movieId string) (events.APIGatewayProxyResponse, error) {
	log.Print("Inside getMoviesSummary func")

	if movieId == "" {
		return response(http.StatusBadRequest, false, "movieId cannot be empty", nil), nil
	}
	return response(http.StatusOK, true, fmt.Sprintf("movieId: %v", movieId), nil), nil
}

func response(statusCode int, status bool, message string, data any) events.APIGatewayProxyResponse {
	log.Print("Inside response func")

	res := Response{
		Status:     status,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
	log.Printf("res: %v", res)
	jsonRes, err := json.Marshal(res)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}
	}

	log.Printf("jsonRes: %v", jsonRes)
	log.Printf("string(jsonRes): %v", string(jsonRes))

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(jsonRes),
	}
}
