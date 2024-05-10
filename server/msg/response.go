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
	return statusLine + headers + r.body
}

func (r *HttpResponse) buildStatusLine() string {
	statusLine := fmt.Sprintf("%s %s %s\n", string(r.version), r.Status.String(), r.Status.GetReason())
	return statusLine
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

func NewHttpResponse() *HttpResponse {
	return &HttpResponse{
		Headers: DefaultResponseHeader(),
	}
}

func ServerErrorResponse() *HttpResponse {
	return &HttpResponse{
		Headers: DefaultResponseHeader(),
		Status:  StatusServerInternalError,
	}
}

func BadRequestResponse() *HttpResponse {
	return &HttpResponse{
		Headers: DefaultResponseHeader(),
		Status:  StatusBadRequest,
	}
}

func NotFoundResponse() *HttpResponse {
	return &HttpResponse{
		Headers: DefaultResponseHeader(),
		Status:  StatusNotFound,
	}
}
