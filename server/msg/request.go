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
	buildState RequestBuildState
	headers    RequestHeaders
}

// parse raw data to instantiate HttpRequest according to HTTP/1.1
func NewRequestFromMsg(msg string) (*HttpRequest, error) {
	req := new(HttpRequest)
	req.Headers = make(map[string]string)

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
			req.Headers[kv[0]] = kv[1]
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
	REQUESTLINE_BS RequestBuildState = iota
	HEADERS_BS
	BODY_BS
	COMPLETE_BS
)

func NewRequest() *HttpRequest {
	req := new(HttpRequest)
	req.Headers = make(map[string]string)
	return req
}

// create [HttpRequest] of each line
func (req *HttpRequest) Next(line string) error {
	switch req.buildState {
	case REQUESTLINE_BS:
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

		req.buildState = HEADERS_BS
	case HEADERS_BS:
		if line == "" { // if there is a body
			req.buildState = BODY_BS
		} else {
			headers := strings.Split(line, ":")
			key := strings.ToLower(strings.TrimSpace(headers[0]))
			value := strings.TrimSpace(headers[1])
			req.Headers[key] = value
			req.addKnownHeaders(key, value)
		}
	case BODY_BS:
		contentLength := req.headers.ContentLength
		if contentLength > len(req.body) {
			req.body += line
			if contentLength > len(req.body) {
				req.body += "\n"
			} else {
				req.buildState = COMPLETE_BS
			}
		}
	case COMPLETE_BS:
		return errors.New("complete state should not be called")
	default:
		return errors.New("unknown [HttpRequest] build state")
	}
	return nil
}

func (req *HttpRequest) Complete() {
	hasNoBody := req.buildState == HEADERS_BS && req.body == "" && req.headers.ContentLength == 0
	if hasNoBody || req.buildState == BODY_BS {
		req.buildState = COMPLETE_BS
	}
}

func (req *HttpRequest) IsReady() bool {
	return req.buildState == COMPLETE_BS
}

func (req *HttpRequest) addKnownHeaders(key string, value string) {
	switch key {
	case "content-length":
		i, err := strconv.Atoi(value)
		if err == nil {
			req.headers.ContentLength = i
		}
	}

}
