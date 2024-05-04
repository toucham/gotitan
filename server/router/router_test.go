package router

import (
	"testing"

	"github.com/toucham/gotitan/server/msg"
)

func createMockRouter() (*Router, []MockToInput) {
	inputs := []MockToInput{
		MockToInput{
			&msg.HttpRequest{},
			&msg.HttpResponse{},
		},
		MockToInput{
			&msg.HttpRequest{},
			&msg.HttpResponse{},
		},
	}

	// add fake action
	router := New()

	return &router, inputs
}

type MockAddRouterInput struct {
	method msg.HttpMethod
	route  string
}

func TestRouter_AddRouter(t *testing.T) {
	router := New()
	mockInput := []MockAddRouterInput{
		MockAddRouterInput{},
		MockAddRouterInput{},
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

type MockToInput struct {
	req *msg.HttpRequest
	res *msg.HttpResponse
}

// func TestRouter_To(t *testing.T) {
// 	router, inputs := createMockRouter()
// 	for _, input := range inputs {
// 		res := router.To(input.req)

// 		// validate
// 		checkStatus := input.res.Status != res.Status
// 		if checkStatus {
// 			t.Fail()
// 		}
// 	}
// }
