package conn

import (
	"bufio"
	"testing"
)

const MOCK_PERSISTENT_PIPELINE_REQUEST = `GET /index.html HTTP/1.1
Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)

GET /index.html HTTP/1.1
Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)

`
const MOCK_PERSISTENT_PIPELINE_REQUEST_WITH_POST = `GET /index.html HTTP/1.1
Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)

GET /index.html HTTP/1.1
Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)

POST /help.txt HTTP/1.1
Host: www.example.re
Content-Type: text/plain
Content-Length: 91
Connection: close

Please visit www.example.re for the latest updates!
Another cool body. Hopefully this works`

func TestReadPersistentConnection(t *testing.T) {
	mock, input := createMockHttpConn()

	mockReqs := []MockResult{
		{
			MOCK_PERSISTENT_PIPELINE_REQUEST,
			2,
		},
		{
			MOCK_PERSISTENT_PIPELINE_REQUEST_WITH_POST,
			3,
		},
	}

	for _, mockReq := range mockReqs {
		writer := bufio.NewWriter(input)
		if _, err := writer.WriteString(mockReq.Mock); err != nil {
			t.Fatal("test failed; unable to write")
		} else {
			go func() {
				writer.Flush()
			}()
		}
		go mock.Read()

		// wait for Read() to send to Route.To()
		for i := 0; i < mockReq.NumOfReqs; i++ {
			result, ok := <-mock.queue
			if result != nil && ok {
				if _, ok := <-result.Ready; !ok {
					t.Fatal("channel to know if result is ready is closed")
				}
			} else {
				t.Fatal("result is returned as nil or channel is closed before")
			}
		}
	}
}

const MOCK_GET_REQUEST = `GET /index.html HTTP/1.1
Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)
Connection: close

`

const MOCK_POST_REQUEST = `POST /help.txt HTTP/1.1
Host: www.example.re
Content-Type: text/plain
Content-Length: 91
Connection: close

Please visit www.example.re for the latest updates!
Another cool body. Hopefully this works`

// Expect:
// - only process the first request => return only one result
const MOCK_CLOSE_PIPELINE_REQUEST = `GET /index.html HTTP/1.1
Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)
Connection: close

GET /index.html HTTP/1.1
Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)

`

func TestReadCloseConnection(t *testing.T) {
	mockReqs := []MockResult{
		{
			MOCK_CLOSE_PIPELINE_REQUEST,
			2,
		},
		{
			MOCK_POST_REQUEST,
			1,
		},
		{
			MOCK_GET_REQUEST,
			1,
		},
	}

	for _, req := range mockReqs {
		mock, input := createMockHttpConn()
		writer := bufio.NewWriter(input)
		if _, err := writer.WriteString(req.Mock); err != nil {
			t.Fatal("test setup failed; unable to write")
		} else {
			go func() {
				writer.Flush()
				mock.conn.Close()
			}()
		}
		go mock.Read()

		// wait for Read() to send to Route.To()
		for i := 0; i < req.NumOfReqs; i++ {
			result, ok := <-mock.queue
			if result != nil && ok {
				if _, ok := <-result.Ready; !ok {
					t.Fatal("channel to know if result is ready is closed")
				}
			} else {
				t.Fatal("result is returned as nil or channel is closed before")
			}
		}
	}
}

// TEST: Discard incorrect requests

// incorrect first request (no newline)
// expect: close connection
const MOCK_INCORRECT_FORMAT_REQ = `Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)
GET /index.html HTTP/1.1
Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)

`

func TestReadIncorrectCloseConn(t *testing.T) {
	mock, input := createMockHttpConn()
	writer := bufio.NewWriter(input)
	if _, err := writer.WriteString(MOCK_INCORRECT_FORMAT_REQ); err != nil {
		t.Fatal("test failed; unable to write")
	} else {
		go func() {
			writer.Flush()
		}()
	}
	go mock.Read()

	// assert that the channel is closed without sending anything to ch
	_, ok := <-mock.queue
	if ok {
		t.Fatal("A result is sent to channel when none is expected")
	}
}

// expect: Send all requests to Route.To()
const MOCK_UNSAFE_METHOD_PIPELINE = `GET /index.html HTTP/1.1
Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)

GET /index.html HTTP/1.1
Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)

POST /help.txt HTTP/1.1
Host: www.example.re
Content-Type: text/plain
Content-Length: 92

Please visit www.example.re for the latest updates!
Another cool body. Hopefully this works
GET /index.html HTTP/1.1
Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)

POST /help.txt HTTP/1.1
Host: www.example.re
Content-Type: text/plain
Content-Length: 6

Please`

func TestReadUnsafeMethod(t *testing.T) {
	mock, input := createMockHttpConn()
	writer := bufio.NewWriter(input)
	if _, err := writer.WriteString(MOCK_UNSAFE_METHOD_PIPELINE); err != nil {
		t.Fatal("test failed; unable to write")
	} else {
		go func() {
			writer.Flush()
			mock.conn.Close()
		}()
	}
	go mock.Read()

	count := 1
	for ; ; count++ {
		ctx, ok := <-mock.queue
		if !ok {
			break
		}
		<-ctx.Ready
	}

	if mock.isSafePipeline {
		t.Errorf("Should be considered as safe method")
	}
	if count == 4 {
		t.Errorf("There should be 4 requests send to Route.To()")
	}
}
