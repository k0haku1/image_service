package main

import (
	"log"
	"net/http"
)

func main() {
	router := InitRouter()

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
