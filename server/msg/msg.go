package msg

import (
	"errors"
	"strings"

	"github.com/toucham/gotitan/logger"
	"github.com/toucham/gotitan/server/url"
)

type HttpMethod string

const (
	HTTP_GET    HttpMethod = "get"
	HTTP_POST   HttpMethod = "post"
	HTTP_DELETE HttpMethod = "delete"
	HTTP_PUT    HttpMethod = "put"
)

type HttpMessage struct {
	Headers map[string]string
	method  HttpMethod
	body    string
	url     *url.Url
	version string
}

// getter method for body field in HttpRequest
func (r *HttpMessage) GetBody() string {
	return r.body
}

// getter method for HTTP method in HttpRequest
func (r *HttpMessage) GetMethod() HttpMethod {
	return r.method
}

// getter method for uri in HttpRequest
func (r *HttpMessage) GetUri() string {
	return r.url.String()
}

// getter method for path in HttpRequest
func (r *HttpMessage) GetPath() string {
	return r.url.String() // TODO: change to get path
}

type HttpRequest struct {
	HttpMessage
}

// parse raw data to instantiate HttpRequest according to HTTP/1.1
func NewRequest(msg string) (*HttpRequest, error) {
	req := new(HttpRequest)
	req.Headers = make(map[string]string)

	logger := logger.New()
	lines := strings.Split(msg, "\n")

	// 1) get request line
	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) != 3 {
		logger.Fatal("Message have %d words; there must be 3", len(requestLine))
		return nil, errors.New("incorrect string to parse in request-line")
	}
	req.method = HttpMethod(strings.ToLower(requestLine[0]))
	// req.url = requestLine[1]
	req.version = requestLine[2]
	if req.version != "HTTP/1.1" {
		logger.Fatal("HTTP request is of version %s", req.version)
		return nil, errors.New("incorrect HTTP version, currently only support 1.1")
	}

	// 2) get headers
	bodyIndex := len(lines) + 1
	for i, field := range lines[1:] {
		if field == "" {
			logger.Info("empty field")
			bodyIndex = i + 1
			break
		}
		kv := strings.SplitN(field, ":", 2)
		if len(kv) == 2 {
			req.Headers[kv[0]] = kv[1]
		} else {
			logger.Info("Incorrect header format: %s", field)
		}
	}

	// 3) get body
	if bodyIndex > len(lines) {
		logger.Info("No body")
	} else {
		req.body = strings.Join(lines[bodyIndex+1:], "\n")
	}

	return req, nil
}

// Http response structure
type HttpResponse struct {
	HttpMessage
	Status int16
}

func (r *HttpResponse) SetStatus(status int16) *HttpResponse {
	return r
}

func (r *HttpResponse) String() (string, error) {
	return "HTTP OK", nil
}
