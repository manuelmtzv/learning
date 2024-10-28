package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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

	u := User{
		Id:   strconv.Itoa(len(users) + 1),
		Age:  payload.Age,
		Name: payload.Name,
	}

	err = insertUser(u)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func insertUser(u User) error {
	if u.Name == "" {
		return errors.New("name is required")
	}

	if u.Username == "" {
		return errors.New("username is required")
	}

	if u.Age == 0 || u.Age > 100 {
		return errors.New("age should be between 0 and 100")
	}

	for _, user := range users {
		if user.Username == u.Username {
			return errors.New("this user already exists")
		}
	}

	users = append(users, u)

	return nil
}
