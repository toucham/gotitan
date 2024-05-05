package msg

import (
	"strconv"
)

type HttpRequest struct {
	HttpMessage
	Headers RequestHeaders
}

type Request interface {
	Message
	IsSafeMethod() bool
}

func NewRequest() *HttpRequest {
	req := new(HttpRequest)
	req.headers = make(map[string]string)
	return req
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
