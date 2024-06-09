package msg

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/toucham/gotitan/logger"
)

type Response interface {
	String() (string, error)      // String transforms [Response] into string that can be sent to the client
	SetBody(string, string) error // SetBody sets body into the response with first as body and second as content-type
}

// Http response structure
type HttpResponse struct {
	HttpMessage
	Status HttpStatus
	body   string
	log    logger.Logger
}

func (r *HttpResponse) SetBody(body string, contentType string) error {
	r.headers["content-length"] = strconv.Itoa(len(body))
	r.headers["content-type"] = contentType
	r.body = body
	return nil
}

func (r *HttpResponse) String() (string, error) {
	statusLine, statusErr := r.buildStatusLine()
	if statusErr != nil {
		return "", statusErr
	}
	headers, headerErr := r.buildHeaders()
	if headerErr != nil {
		return "", statusErr
	}
	body, bodyErr := r.buildBody()
	if bodyErr != nil {
		return "", bodyErr
	}
	resString := statusLine + headers + body
	r.log.Info(fmt.Sprintf("response: %s", statusLine))
	return resString, nil
}

// validate and build status line
func (r *HttpResponse) buildStatusLine() (string, error) {
	if !r.Status.IsValid() {
		return "", errors.New("status is not valid")
	}
	statusLine := fmt.Sprintf("%s %s %s\n", string(r.GetVersion()), r.Status.String(), r.Status.GetReason())
	return statusLine, nil
}

// validate and build headers
func (r *HttpResponse) buildHeaders() (string, error) {
	headers := ""
	for k, v := range r.headers {
		headers += fmt.Sprintf("%s: %s\n", k, v)
	}
	headers += "\n"
	return headers, nil
}

func (r *HttpResponse) buildBody() (string, error) {
	// validate body
	// return body
	return r.body, nil
}

// CreateHttpResponse create a [HttpResponse]
func CreateHttpResponse(status HttpStatus) *HttpResponse {
	return &HttpResponse{
		Status: status,
		HttpMessage: HttpMessage{
			version: HttpVersion("HTTP/1.1"),
			headers: make(map[string]string),
		},
		log: logger.New("HttpResponse"),
	}
}
