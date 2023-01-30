package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const port string = "80"

type Config struct {
	Rabbit *amqp.Connection
}

type Note struct {
	ID              int       `json:"note_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	TextColor       string    `json:"text_color"`
	BackgroundColor string    `json:"background_color"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func main() {
	// connect to rabbitmq
	rabbitConn, err := connect()

	if err != nil {
		log.Panic(err)
	}

	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Panic("error when starting servser: %s", err)
	}
}

func connect() (*amqp.Connection, error) {
	var attempts int64
	var delay = 1 * time.Second
	var connection *amqp.Connection

	for {
		// create new connection
		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq")

		if err != nil {
			fmt.Println("RabbitMQ not yet ready")
			attempts++
		} else {
			// if connection is ready break the loop
			connection = conn
			break
		}

		// if number of connection attempts is bigger than 5 return error
		if attempts > 5 {
			fmt.Println(err)
			return nil, err
		}

		delay = time.Duration(math.Pow(float64(attempts), 2)) * time.Second
		log.Println("Delay attempt...")

		time.Sleep(delay)
		continue
	}

	return connection, nil
}
