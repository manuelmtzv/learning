package main

import (
	"log"
	"net/http"
)

func main() {
	api := &api{
		address: ":8080",
	}

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    api.address,
		Handler: mux,
	}

	mux.HandleFunc("GET /users", api.getUsersHandler)
	mux.HandleFunc("POST /users", api.createUserHandler)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
