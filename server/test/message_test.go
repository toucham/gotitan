package server

import (
	"server"
	"testing"
)

type MockRequests struct {
	mock       string
	uri        string
	headersLen int
	body       string
	method     server.HttpMethod
}

const MOCK_GET_REQUEST = `GET /index.html HTTP/1.1
	Host: www.example.re
	User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)
	Accept: text/html
	Accept-Language: en-US, en; q=0.5
	Accept-Encoding: gzip, deflate`

const MOCK_POST_REQUEST = `POST /help.txt HTTP/1.1
Host: www.example.re
Content-Type: text/plain
Content-Length: 51

Please visit www.example.re for the latest updates!
Another cool body. Hopefully this works`

func TestExtractRequest(t *testing.T) {
	mockRequests := []MockRequests{
		{
			mock:       MOCK_GET_REQUEST,
			uri:        "/index.html",
			headersLen: 5,
			body:       "",
			method:     server.HTTP_GET,
		},
		{
			mock:       MOCK_POST_REQUEST,
			uri:        "/help.txt",
			headersLen: 3,
			body:       "Please visit www.example.re for the latest updates!\nAnother cool body. Hopefully this works",
			method:     server.HTTP_POST,
		},
	}
	for _, mr := range mockRequests {
		req, err := server.ExtractRequest(mr.mock)
		if err != nil {
			t.Fatal(err)
		}

		checkMethod := req.GetMethod() == mr.method
		if !checkMethod {
			t.Errorf("Incorrect method: %s", req.GetMethod())
		}
		checkUri := req.GetUri() == mr.uri
		if !checkUri {
			t.Errorf("Incorrect uri: %s", req.GetUri())
		}
		checkHeaders := len(req.Headers) == mr.headersLen
		if !checkHeaders {
			t.Errorf("Incorrect len of headers: %d", len(req.Headers))
		}
		checkBody := req.GetBody() == mr.body
		if !checkBody {
			t.Errorf("Incorrect body: %s", req.GetBody())
		}
	}
}
