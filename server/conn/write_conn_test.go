package conn

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func readFromNet(conn net.Conn, ch chan<- string, len int) {
	scanner := bufio.NewScanner(conn)
	scanner.Split(bufio.ScanBytes)
	text := make([]byte, 0, 100)
	for i := 0; i < len && scanner.Scan(); i++ {
		t := scanner.Bytes()
		text = append(text, t...)
	}
	textStr := string(text)
	go func() {
		ch <- textStr
		close(ch)
	}()
}

func TestWriteResponseInPersistentChannel(t *testing.T) {
	writeTo, readFrom := net.Pipe()
	channel := make(chan *routerContext)
	go write(writeTo, channel, new(MockLogger))

	// sends input
	ctx := createMockCtx()
	channel <- &ctx
	close(ctx.Done)

	result := make(chan string)
	go readFromNet(readFrom, result, len(EXPECTED_RESP_STRING))
	select {
	case got := <-result:
		if EXPECTED_RESP_STRING != got {
			t.Fatal("Not expected string")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout after 2 seconds")
	}
}

func TestWriteResponseInClosedChannel(t *testing.T) {
	writeTo, readFrom := net.Pipe()
	channel := make(chan *routerContext)
	go write(writeTo, channel, new(MockLogger))

	// sends input
	ctx := createMockCtx()
	channel <- &ctx
	close(ctx.Done)
	close(channel)

	result := make(chan string)
	readFromNet(readFrom, result, len(EXPECTED_RESP_STRING)+1)
	got := <-result
	if EXPECTED_RESP_STRING != got {
		t.Fatal("Not expected string")
	}
}
