package main

import (
	"encoding/json"
	"net/http"
)

type api struct {
	address string
}

var users = []User{
	{
		Id:   "1",
		Age:  22,
		Name: "Manuel",
	},
}

func (api *api) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (api *api) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload User

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users = append(users, payload)

	w.WriteHeader(http.StatusOK)
}
