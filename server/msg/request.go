package msg

import (
	"errors"
	"strconv"
	"strings"

	"github.com/toucham/gotitan/logger"
	"github.com/toucham/gotitan/server/url"
)

type RequestHeaders struct {
	ContentLength int
}

type HttpRequest struct {
	HttpMessage
	Headers RequestHeaders
}

type Request interface {
	AddRequestLine(string) error
	AddHeader(string) error
	AddBody(string) error
}

// parse raw data to instantiate HttpRequest according to HTTP/1.1
func NewRequestFromMsg(msg string) (*HttpRequest, error) {
	req := new(HttpRequest)
	req.headers = make(map[string]string)

	logger := logger.New("HttpRequest")
	lines := strings.Split(msg, "\n")

	// 1) get request line
	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) != 3 {
		return nil, errors.New("incorrect string to parse in request-line")
	}
	req.method = HttpMethod(strings.ToLower(requestLine[0]))
	req.url = url.NewFromReqLine(requestLine[1])
	req.version = requestLine[2]
	if req.version != "HTTP/1.1" {
		return nil, errors.New("incorrect HTTP version, currently only support 1.1")
	}

	// 2) get headers
	bodyIndex := len(lines) + 1
	for i, field := range lines[1:] {
		if field == "" {
			logger.Debug("empty field")
			bodyIndex = i + 1
			break
		}
		kv := strings.SplitN(field, ":", 2)
		if len(kv) == 2 {
			req.headers[kv[0]] = kv[1]
		} else {
			logger.Debug("Incorrect header format: %s", field)
		}
	}

	// 3) get body
	if bodyIndex > len(lines) {
		logger.Debug("No body")
	} else {
		req.body = strings.Join(lines[bodyIndex+1:], "\n")
	}

	return req, nil
}

type RequestBuildState int

const (
	REQUESTLINE_BS RequestBuildState = iota + 1
	HEADERS_BS
	BODY_BS
	COMPLETE_BS
)

func NewRequest() *HttpRequest {
	req := new(HttpRequest)
	req.headers = make(map[string]string)
	return req
}

func (req *HttpRequest) AddRequestLine(line string) error {
	requestLine := strings.Split(line, " ")
	if len(requestLine) != 3 {
		return errors.New("incorrect string to parse in request-line")
	}
	req.method = HttpMethod(strings.ToLower(requestLine[0]))
	req.url = url.NewFromReqLine(requestLine[1])
	req.version = requestLine[2]

	if req.version != "HTTP/1.1" {
		return errors.New("incorrect HTTP version, currently only support 1.1")
	}
	return nil
}

func (req *HttpRequest) AddHeader(line string) error {
	headers := strings.SplitN(line, ":", 2)
	if len(headers) != 2 {
		return errors.New("headers are split into more than two elements")
	}
	key := strings.ToLower(strings.TrimSpace(headers[0]))
	value := strings.TrimSpace(headers[1])
	req.headers[key] = value
	req.addKnownHeaders(key, value)

	return nil
}

func (req *HttpRequest) AddBody(body string) error {
	if len(req.body) > 0 {
		return errors.New("request already contain body")
	}
	req.body = body
	return nil
}

// create [HttpRequest] of each line
func (req *HttpRequest) Next(line string, buildState RequestBuildState) (RequestBuildState, error) {
	var nextState RequestBuildState = 0
	switch buildState {
	case REQUESTLINE_BS:
		requestLine := strings.Split(line, " ")

		if len(requestLine) != 3 {
			return 0, errors.New("incorrect string to parse in request-line")
		}
		req.method = HttpMethod(strings.ToLower(requestLine[0]))
		req.url = url.NewFromReqLine(requestLine[1])
		req.version = requestLine[2]

		if req.version != "HTTP/1.1" {
			return 0, errors.New("incorrect HTTP version, currently only support 1.1")
		}

		nextState = HEADERS_BS
	case HEADERS_BS:
		if line == "" && req.Headers.ContentLength == 0 { // if there is a body
			nextState = COMPLETE_BS
		} else if line == "" {
			nextState = BODY_BS
		} else {
			headers := strings.Split(line, ":")
			key := strings.ToLower(strings.TrimSpace(headers[0]))
			value := strings.TrimSpace(headers[1])
			req.headers[key] = value
			req.addKnownHeaders(key, value)
		}
	case BODY_BS:
		contentLength := req.Headers.ContentLength
		if contentLength > len(req.body) {
			req.body += line
			if contentLength > len(req.body) {
				req.body += "\n"
			} else {
				nextState = COMPLETE_BS
			}
		}
	case COMPLETE_BS:
		return 0, errors.New("complete state should not be called")
	default:
		return 0, errors.New("unknown [HttpRequest] build state")
	}
	return nextState, nil
}

func (req *HttpRequest) addKnownHeaders(key string, value string) {
	switch key {
	case "content-length":
		i, err := strconv.Atoi(value)
		if err == nil {
			req.Headers.ContentLength = i
		}
	}
}
