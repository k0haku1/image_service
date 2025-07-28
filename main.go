package main

import (
	"log"
	"net/http"
)

func main() {
	router := InitRouter()

	log.Println("Server listening on :8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatal(err)
	}
}
