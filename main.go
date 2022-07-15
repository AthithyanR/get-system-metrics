package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "80"
	}
	return port
}

func Init() {
	dbURL := os.Getenv("dbURL")

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&User{})

	DB = db

}

func initHandlers(r *mux.Router) {

	var dir string

	flag.StringVar(&dir, "dir", "./fe/dist", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	// r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/health-check", HandleHealthCheck)
	apiRouter.HandleFunc("/users/{count}", GetUsers).Methods("OPTIONS", "GET")
	// apiRouter.HandleFunc("/users", CreateUser).Methods("OPTIONS", "POST")
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
	// Init()
	initHandlers(r)
	initServer(serverPort, r)
}
