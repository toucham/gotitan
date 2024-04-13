package server

import "github.com/toucham/gotitan/server/msg"

type HttpAction func(req *msg.HttpRequest) *msg.HttpResponse

func (s *HttpServer) AddRoute(method msg.HttpMethod, route string, action HttpAction) {

}
