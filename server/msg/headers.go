package msg

import "time"

// connection options for HTTP message
type ConnectionOpt string

const (
	ConnOpt_KeepAlive ConnectionOpt = "keep-alive"
	ConnOpt_Close     ConnectionOpt = "close"
)

type Headers struct {
	ContentLength int
	ContentType   string
}

// headers in HTTP request
type RequestHeaders struct {
	Headers
	Host             string
	Connection       ConnectionOpt
	TransferEncoding string
}

// headers in HTTP response
type ResponseHeaders struct {
	Headers
	Date         time.Time
	LastModified time.Time
	Server       string
	Vary         string
}

func DefaultResponseHeader() ResponseHeaders {
	return ResponseHeaders{
		Date: time.Now(),
	}
}
