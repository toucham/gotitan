package msg

import "strconv"

type HttpStatus int16

const (
	// 1xx - Informational
	StatusContinue HttpStatus = 100
	// 2xx - Successful
	StatusOk        HttpStatus = 200
	StatusCreated   HttpStatus = 201
	StatusAccepted  HttpStatus = 202
	StatusNoContent HttpStatus = 204
	// 3xx - Redirection

	// 4xx - Client Error
	StatusBadRequest   HttpStatus = 400
	StatusUnauthorized HttpStatus = 401
	StatusNotFound     HttpStatus = 404

	// 5xx - Server Error
	StatusServerInternalError HttpStatus = 500
	StatusNotImplemented      HttpStatus = 501
	StatusBadGateway          HttpStatus = 502
	StatusServiceUnavailable  HttpStatus = 503
)

func (s HttpStatus) IsValid() bool {
	if s >= 100 && s < 600 {
		return true
	}
	return false
}

func (s HttpStatus) String() string {
	return strconv.Itoa(int(s))
}

func (s HttpStatus) GetReason() string {
	switch s {
	case StatusOk:
		return "OK"
	case StatusCreated:
		return "Created"
	case StatusBadRequest:
		return "Bad Request"
	case StatusServerInternalError:
		return "Server Internal Error"
	case StatusNotFound:
		return "Not Found"
	default:
		return "Not Found Status"
	}
}
