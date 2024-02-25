package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	st := fmt.Sprintf("Hi there, I love %s!", r.URL.Path[1:])
	var body = make([]byte, 8) // byte array's capacity is 8
	for {
		n, err := r.Body.Read(body) // it will read 8 byte
		fmt.Printf("n = %v err = %v b = %v\n", n, err, body)
		if err == io.EOF {
			break
		}
	}
	fmt.Fprint(w, st)
}
