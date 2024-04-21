package msg

import (
	"github.com/toucham/gotitan/server/url"
)

type HttpMethod string

const (
	HTTP_GET     HttpMethod = "get"
	HTTP_POST    HttpMethod = "post"
	HTTP_DELETE  HttpMethod = "delete"
	HTTP_PUT     HttpMethod = "put"
	HTTP_OPTIONS HttpMethod = "options"
	HTTP_TRACE   HttpMethod = "trace"
	HTTP_HEAD    HttpMethod = "head"
)

type HttpMessage struct {
	headers map[string]string
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

func (r *HttpMessage) IsSafeMethod() bool {
	return r.method == HTTP_GET ||
		r.method == HTTP_OPTIONS ||
		r.method == HTTP_TRACE ||
		r.method == HTTP_HEAD
}
