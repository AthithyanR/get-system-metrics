package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
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

func monitor() {
	for {
		memory, err := memory.Get()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
		stat.MemoryTotal = memory.Total / 1048576
		stat.MemoryUsed = memory.Used / 1048576
		stat.MemoryCached = memory.Cached / 1048576
		stat.MemoryFree = memory.Free / 1048576

		before, err := cpu.Get()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
		time.Sleep(time.Duration(1) * time.Second)
		after, err := cpu.Get()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
		total := float64(after.Total - before.Total)
		stat.CpuUser = float64(after.User-before.User) / total * 100
		stat.CpuSystem = float64(after.System-before.System) / total * 100
		stat.CpuIdle = float64(after.Idle-before.Idle) / total * 100

		SendStat()
	}
}

func initHandlers(r *mux.Router) {

	var dir string

	flag.StringVar(&dir, "dir", "./fe/dist", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	r.HandleFunc("/ws", HandleWS)

	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/health-check", HandleHealthCheck)
	apiRouter.HandleFunc("/users/{count}", GetUsers).Methods("OPTIONS", "GET")
	// apiRouter.HandleFunc("/users", CreateUser).Methods("OPTIONS", "POST")

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))
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
	go monitor()
	initHandlers(r)
	initServer(serverPort, r)
}
