package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received webhook:")
		fmt.Println(r.Body)
		w.WriteHeader(200)
		w.Write([]byte("hello"))
	})

	fmt.Println("Started http server...")
	http.ListenAndServe(":3000", nil)
}
