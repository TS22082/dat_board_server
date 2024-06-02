package main

import (
	"fmt"
	"net/http"

	"github.com/TS22082/dat_board_server/api/test"
)

func main() {
	http.HandleFunc("/", test.HelloHandler)
	http.HandleFunc("/2", test.HelloHandler2)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Server failed to start: %v", err)
	}
}
