package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", handleConnection)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
