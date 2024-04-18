package conn

import (
	"bufio"
	"net"
	"testing"

	"github.com/toucham/gotitan/logger"
	"github.com/toucham/gotitan/server/msg"
	"github.com/toucham/gotitan/server/router"
)

type MockRoute struct {
	IsCalled chan bool
}

func (m *MockRoute) To(*msg.HttpRequest, *router.RouterResult) {
	m.IsCalled <- true
}

func (r *MockRoute) AddRoute(method msg.HttpMethod, route string, action router.RouterAction) error {
	return nil
}

func (r *MockRoute) ContainRoute(method msg.HttpMethod, route string) bool {
	return false
}

const MOCK_GET_REQUEST = `GET /index.html HTTP/1.1
Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)
`

func createMockHttpConn(isCalled chan bool) (HttpConn, net.Conn) {
	conn, input := net.Pipe()
	var timeout int32 = 10000
	ch := make(chan *router.RouterResult, 1)
	route := new(MockRoute)
	route.IsCalled = isCalled
	return HttpConn{
		conn,
		timeout,
		ch,
		route,
		logger.New("MockRouter"),
	}, input
}

// TODO: create unit test for Read()
func TestRead(t *testing.T) {
	isCalled := make(chan bool)
	mock, input := createMockHttpConn(isCalled)
	writer := bufio.NewWriter(input)
	if _, err := writer.WriteString(MOCK_GET_REQUEST); err != nil {
		t.Fatal("test failed; unable to write")
	} else {
		go writer.Flush()
	}
	go mock.Read()
	<-isCalled
}

// TODO: create unit test for Write()
func TestWrite(t *testing.T) {

}
