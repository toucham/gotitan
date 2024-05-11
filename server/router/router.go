package router

import (
	"github.com/toucham/gotitan/server/msg"
)

type RouterAction func(req msg.Request) msg.Response

type Router struct {
	routes []map[string]RouterAction // An array of map for each method ordered as [get, post, put, delete]
}

type Route interface {
	To(msg.Request) msg.Response
	ContainRoute(method msg.HttpMethod, route string) bool
	AddRoute(method msg.HttpMethod, route string, action RouterAction) error
}

func New() Router {
	return Router{
		routes: make([]map[string]RouterAction, 4),
	}
}

func (r *Router) AddRoute(method msg.HttpMethod, route string, action RouterAction) error {
	return nil
}

func (r *Router) ContainRoute(method msg.HttpMethod, route string) bool {
	return false
}

// Route [HttpRequest] to the correct action depending on the path
func (r *Router) To(req msg.Request) msg.Response {
	// defer close(rc.Done) // closes the channel at the end
	if req == nil {
		return msg.ServerErrorResponse()
	} else {
		var action RouterAction

		switch req.GetMethod() {
		case msg.HTTP_GET:
			action = r.routes[0][req.GetPath()]
		case msg.HTTP_POST:
			action = r.routes[1][req.GetPath()]
		case msg.HTTP_PUT:
			action = r.routes[2][req.GetPath()]
		case msg.HTTP_DELETE:
			action = r.routes[3][req.GetPath()]
		}
		if action == nil {
			return msg.NotFoundResponse()
		} else {
			return action(req)
		}
	}
}
