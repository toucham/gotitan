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
	httpConn, netConn := createMockHttpConn()
	go httpConn.Write()

	// sends input
	ctx := createMockCtx()
	httpConn.queue <- &ctx
	ctx.Ready <- true

	result := make(chan string)
	go readFromNet(netConn, result, len(EXPECTED_RESP_STRING))
	select {
	case got := <-result:
		if EXPECTED_RESP_STRING != got {
			t.Fatal("Not expected string")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout after a second")
	}
}

func TestWriteResponseInClosedChannel(t *testing.T) {
	httpConn, netConn := createMockHttpConn()
	go httpConn.Write()

	// sends input
	ctx := createMockCtx()
	httpConn.queue <- &ctx
	ctx.Ready <- true
	close(httpConn.queue)

	result := make(chan string)
	readFromNet(netConn, result, len(EXPECTED_RESP_STRING)+1)
	got := <-result
	if EXPECTED_RESP_STRING != got {
		t.Fatal("Not expected string")
	}
}
