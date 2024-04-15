# HTTP Server

Implementing a simple HTTP WebServer in Go using only the `net` package that supports HTTP/1.1 according [RFC9110](https://www.rfc-editor.org/rfc/rfc9110.html) and [RFC9112](https://datatracker.ietf.org/doc/html/rfc9112).

This is for educational purposes to get better understanding of Go programming language.

## Features

These are the features that webserver will offer:

- support for HTTP/1.1
- able to concurrently process requests
- apply middleware on both requests and responses

## Architecture

This is the high-level architecture of the web server:

## Learn

### IO in Golang

`io.Writer` and `io.Reader` are interfaces that is wrapped around a file descriptor for writing/reading bytes. The `net.Conn` implements both interfaces, allowing to read/write from/to the socket file descriptor.
