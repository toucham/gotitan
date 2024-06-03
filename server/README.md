# HTTP WebServer v0.1

Implementing a simple HTTP WebServer in Go using only the `net` package that supports HTTP/1.1 according [RFC9110](https://www.rfc-editor.org/rfc/rfc9110.html) and [RFC9112](https://datatracker.ietf.org/doc/html/rfc9112).

## Motivation

This is for educational purposes to better understand Golang and its set of tools. The is project is completed after the server is able handle the following requests using `curl`:

```bash
# GET
curl 127.0.0.1:8080

# POST

# PUT

# DELETE

# Concurrent Requests

```

## Server Features

These are the features that webserver is offering:

- support HTTP/1.1
- easy routing
- able to concurrently process requests
- connection management (persistent connection, timeout, pipelining)

### HTTP/1.1

Implemented as direted in [RFC9112](https://www.rfc-editor.org/rfc/rfc9112.html#name-connection-management).

#### Persistent Connection

HTTP/1.1 allows persistent connection. This means that client can reuse the TCP connection to send separate HTTP messages to the server. This decreases the time it takes to send each HTTP message as TCP connection must be established between the client and server through multiple handshakes.

#### Pipeline Requests

HTTP/1.1 states the following regarding pipeline:

```text
A client that supports persistent connections MAY "pipeline" its requests (i.e., send multiple requests without waiting for each response). A server MAY process a sequence of pipelined requests in parallel if they all have safe methods (Section 9.2.1 of [HTTP]), but it MUST send the corresponding responses in the same order that the requests were received.
```

A safe method is GET, HEAD, OPTIONS, and TRACE methods (as defined in [RFC9110](https://www.rfc-editor.org/rfc/rfc9110#name-safe-methods)).

The RFC states that "a server MAY process a sequence of pipelined requests in parallel if they all have safe method". Therefore, if an unsafe method is received, the webserver will immediately stop processing requests using goroutine.

#### Closing Connection

To close a connection, it must be done in stage so that there's no problem with 

## Architecture

This is the high-level architecture of the web server:

### Message Lifecycle

The following demonstates the HTTP message lifecycle from request into response.