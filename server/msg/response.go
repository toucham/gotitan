package msg

// Http response structure
type HttpResponse struct {
	HttpMessage
	Status  HttpStatus
	Headers ResponseHeaders
}

func (r *HttpResponse) String() string {
	return "HTTP OK"
}

func NewResponse() *HttpResponse {
	return &HttpResponse{}
}

type Response interface {
	String() string
}
