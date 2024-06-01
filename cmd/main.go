package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func main() {
	http.HandleFunc("/", helloHandler)

	err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v", err)
	}
}
