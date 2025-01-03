package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/SauravNaruka/Peril/internal/gamelogic"
	"github.com/SauravNaruka/Peril/internal/pubsub"
	"github.com/SauravNaruka/Peril/internal/routing"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Can not read env file")
	}

	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		log.Fatal("Peril game client can not read AMQP_URL from env file")
	}

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Peril game client could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game client connected to RabbitMQ!")

	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Peril game client couldn't get username: %v", err)
	}

	pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+userName, routing.PauseKey, pubsub.SimpleQueueTransient)

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
