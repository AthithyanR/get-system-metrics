package main

import (
	"encoding/json"
	"net/http"
)

type HealthCheckResponse struct {
	Success bool `json:"success"`
}

type User struct {
	Id   int    `json:"success"`
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
