package router

import (
	"testing"

	"github.com/toucham/gotitan/server/msg"
)

type MockAddRouterInput struct {
	method msg.HttpMethod
	route  string
}

func TestRouter_AddRouter(t *testing.T) {
	router := New()
	mockInput := []MockAddRouterInput{
		{
			msg.HTTP_GET,
			"/index",
		},
		{
			msg.HTTP_POST,
			"/route",
		},
	}

	for _, input := range mockInput {
		router.AddRoute(input.method, input.route, func(req msg.Request) msg.Response {
			return nil
		})
	}

	for _, input := range mockInput {
		if !router.ContainRoute(input.method, input.route) {
			t.Fatalf("Doesn't contain route: %s", input.route)
		}
	}
}
