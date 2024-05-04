package msg

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

type MockRequests struct {
	mock       string
	uri        string
	headersLen int
	body       string
	method     HttpMethod
}

const MOCK_GET_REQUEST = `GET /index.html HTTP/1.1
Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)
Accept: text/html
Accept-Language: en-US, en; q=0.5
Accept-Encoding: gzip, deflate
`

const MOCK_POST_REQUEST = `POST /help.txt HTTP/1.1
Host: www.example.re
Content-Type: text/plain
Content-Length: 90

Please visit www.example.re for the latest updates!
Another cool body. Hopefully this works`

func createMockRequests() []MockRequests {
	return []MockRequests{
		{
			mock:       MOCK_GET_REQUEST,
			uri:        "/index.html",
			headersLen: 5,
			body:       "",
			method:     HTTP_GET,
		},
		{
			mock:       MOCK_POST_REQUEST,
			uri:        "/help.txt",
			headersLen: 3,
			body:       "Please visit www.example.re for the latest updates!\nAnother cool body. Hopefully this works",
			method:     HTTP_POST,
		},
	}
}

func checkAnswer(req *HttpRequest, mr *MockRequests) error {
	checkMethod := req.GetMethod() == mr.method
	if !checkMethod {
		errorMsg := fmt.Sprintf("Incorrect method: %s", req.GetMethod())
		return errors.New(errorMsg)
	}
	checkUri := req.GetUri() == mr.uri
	if !checkUri {
		errorMsg := fmt.Sprintf("Incorrect uri: %s", req.GetUri())
		return errors.New(errorMsg)
	}
	checkHeaders := len(req.headers) == mr.headersLen
	if !checkHeaders {
		errorMsg := fmt.Sprintf("Incorrect len of headers: %d", len(req.headers))
		return errors.New(errorMsg)
	}
	checkBody := req.GetBody() == mr.body
	if !checkBody {
		errorMsg := fmt.Sprintf("Incorrect body: %s", req.GetBody())
		return errors.New(errorMsg)
	}
	return nil
}

func TestNewRequestFromMsg(t *testing.T) {
	mockRequests := createMockRequests()
	for _, mr := range mockRequests {
		req, err := NewRequestFromMsg(mr.mock)
		if err != nil {
			t.Error(err)
		}
		err = checkAnswer(req, &mr)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestAddHeader(t *testing.T) {
	headers := strings.Split(MOCK_POST_REQUEST, "\n")[1:]
	req := NewRequest()
	for _, line := range headers {
		if line == "" {
			break
		}
		req.AddHeader(line)
	}
	if req.Headers.ContentLength == 0 {
		t.Fatal("content lenght is 0 when it is not supposed to")
	}
	if len(req.headers) != 3 {
		t.Fatalf("there are less than expected headers: %d", len(req.headers))
	}
}

func TestAddRequestLine(t *testing.T) {
	inputs := []string{MOCK_GET_REQUEST, MOCK_POST_REQUEST}
	mockReq := createMockRequests()
	for i, input := range inputs {
		req := NewRequest()
		rl := strings.Split(input, "\n")[0]
		t.Logf("request line: %s", rl)
		req.AddRequestLine(rl)
		if req.GetMethod() != mockReq[i].method {
			t.Fatalf("method is not the same: %s", req.GetMethod())
		}
		if req.GetUri() != mockReq[i].uri {
			t.Fatalf("uri is not the same: %s", req.GetUri())
		}
		if req.GetVersion() != "HTTP/1.1" {
			t.Fatalf("uri is not the same: %s", req.GetVersion())
		}
	}
}
