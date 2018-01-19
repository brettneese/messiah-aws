package Messiah

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	. "github.com/smartystreets/goconvey/convey"
)

var MockAPIGatewayProxyRequest = events.APIGatewayProxyRequest{
	Resource:   "/",
	Path:       "/",
	HTTPMethod: "GET",
	Headers: map[string]string{
		"Accept":                       "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Encoding":              "gzip, deflate, br",
		"Accept-Language":              "en-US,en;q=0.5",
		"CloudFront-Forwarded-Proto":   "https",
		"CloudFront-Is-Desktop-Viewer": "true",
		"CloudFront-Is-Mobile-Viewer":  "false",
		"CloudFront-Is-SmartTV-Viewer": "false",
		"CloudFront-Is-Tablet-Viewer":  "false",
		"CloudFront-Viewer-Country":    "US",
		"DNT":                       "1",
		"Host":                      "123abcdefg.execute-api.us-east-1.amazonaws.com",
		"Referer":                   "https://google.com",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:57.0) Gecko/20100101 Firefox/57.0",
		"Via":                       "1.1 abc.cloudfront.net (CloudFront)",
		"X-Amz-Cf-Id":               "abc",
		"X-Amzn-Trace-Id":           "Root=abc",
		"X-Forwarded-For":           "0.0.0.0, 0.0.0.0",
		"X-Forwarded-Port":          "443",
		"X-Forwarded-Proto":         "https",
	},
	QueryStringParameters: map[string]string{},
	PathParameters:        map[string]string{},
	StageVariables:        map[string]string{},
	RequestContext: events.APIGatewayProxyRequestContext{
		AccountID:  "012345678912",
		ResourceID: "abcdefg123",
		Stage:      "Prod",
		RequestID:  "iaaaaaaaa-bbbb-cccc-dddd-123456789abc",
		Identity: events.APIGatewayRequestIdentity{
			CognitoIdentityPoolID:         "",
			AccountID:                     "",
			CognitoIdentityID:             "",
			Caller:                        "",
			APIKey:                        "",
			SourceIP:                      "0.0.0.0",
			CognitoAuthenticationType:     "",
			CognitoAuthenticationProvider: "",
			UserArn:   "",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:57.0) Gecko/20100101 Firefox/57.0",
			User:      "",
		},
		ResourcePath: "/",
		HTTPMethod:   "GET",
		APIID:        "123abcdefg",
	},
}

var MockLambdaContext context.Context

func TestParsingRequest(t *testing.T) {

	Convey("given an ApiGatewayProxyRequest with JSON data in the req.RequestData", t, func() {
		MockAPIGatewayProxyRequest.Body = "{\"Hello\":\"World\"}"
		request := parseRequest(MockLambdaContext, MockAPIGatewayProxyRequest)

		Convey("it should  return a properly unmarshalled RequestData map[string]", func() {
			So(request.RequestData, ShouldContainKey, "Hello")
			So(request.RequestData["Hello"], ShouldEqual, "World")
		})

		Convey("it should return the rest of the properties from ApiGatewayProxyRequest", func() {
			So(MockAPIGatewayProxyRequest, ShouldResemble, MockAPIGatewayProxyRequest)
		})

	})

	Convey("given an ApiGatewayProxyRequest with a string in the req.Body", t, func() {
		MockAPIGatewayProxyRequest.Body = "Hello World"
		request := parseRequest(MockLambdaContext, MockAPIGatewayProxyRequest)

		Convey("it should return the data as a string in the RequestData property", func() {
			So(request.Body, ShouldEqual, "Hello World")
		})

		Convey("it should return the rest of the properties from ApiGatewayProxyRequest", func() {
			So(MockAPIGatewayProxyRequest, ShouldResemble, MockAPIGatewayProxyRequest)
		})

	})

}

func TestParsingResponse(t *testing.T) {
	var MockResponse Response

	Convey("given a Response with a `res.headers` map[string]string property", t, func() {

		MockResponse.Headers = map[string]string{"Hello": "World"}
		APIGatewayProxyResponse := parseResponse(MockResponse)

		Convey("it should return the headers as a map[string]string in `res.APIGatewayProxyResponse.Headers`", func() {

			So(APIGatewayProxyResponse.Headers, ShouldContainKey, "Hello")
			So(APIGatewayProxyResponse.Headers["Hello"], ShouldEqual, "World")
		})
	})

	Convey("given a Response with a valid `res.StatusCode` property", t, func() {
		MockResponse.StatusCode = 200
		APIGatewayProxyResponse := parseResponse(MockResponse)

		Convey("it should return the proper `res.APIGatewayProxyResponse.StatusCode", func() {
			So(APIGatewayProxyResponse.StatusCode, ShouldEqual, 200)
		})
	})

	Convey("given a Response with a valid `res.StatusCode` property set by the net/http package", t, func() {

		MockResponse.StatusCode = http.StatusOK
		APIGatewayProxyResponse := parseResponse(MockResponse)

		Convey("it should return the proper `res.APIGatewayProxyResponse.StatusCode", func() {
			So(APIGatewayProxyResponse.StatusCode, ShouldEqual, 200)
		})

	})

	// @todo
	// Convey("given a Response with an invalid `res.StatusCode` property", t, func() {
	// 	// MockResponse.StatusCode = 10
	// 	// APIGatewayProxyResponse := parseResponse(MockResponse)

	// 	// Convey("it should return the proper `res.APIGatewayProxyResponse.StatusCode", func() {
	// 	// 	So(APIGatewayProxyResponse.StatusCode, ShouldEqual, 500)
	// 	// })

	// })

	Convey("given a Response with valid JSONifiable data in the `res.ResponseData` property", t, func() {
		MockResponse.ResponseData = map[string]string{
			"Hello": "World",
		}
		APIGatewayProxyResponse := parseResponse(MockResponse)

		Convey("it should properly marshall the JSON and return it as a string in the `res.APIGatewayProxyResponse.Body`", func() {
			So(APIGatewayProxyResponse.Body, ShouldEqual, "{\"Hello\":\"World\"}")
		})

	})

	Convey("given a response with invalid JSON in the `res.ResponseData` property", t, func() {
		MockResponse.ResponseData = "hello world"
		APIGatewayProxyResponse := parseResponse(MockResponse)

		Convey("it should return a string in the `res.APIGatewayProxyResponse.Body`", func() {
			So(APIGatewayProxyResponse.Body, ShouldEqual, "\"hello world\"")
		})
	})

}

type MockHandler struct{}

func (handler MockHandler) Handle(req Request) interface{} {
	res := Response{
		StatusCode:   200,
		ResponseData: map[string]interface{}{"hello": "world"},
	}
	return res
}

func TestGetLambdaHandler(t *testing.T) {
	handler := MockHandler{}

	Convey("given a MockHandler, should return a non-nil lambdaHandler function", t, func() {
		So(GetLambdaHandler(handler), ShouldNotBeNil)
	})
}
