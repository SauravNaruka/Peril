package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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
		log.Fatal("Can not read AMQP_URL from env file")
	}

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	err = pubsub.PublishJSON(publishCh, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})
	if err != nil {
		log.Printf("could not publish time: %v", err)
	}
	fmt.Println("Pause message sent!")

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
