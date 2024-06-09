package router

import (
	"fmt"

	"github.com/toucham/gotitan/server/msg"
)

type RouterAction func(req msg.Request) msg.Response

type Router struct {
	routes []map[string]RouterAction // An array of map for each method ordered as [get, post, put, delete]
}

type Route interface {
	To(msg.Request) msg.Response
	ContainRoute(method msg.HttpMethod, route string) bool
	AddRoute(method msg.HttpMethod, route string, action RouterAction)
}

func New() Router {
	r := make([]map[string]RouterAction, 4)
	for i := range r {
		r[i] = make(map[string]RouterAction)
	}
	return Router{
		routes: r,
	}
}

// Simple add route to dictionary
func (r *Router) AddRoute(method msg.HttpMethod, route string, action RouterAction) {
	switch method {
	case msg.HTTP_GET:
		r.routes[0][route] = action
	case msg.HTTP_POST:
		r.routes[1][route] = action
	case msg.HTTP_PUT:
		r.routes[2][route] = action
	case msg.HTTP_DELETE:
		r.routes[3][route] = action
	default:
		panic(fmt.Sprintf("Undefined method: %s", route))
	}
}

func (r *Router) ContainRoute(method msg.HttpMethod, route string) bool {
	var action RouterAction = nil
	switch method {
	case msg.HTTP_GET:
		action = r.routes[0][route]
	case msg.HTTP_POST:
		action = r.routes[1][route]
	case msg.HTTP_PUT:
		action = r.routes[2][route]
	case msg.HTTP_DELETE:
		action = r.routes[3][route]
	default:
		return false
	}
	return action != nil
}

// Route [HttpRequest] to the correct action depending on the path
func (r *Router) To(req msg.Request) msg.Response {
	// defer close(rc.Done) // closes the channel at the end
	if req == nil {
		return msg.CreateHttpResponse(msg.StatusServerInternalError)
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
			return msg.CreateHttpResponse(msg.StatusNotFound)
		} else {
			return action(req)
		}
	}
}
