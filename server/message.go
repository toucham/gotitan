package server

type HttpMethod string

const (
	HTTP_GET    HttpMethod = "get"
	HTTP_POST   HttpMethod = "post"
	HTTP_DELETE HttpMethod = "delete"
	HTTP_PUT    HttpMethod = "put"
)

type HttpResponse struct {
	headers map[string]string
	method  HttpMethod
	body    string
}

type HttpRequest struct {
	headers map[string]string
	method  HttpMethod
	body    string
}
