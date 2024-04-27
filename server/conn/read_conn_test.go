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
	mockLogger := MockLogger{T: t}
	mock, input, ch := createMockHttpConn(&mockLogger)

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
			result, ok := <-ch
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
	mockLogger := MockLogger{T: t}

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
		mock, input, ch := createMockHttpConn(&mockLogger)
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
			result, ok := <-ch
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
// expect: discard first and second message
const MOCK_INCORRECT_FORMAT_REQ = `Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)
GET /index.html HTTP/1.1
Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)

`

// should discard request if http request is in incorrect format
func TestReadIncorrectRequestDiscard(t *testing.T) {
	mockLogger := MockLogger{T: t}

	mock, input, ch := createMockHttpConn(&mockLogger)
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
	_, ok := <-ch
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
Content-Length: 91

Please visit www.example.re for the latest updates!
Another cool body. Hopefully this works
GET /index.html HTTP/1.1
Host: www.example.re
User-Agent: Mozilla/5.0 (Windows; U; Windows NT 5.0; en-US; rv:1.1)

`

func TestReadUnsafeMethod(t *testing.T) {
	mockLogger := MockLogger{T: t}

	mock, input, ch := createMockHttpConn(&mockLogger)
	writer := bufio.NewWriter(input)
	if _, err := writer.WriteString(MOCK_UNSAFE_METHOD_PIPELINE); err != nil {
		t.Fatal("test failed; unable to write")
	} else {
		go func() {
			writer.Flush()
		}()
	}
	go mock.Read()

	count := 0
	for {
		_, ok := <-ch
		if !ok {
			break
		}
		count++
	}

	// TODO: is there a way to test it is not use gouritine
	if !mock.isSafeMethod {
		t.Errorf("Still considered as safe method")
	}
	if count == 4 {
		t.Errorf("There should be 4 requests sent to Route.To()")
	}
}
