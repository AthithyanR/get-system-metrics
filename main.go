package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func getPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return port
}

func initDB() {}

func initHandlers(r *mux.Router) {

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./fe/dist")))

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/health-check", HandleHealthCheck)
	apiRouter.HandleFunc("/users", HandleHealthCheck).Methods("GET")
	apiRouter.HandleFunc("/users", HandleHealthCheck).Methods("POST")
}

func initServer(serverPort string, r *mux.Router) {
	log.Println("Starting server on port: " + serverPort)
	if err := http.ListenAndServe(":"+serverPort, r); err != nil {
		log.Println("Error starting server on port: " + serverPort)
	}
}

func main() {
	serverPort := getPort()
	r := mux.NewRouter()
	initHandlers(r)
	initServer(serverPort, r)
}
