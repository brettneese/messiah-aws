package Messiah

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// Endpoint is an endpoint handler. Body is the body of the API gateway request.
type Handler interface {
	Handle(Request) interface{}
}

type Request struct {
	events.APIGatewayProxyRequest
	RequestData map[string]interface{}
}

type Response struct {
	events.APIGatewayProxyResponse
	Headers      map[string]string
	StatusCode   int
	ResponseData interface{}
}

// LambdaHandler is an API Gateway Lambda handler, as defined by: https://github.com/aws/aws-lambda-go/blob/master/events/README_ApiGatewayEvent.md
type LambdaHandler func(_ context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func parseRequest(r events.APIGatewayProxyRequest) Request {
	var apiRequest map[string]interface{}

	arr := []byte(r.Body)
	json.Unmarshal(arr, &apiRequest)

	request := Request{
		APIGatewayProxyRequest: r,
		RequestData:            apiRequest,
	}

	return request
}

func parseResponse(res Response) events.APIGatewayProxyResponse {
	if res.Headers != nil {
		res.APIGatewayProxyResponse.Headers = res.Headers
	}

	if res.StatusCode > 0 {
		res.APIGatewayProxyResponse.StatusCode = res.StatusCode
	} else {
		res.APIGatewayProxyResponse.StatusCode = 200
	}

	if body, err := json.Marshal(res.ResponseData); err != nil {
		res.APIGatewayProxyResponse.Body = res.ResponseData.(string)
	} else {
		res.APIGatewayProxyResponse.Body = string(body)
	}

	return res.APIGatewayProxyResponse
}

// GetLambdaHandler accepts an endpoint handler and returns an API Gateway response
func GetLambdaHandler(handler Handler) LambdaHandler {
	return func(_ context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		res := handler.Handle(parseRequest(r)).(Response)

		return parseResponse(res), nil
	}
}
