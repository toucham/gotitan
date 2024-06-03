package msg

import (
	"errors"
	"fmt"
	"strings"
)

type RequestBuilderState int

const (
	RequestLineBuildState RequestBuilderState = iota + 1
	HeadersBuildState
	BodyBuildState
	CompleteBuildState
)

type RequestBuilder interface {
	Build() (Request, error)
	AddRequestLine(string) error
	AddHeader(string) error
	AddBody(string) error
	State() RequestBuilderState
}

type HttpRequestBuilder struct {
	state   RequestBuilderState
	request *HttpRequest
}

func NewHttpReqBuilder() *HttpRequestBuilder {
	return &HttpRequestBuilder{
		state:   RequestLineBuildState,
		request: NewRequest(),
	}
}

// AddRequestLine parses a line to fill request line data into a [HttpRequest]
//
// A "line" should not include CRLF ("\n")
func (builder *HttpRequestBuilder) AddRequestLine(line string) error {
	if builder.state != RequestLineBuildState {
		return errors.New("incorrect build state")
	}

	requestLine := strings.Split(line, " ")
	if len(requestLine) != 3 {
		return errors.New("incorrect string to parse in request-line")
	}
	req := builder.request
	if req.method = toHttpMethod(requestLine[0]); req.method == "" {
		return errors.New("incorrect HTTP method")
	}
	if req.version = toHttpVersion(requestLine[2]); req.version == "" {
		return errors.New("incorrect HTTP version")
	}
	if req.url = ParseUri(requestLine[1]); req.url == nil {
		return errors.New("incorrect HTTP path")
	}
	builder.state = HeadersBuildState
	return nil
}

// AddHeader parses a "line" to fill in headers property and add to headers map.
// It will parse successfully if the line is aligned with the format stated in HTTP/1.1 RFC
//
// A "line" should include CRLF ("\n") at the end
func (builder *HttpRequestBuilder) AddHeader(line string) error {
	if builder.state != HeadersBuildState {
		return errors.New("incorrect build state")
	}

	req := builder.request

	if line == "" || line == "\r" { // curl -> carriage return
		if req.Headers.ContentLength > 0 { // expect body
			builder.state = BodyBuildState
		} else {
			builder.state = CompleteBuildState
		}
		return nil
	}

	headers := strings.SplitN(line, ":", 2)
	if len(headers) != 2 {
		errorMsg := fmt.Sprintf("headers are not split into two elements: %d", len(headers))
		return errors.New(errorMsg)
	}

	key := strings.ToLower(strings.TrimSpace(headers[0]))
	value := strings.TrimSpace(headers[1])
	req.headers[key] = value
	req.addKnownHeaders(key, value)

	return nil
}

// AddBody expects an entire string that is extracted from the message body
func (builder *HttpRequestBuilder) AddBody(body string) error {
	if builder.state != BodyBuildState {
		return errors.New("incorrect build state")
	}
	req := builder.request

	if len(req.body) > 0 {
		return errors.New("request already contain body")
	}
	if len(body) == req.Headers.ContentLength {
		req.body = body
		builder.state = CompleteBuildState
		return nil
	} else if len(body) > req.Headers.ContentLength {
		return errors.New("body has is longer than content length")
	} else {
		return nil // not enough byte & not moved to next state
	}
}

func (builder *HttpRequestBuilder) Build() (Request, error) {
	if builder.state != CompleteBuildState {
		return nil, errors.New("incorrect state")
	}
	req := builder.request
	builder.request = NewRequest()
	builder.state = RequestLineBuildState
	return req, nil
}

func (builder *HttpRequestBuilder) State() RequestBuilderState {
	return builder.state
}
