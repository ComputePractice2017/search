package model

import (
	"os"
	"fmt"
	"log"
	"net/http"

	"github.com/streadway/amqp"
	"github.com/gorilla/mux"
)

var conn *amqp.Connection
var channel *amqp.Channel
var queue amqp.Queue

func InitRabbitConnection() {
	rabbitMqAddress := os.Getenv("RABBITMQ_HOST")
	if rabbitMqAddress == "" {
		rabbitMqAddress = "localhost"
	}
	log.Printf("RABBITMQ_HOST: %s\n", rabbitMqAddress)

	rabbitMqUser := os.Getenv("RABBITMQ_DEFAULT_USER")
	if rabbitMqUser == "" {
		rabbitMqUser = "guest"
	}
	log.Printf("RABBITMQ_DEFAULT_USER: %s\n", rabbitMqUser)

	rabbitMqPassword := os.Getenv("RABBITMQ_DEFAULT_PASS")
	if rabbitMqPassword == "" {
		rabbitMqPassword = "guest"
	}
	log.Printf("RABBITMQ_DEFAULT_PASS: %s\n", rabbitMqPassword)

	rabbitMqChannel := os.Getenv("RABBITMQ_CHANNEL")
	if rabbitMqChannel == "" {
		rabbitMqChannel = "main"
	}
	log.Printf("RABBITMQ_CHANNEL: %s\n", rabbitMqChannel)

	connectionStringTemplate := "amqp://%s:%s@%s:5672/"
	connectionString := fmt.Sprintf(connectionStringTemplate, rabbitMqUser, rabbitMqPassword, rabbitMqAddress)

	var err error
	conn, err = amqp.Dial(connectionString)
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err = conn.Channel()
	failOnError(err, "Failed to open a send channel")

	queue, err = channel.QueueDeclare(
		rabbitMqChannel,
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(vars)

	err := channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(vars["url"]),
		})

	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ReceiveMessageHandler(w http.ResponseWriter, r *http.Request) {
	msgs, err := channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for d := range msgs {
		w.Write(d.Body)
		d.Ack(false)
		return
	}
}