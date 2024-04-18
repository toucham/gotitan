package msg

import (
	"errors"
	"strings"

	"github.com/toucham/gotitan/logger"
	"github.com/toucham/gotitan/server/url"
)

type HttpRequest struct {
	HttpMessage
	buildState RequestBuildState
}

// parse raw data to instantiate HttpRequest according to HTTP/1.1
func NewRequest(msg string) (*HttpRequest, error) {
	req := new(HttpRequest)
	req.Headers = make(map[string]string)

	logger := logger.New("HttpRequest")
	lines := strings.Split(msg, "\n")

	// 1) get request line
	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) != 3 {
		logger.Fatal("Message have %d words; there must be 3", len(requestLine))
		return nil, errors.New("incorrect string to parse in request-line")
	}
	req.method = HttpMethod(strings.ToLower(requestLine[0]))
	req.url = url.NewFromReqLine(requestLine[1])
	req.version = requestLine[2]
	if req.version != "HTTP/1.1" {
		logger.Fatal("HTTP request is of version %s", req.version)
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
			logger.Info("Incorrect header format: %s", field)
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
		// if line == "" || line == "\n" {
		req.buildState = BODY_BS
		// }
	case BODY_BS:
		// if line == "" || line == "\n" {
		req.buildState = COMPLETE_BS
		// }
	case COMPLETE_BS:
		return errors.New("complete state should not be called")
	default:
		return errors.New("unknown [HttpRequest] build state")
	}
	return nil
}

func (req *HttpRequest) Complete() {
	if req.buildState == HEADERS_BS {
		req.buildState = COMPLETE_BS
	}
}

func (req *HttpRequest) IsReady() bool {
	return req.buildState == COMPLETE_BS
}
