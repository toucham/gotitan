package url

import "strings"

type Url struct {
	path string
	uri  string
}

func (u *Url) String() string {
	url := "/" + u.path
	if u.uri != "" {
		url = u.uri + url
	}
	return url
}

// Instantiate [Url] from request line in HTTP message
func NewFromReqLine(requestLine string) *Url {
	// TODO: parse request-line to distinguish type of uri

	return &Url{
		path: strings.Trim(requestLine, "/"),
	}
}

// Add value from host to set the correct url
func (u *Url) AddHostHeader(host string) *Url {
	return u
}
