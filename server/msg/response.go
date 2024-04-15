package msg

// Http response structure
type HttpResponse struct {
	HttpMessage
	Status int16
}

func (r *HttpResponse) SetStatus(status int16) *HttpResponse {
	return r
}

func (r *HttpResponse) String() (string, error) {
	return "HTTP OK", nil
}
