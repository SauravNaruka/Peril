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

type config struct {
	ch *amqp.Channel
}

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

	cfg := &config{
		ch: publishCh,
	}
	gamelogic.PrintServerHelp()

	cfg.repl()

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}

func (cfg *config) repl() {
	for {
		words := gamelogic.GetInput()

		switch words[0] {
		case "pause":
			fmt.Println("Publishing paused game state")
			err := pubsub.PublishJSON(cfg.ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		case "resume":
			fmt.Println("Publishing resumes game state")
			err := pubsub.PublishJSON(cfg.ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: false,
			})
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		case "quit":
			log.Println("goodbye")
			return
		default:
			fmt.Println("unknown command")
		}
	}

}
