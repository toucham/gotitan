package msg

import "strings"

type Uri struct {
	path string
	uri  string
}

func (u *Uri) String() string {
	url := "/" + u.path
	if u.uri != "" {
		url = u.uri + url
	}
	return url
}

// Instantiate [Uri] from request line in HTTP message
func ParseUri(uri string) *Uri {
	// TODO: parse request-line to distinguish type of uri

	return &Uri{
		path: strings.Trim(uri, "/"),
	}
}

// Add value from host to set the correct url
func (u *Uri) AddHostHeader(host string) *Uri {
	return u
}
