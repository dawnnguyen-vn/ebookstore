package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/aws/aws-lambda-go/events"
	"github.com/labstack/echo/v4"
)

type LambdaAdapter struct {
	Echo *echo.Echo
}

func (a *LambdaAdapter) Handler(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	req, err := a.ProxyEventToHTTPRequest(event)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	resWriter := NewProxyResponseWriter()

	a.Echo.ServeHTTP(http.ResponseWriter(resWriter), req)

	res, err := resWriter.GetProxyResponse()

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return res, nil
}

func (a *LambdaAdapter) ProxyEventToHTTPRequest(event events.APIGatewayProxyRequest) (*http.Request, error) {
	decodedBody := []byte(event.Body)
	if event.IsBase64Encoded {
		base64Body, err := base64.StdEncoding.DecodeString(event.Body)
		if err != nil {
			return nil, err
		}
		decodedBody = base64Body
	}

	queryString := ""
	if len(event.QueryStringParameters) > 0 {
		queryString = "?"
		queryCnt := 0
		for q := range event.QueryStringParameters {
			if queryCnt > 0 {
				queryString += "&"
			}
			queryString += url.QueryEscape(q) + "=" + url.QueryEscape(event.QueryStringParameters[q])
			queryCnt++
		}
	}

	path := event.Path
	httpRequest, err := http.NewRequest(
		strings.ToUpper(event.HTTPMethod),
		path+queryString,
		bytes.NewReader(decodedBody),
	)

	if err != nil {
		fmt.Printf("Could not convert request %s:%s to http.Request\n", event.HTTPMethod, event.Path)
		return nil, err
	}

	for h := range event.Headers {
		httpRequest.Header.Add(h, event.Headers[h])
	}

	return httpRequest, nil
}

// ProxyResponseWriter implements http.ResponseWriter and adds the method
// necessary to return an events.APIGatewayProxyResponse object
type ProxyResponseWriter struct {
	headers http.Header `json:"headers"`
	body    []byte      `json:"body"`
	status  int         `json:"statusCode"`
}

// NewProxyResponseWriter returns a new ProxyResponseWriter object.
// The object is initialized with an empty map of headers and a
// status code of -1
func NewProxyResponseWriter() *ProxyResponseWriter {
	return &ProxyResponseWriter{
		headers: make(http.Header),
		status:  http.StatusOK,
	}

}

// Header implementation from the http.ResponseWriter interface.
func (r *ProxyResponseWriter) Header() http.Header {
	return r.headers
}

// Write sets the response body in the object. If no status code
// was set before with the WriteHeader method it sets the status
// for the response to 200 OK.
func (r *ProxyResponseWriter) Write(body []byte) (int, error) {
	r.body = body
	if r.status == -1 {
		r.status = http.StatusOK
	}

	return len(body), nil
}

// WriteHeader sets a status code for the response. This method is used
// for error responses.
func (r *ProxyResponseWriter) WriteHeader(status int) {
	r.status = status
}

// GetProxyResponse converts the data passed to the response writer into
// an events.APIGatewayProxyResponse object.
// Returns a populated proxy response object. If the reponse is invalid, for example
// has no headers or an invalid status code returns an error.
func (r *ProxyResponseWriter) GetProxyResponse() (events.APIGatewayProxyResponse, error) {
	if len(r.headers) == 0 {
		return events.APIGatewayProxyResponse{}, errors.New("no headers generated for response")
	}

	var output string
	isBase64 := false

	if utf8.Valid(r.body) {
		output = string(r.body)
	} else {
		output = base64.StdEncoding.EncodeToString(r.body)
		isBase64 = true
	}

	proxyHeaders := make(map[string]string)

	for h := range r.headers {
		proxyHeaders[h] = r.headers.Get(h)
	}

	return events.APIGatewayProxyResponse{
		StatusCode:      r.status,
		Headers:         proxyHeaders,
		Body:            output,
		IsBase64Encoded: isBase64,
	}, nil
}
