package Messiah

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is messiah's generic function type. Handlers accept a generic messiah.Request type and return a messiah.Response. Those values then get set on the `events.APIGatewayProxyResponse.`
type Handler interface {
	Handle(Request) interface{}
}

// Request is messiah's generic HTTP-like Request type. RequestData is the unmarshalled JSON data from the request. By embedding the events.APIGatewayProxyRequest, you also have access to all of the other properties from events.APIGatewayProxyRequest.
type Request struct {
	Context context.Context
	events.APIGatewayProxyRequest
	RequestData map[string]interface{}
}

// Response is messiah's generic HTTP-like Response type. Response abstract away the essential API Gateway properties from the events.APIGatewayProxyResponse type, ensuring clean seperation from Lambda/API Gateway and your business logic.
// These properties are then processed and appended to a before returning that function to be processed by Lambda. If ResponseData is able to be marshalled into JSON, then it is and. Otherwise it is returned to the client as a string.
type Response struct {
	events.APIGatewayProxyResponse
	Headers      map[string]string
	StatusCode   int
	ResponseData interface{}
}

// LambdaHandler is an API Gateway Lambda handler, as defined by: https://github.com/aws/aws-lambda-go/blob/master/events/README_ApiGatewayEvent.md
type LambdaHandler func(_ context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// parseRequest accepts a events.APIGatewayProxyRequest type, attempts to unmarshall the APIGatewayProxyRequest from JSON, and returns a generic Messiah.Request type, merging in all the values from events.APIGatewayProxyRequest
func parseRequest(ctx context.Context, r events.APIGatewayProxyRequest) Request {
	var apiRequest map[string]interface{}

	arr := []byte(r.Body)
	json.Unmarshal(arr, &apiRequest)

	request := Request{
		Context:                ctx,
		APIGatewayProxyRequest: r,
		RequestData:            apiRequest,
	}

	return request
}

// ParseResponse accepts a generic messiah.Response type and returns a properly parsed res.APIGatewayProxyResponse with the proper values set from the messiah.Response on to the res.APIGatewayProxyResponse
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
		// if the marshalling fails, just stringify whatever's in `res.ResponseData`
		res.APIGatewayProxyResponse.Body = string(res.ResponseData.(string))
	} else {
		res.APIGatewayProxyResponse.Body = string(body)
	}

	return res.APIGatewayProxyResponse
}

// GetLambdaHandler accepts a generic endpoint handler and returns a generic lambdaHandler function, which should be fed into lambda.Start().
func GetLambdaHandler(handler Handler) LambdaHandler {
	return func(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		res := handler.Handle(parseRequest(ctx, r)).(Response)

		return parseResponse(res), nil
	}
}

// Start accepts a generic Messiah endpoint handler (as defined by Messiah.Handler), translates it into a vendor-specific handler, and passes that handler into lambda.Start()
func Start(handler Handler) {
	h := GetLambdaHandler(handler)
	lambda.Start(h)
}
