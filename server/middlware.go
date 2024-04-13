package server

type MiddlwareOptions struct {
}

type ReqMiddleware func(req *HttpRequest)
type ResMiddleware func(req *HttpRequest)

// Add middlware for processing requests
func (s *HttpServer) AddReqMiddlware(m ReqMiddleware, opt MiddlwareOptions) {
	s.reqMw = append(s.reqMw, m)
}

// Add middlware for processing response
// func (s *HttpServer) AddResMiddlware(m ResMiddleware, opt MiddlwareOptions) {
// 	s.resMw = append(s.resMw, m)
// }
