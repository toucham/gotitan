package server

type MiddlwareOptions struct {
}

type ReqMiddleware func(req *HttpRequest, opt MiddlwareOptions)
type ResMiddleware func(req *HttpRequest, opt MiddlwareOptions)

// Add middlware for processing requests
func (s *HttpServer) AddReqMiddlware(m ReqMiddleware) {
	s.reqMw = append(s.reqMw, m)
}

// Add middlware for processing response
func (s *HttpServer) AddResMiddlware(m ResMiddleware) {
	s.resMw = append(s.resMw, m)
}
