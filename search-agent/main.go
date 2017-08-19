package main

import (
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"search-agent/model"
)

func main() {
	log.Println("Connecting to RabbitMQ...")
	model.InitRabbitConnection()
	log.Println("Connected")

	r := mux.NewRouter()
	r.HandleFunc("/send/{url}", model.SendMessageHandler).Methods("GET")
	r.HandleFunc("/receive", model.ReceiveMessageHandler).Methods("GET")

	log.Println("Running the server on 8000...")
	http.ListenAndServe(":8000", r)
}