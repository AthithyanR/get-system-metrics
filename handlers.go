package main

import (
	"encoding/json"
	"net/http"
)

type HealthCheckResponse struct {
	Success bool `json:"success"`
}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users []User

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	response := &HealthCheckResponse{
		Success: true,
	}
	json.NewEncoder(w).Encode(response)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	decoder := json.NewDecoder(r.Body)
	var body User
	err := decoder.Decode(&body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users = append(users, body)
	response := &HealthCheckResponse{
		Success: true,
	}
	json.NewEncoder(w).Encode(response)
}
