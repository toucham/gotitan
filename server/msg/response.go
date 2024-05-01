package msg

// Http response structure
type HttpResponse struct {
	*HttpMessage
	Status int16
}

func (r *HttpResponse) String() string {
	return "HTTP OK"
}

type Response interface {
	String() string
}
