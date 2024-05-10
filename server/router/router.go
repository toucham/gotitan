package router

import "github.com/toucham/gotitan/server/msg"

type RouterAction func(req msg.Request) msg.Response

type Router struct {
	routes []map[string]RouterAction // An array of map for each method ordered as [get, post, put, delete]
}

type Route interface {
	To(msg.Request, *RouterContext)
	ContainRoute(method msg.HttpMethod, route string) bool
	AddRoute(method msg.HttpMethod, route string, action RouterAction) error
}

// RouterContext implements the [context.Context] interface for passing in message info across goroutines
type RouterContext struct {
	Response  msg.Response // response from [RouterAction]
	CloseConn bool         // should close connection after sending response
	Ready     chan bool    // if data in channel, then result is ready to be sent
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
func (r *Router) To(req msg.Request, result *RouterContext) {
	if req == nil {
		result.Response = msg.ServerErrorResponse()
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
		default:
			result.Ready <- false
		}
		if action == nil {
			result.Response = msg.NotFoundResponse()
		} else {
			result.Response = action(req)
		}
	}

	// response ready to be sent to client
	result.Ready <- true
}

func CreateContext() *RouterContext {
	// parse header and create result accordingly
	return &RouterContext{
		CloseConn: false,
		Ready:     make(chan bool),
	}
}
