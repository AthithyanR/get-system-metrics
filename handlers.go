package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var DB *gorm.DB

type HealthCheckResponse struct {
	Success bool `json:"success"`
}

type User struct {
	Uid  string `json:"uid" gorm:"primaryKey"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Stats struct {
	MemoryTotal  uint64  `json:"MemoryTotal"`
	MemoryUsed   uint64  `json:"MemoryUsed"`
	MemoryCached uint64  `json:"MemoryCached"`
	MemoryFree   uint64  `json:"MemoryFree"`
	CpuUser      float64 `json:"CpuUser"`
	CpuSystem    float64 `json:"CpuSystem"`
	CpuIdle      float64 `json:"CpuIdle"`
}

var stat Stats

var wsClients []*websocket.Conn

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
	// DB.Find(&users)
	params := mux.Vars(r)
	count, err := strconv.Atoi(params["count"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for i := 0; i < count; i++ {
		users = append(users, User{
			Uid:  uuid.NewV4().String(),
			Name: "Athi " + strconv.Itoa(i+1),
			Age:  23,
		})
	}
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

func HandleWS(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	wsClients = append(wsClients, conn)
}

func SendStat() {
	for _, conn := range wsClients {
		conn.WriteJSON(stat)
	}
}
