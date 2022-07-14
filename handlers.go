package main

import (
	"encoding/json"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var DB *gorm.DB

type HealthCheckResponse struct {
	Success bool `json:"success"`
}

type User struct {
	Uid  string `json:"uid" gorm:"primaryKey"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// var users []User

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	response := &HealthCheckResponse{
		Success: true,
	}
	json.NewEncoder(w).Encode(response)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var users []User
	DB.Find(&users)
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

	// users = append(users, body)
	body.Uid = uuid.NewV4().String()
	DB.Create(&body)
	response := &HealthCheckResponse{
		Success: true,
	}
	json.NewEncoder(w).Encode(response)
}
