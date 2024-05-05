package conn

import (
	"net"

	"github.com/toucham/gotitan/server/msg"
	"github.com/toucham/gotitan/server/router"
)

type MockRoute struct {
}

func (m *MockRoute) To(msg msg.Request, r *router.RouterContext) {
	r.Ready <- true
}

func (r *MockRoute) AddRoute(method msg.HttpMethod, route string, action router.RouterAction) error {
	return nil
}

func (r *MockRoute) ContainRoute(method msg.HttpMethod, route string) bool {
	return false
}

type MockLogger struct {
}

func (l *MockLogger) Debug(format string, v ...any) {
}

func (l *MockLogger) Info(format string, v ...any) {
}

func (l *MockLogger) Warn(format string, v ...any) {
}

func (l *MockLogger) Fatal(format string, v ...any) {
}

func createMockHttpConn() (HttpConn, net.Conn) {
	conn, input := net.Pipe()
	var timeout int32 = 10000
	route := new(MockRoute)
	queue := make(chan *router.RouterContext)
	return HttpConn{
		conn,
		timeout,
		queue,
		true,
		msg.NewHttpReqBuilder(),
		route,
		new(MockLogger),
	}, input
}

type MockResult struct {
	Mock      string
	NumOfReqs int
}

type MockResponse struct{}

func (m MockResponse) SetBody(body string, contentType string) error {
	return nil
}

const EXPECTED_RESP_STRING = "HTTP OK"

func (r MockResponse) String() string {
	return EXPECTED_RESP_STRING
}

func createMockCtx() router.RouterContext {
	res := MockResponse{}
	return router.RouterContext{
		Response: res,
		Ready:    make(chan bool),
	}
}
