package msg

import (
	"strings"
)

type HttpMethod string

type HttpVersion string

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
	body    string
	url     *Uri
	method  HttpMethod
	version HttpVersion
}

type Message interface {
	GetMethod() HttpMethod
	GetVersion() HttpVersion
	GetBody() string
	GetPath() string
	GetUri() string
}

// getter method for body field in HttpMessage
func (r *HttpMessage) GetBody() string {
	return r.body
}

// getter method for HTTP method in HttpMessage
func (r *HttpMessage) GetMethod() HttpMethod {
	return r.method
}

// getter method for uri in HttpMessage
func (r *HttpMessage) GetUri() string {
	return r.url.String()
}

// getter method for path in HttpMessage
func (r *HttpMessage) GetPath() string {
	return r.url.String() // TODO: change to get path
}

// getter method for HTTP version in HttpMessage
func (r *HttpMessage) GetVersion() HttpVersion {
	return r.version // TODO: change to get path
}

func (r *HttpMessage) IsSafeMethod() bool {
	return r.method == HTTP_GET ||
		r.method == HTTP_OPTIONS ||
		r.method == HTTP_TRACE ||
		r.method == HTTP_HEAD
}

func toHttpMethod(method string) HttpMethod {
	method = strings.ToLower(method)
	isValid := method == string(HTTP_DELETE) ||
		method == string(HTTP_GET) ||
		method == string(HTTP_POST) ||
		method == string(HTTP_OPTIONS) ||
		method == string(HTTP_PUT)
	if isValid {
		return HttpMethod(method)
	}
	return ""
}

func toHttpVersion(version string) HttpVersion {
	msgs := strings.Split(version, "/")
	if len(msgs) == 2 {
		isHttp := msgs[0] == "HTTP"
		ver := strings.Split(msgs[1], ".")
		isSupportedVer := ver[0] == "1"
		isValid := isHttp && isSupportedVer
		if isValid {
			return HttpVersion(version)
		}
	}
	return ""
}
