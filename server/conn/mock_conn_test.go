package conn

import (
	"net"
	"testing"

	"github.com/toucham/gotitan/server/msg"
	"github.com/toucham/gotitan/server/router"
)

type MockRoute struct {
}

func (m *MockRoute) To(msg *msg.HttpRequest, r *router.RouterContext) {
	r.Ready <- true
}

func (r *MockRoute) AddRoute(method msg.HttpMethod, route string, action router.RouterAction) error {
	return nil
}

func (r *MockRoute) ContainRoute(method msg.HttpMethod, route string) bool {
	return false
}

type MockLogger struct {
	T *testing.T
}

func (l *MockLogger) Debug(format string, v ...any) {
}

func (l *MockLogger) Info(format string, v ...any) {
}

func (l *MockLogger) Warn(format string, v ...any) {
}

func (l *MockLogger) Fatal(format string, v ...any) {
}

func createMockHttpConn(log *MockLogger) (HttpConn, net.Conn, chan *router.RouterContext) {
	conn, input := net.Pipe()
	var timeout int32 = 10000
	route := new(MockRoute)
	ch := make(chan *router.RouterContext)
	return HttpConn{
		conn,
		timeout,
		ch,
		true,
		route,
		log,
	}, input, ch
}

type MockResult struct {
	Mock      string
	NumOfReqs int
}
