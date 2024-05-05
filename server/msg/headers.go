package msg

// connection options for HTTP message
type ConnectionOpt string

const (
	ConnOpt_KeepAlive ConnectionOpt = "keep-alive"
	ConnOpt_Close     ConnectionOpt = "close"
)

// headers in HTTP request
type RequestHeaders struct {
	ContentLength    int
	Host             string
	Connection       ConnectionOpt
	TransferEncoding string
}

// headers in HTTP response
type ResponseHeaders struct {
}
