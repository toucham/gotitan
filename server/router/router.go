package router

import "github.com/toucham/gotitan/server/msg"

type RouterAction func(req *msg.HttpRequest) *msg.HttpResponse

type Router struct {
	routes []map[string]RouterAction // An array of map for each method ordered as [get, post, put, delete]
}

type RouterResult struct {
	Response  *msg.HttpResponse // response from [RouterAction]
	CloseConn bool              // if "Connection" header has "close" as value
	Ready     chan bool         // if data in channel, then result is ready to be sent
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
func (r *Router) To(req *msg.HttpRequest, result *RouterResult) {
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
		result.Response = new(msg.HttpResponse).SetStatus(404)
	} else {
		result.Response = action(req)
	}
	result.Ready <- true
}

func CreateResult(req *msg.HttpRequest) *RouterResult {
	// parse header and create result accordingly
	return &RouterResult{
		CloseConn: false,
		Ready:     make(chan bool),
	}
}
