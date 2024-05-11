package router

import (
	"github.com/toucham/gotitan/server/msg"
)

type RouterAction func(req msg.Request) msg.Response

type Router struct {
	routes []map[string]RouterAction // An array of map for each method ordered as [get, post, put, delete]
}

type Route interface {
	To(*RouterContext)
	ContainRoute(method msg.HttpMethod, route string) bool
	AddRoute(method msg.HttpMethod, route string, action RouterAction) error
}

// RouterContext implements the [context.Context] interface for passing in message info across goroutines
type RouterContext struct {
	Request   msg.Request   // request from read()
	Response  msg.Response  // response from [RouterAction]
	CloseConn bool          // should close connection after sending response
	Done      chan struct{} // if data in channel, then result is ready to be sent
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
func (r *Router) To(rc *RouterContext) {
	defer close(rc.Done) // closes the channel at the end
	req := rc.Request
	if req == nil {
		rc.Response = msg.ServerErrorResponse()
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
			rc.Response = msg.NotFoundResponse()
		} else {
			rc.Response = action(req)
		}
	}
}

func BuildContext(req msg.Request) *RouterContext {
	// parse header and create result accordingly
	return &RouterContext{
		Request:   req,
		CloseConn: false,
		Done:      make(chan struct{}),
	}
}
