package conn

import (
	"bufio"
	"net"
	"testing"

	"github.com/toucham/gotitan/server/msg"
	"github.com/toucham/gotitan/server/router"
)

const MOCK_GET_REQUEST = `GET /index.html HTTP/1.1
Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)

`

const MOCK_POST_REQUEST = `POST /help.txt HTTP/1.1
Host: www.example.re
Content-Type: text/plain
Content-Length: 90

Please visit www.example.re for the latest updates!
Another cool body. Hopefully this works
`

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

type MockLogger struct {
	T *testing.T
}

func (l *MockLogger) Debug(format string, v ...any) {
}

func (l *MockLogger) Info(format string, v ...any) {
}

func (l *MockLogger) Warn(format string, v ...any) {
	l.T.Fatalf(format, v...)
}

func (l *MockLogger) Fatal(format string, v ...any) {
}

func createMockHttpConn(isCalled chan bool, log *MockLogger) (HttpConn, net.Conn) {
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
		log,
	}, input
}

// TODO: create unit test for Read()
func TestRead(t *testing.T) {
	isCalled := make(chan bool)
	mockLogger := MockLogger{T: t}

	mockReqs := []string{MOCK_POST_REQUEST}
	for _, r := range mockReqs {
		mock, input := createMockHttpConn(isCalled, &mockLogger)
		writer := bufio.NewWriter(input)
		if _, err := writer.WriteString(r); err != nil {
			t.Fatal("test failed; unable to write")
		} else {
			go func() {
				writer.Flush()
			}()
		}
		go mock.Read()
		<-isCalled
	}
}

func TestReadPersistConnection(t *testing.T) {
	isCalled := make(chan bool)
	mockLogger := MockLogger{T: t}
	mock, input := createMockHttpConn(isCalled, &mockLogger)
	writer := bufio.NewWriter(input)

	mockReqs := []string{MOCK_POST_REQUEST, MOCK_GET_REQUEST}
	for _, r := range mockReqs {
		if _, err := writer.WriteString(r); err != nil {
			t.Fatal("test failed; unable to write")
		}
	}
	go func() {
		writer.Flush()
		input.Close()
	}()
	go mock.Read()
	<-isCalled
}

// TODO: create unit test for Write()
func TestWrite(t *testing.T) {

}
