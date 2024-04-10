package server

type HttpAction func(req *HttpRequest) *HttpResponse

func (s *HttpServer) AddRoute(method HttpMethod, route string, action HttpAction) {

}
