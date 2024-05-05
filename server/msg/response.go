package msg

import "fmt"

type Response interface {
	String() string               // String transforms [Response] into string that can be sent to the client
	SetBody(string, string) error // SetBody sets body into the response with first as body and second as content-type
}

// Http response structure
type HttpResponse struct {
	HttpMessage
	Status  HttpStatus
	Headers ResponseHeaders
	Body    string
}

func (r *HttpResponse) SetBody(body string, contentType string) error {
	r.Headers.ContentLength = len(body)
	r.Headers.ContentType = contentType
	r.body = body
	return nil
}

func (r *HttpResponse) String() string {
	statusLine := r.buildStatusLine()
	headers := r.buildHeaders()
	return statusLine + "\n" + headers + r.body
}

func (r *HttpResponse) buildStatusLine() string {
	return string(r.version) + " " + r.Status.String() + " " + r.Status.GetReason()
}

func (r *HttpResponse) buildHeaders() string {
	headers := ""
	for k, v := range r.headers {
		headers += fmt.Sprintf("%s: %s\n", k, v)
	}
	if len(headers) == 0 {
		headers += "\n"
	}
	return headers
}

func NewResponse() *HttpResponse {
	return &HttpResponse{
		Headers: DefaultResponseHeader(),
	}
}
