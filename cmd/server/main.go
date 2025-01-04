package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SauravNaruka/Peril/internal/gamelogic"
	"github.com/SauravNaruka/Peril/internal/pubsub"
	"github.com/SauravNaruka/Peril/internal/routing"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type EnvConfig struct {
	AMQPURL string
}

type config struct {
	ch *amqp.Channel
}

func main() {
	env, err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Establish connection
	conn, err := amqp.Dial(env.AMQPURL)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	// Create a publish channel to RabbitMQ.
	// Basically A virtual connection inside a connection that allows you to create queues, exchanges, and publish messages.
	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	_, queue, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", pubsub.SimpleQueueDurable)
	if err != nil {
		log.Fatalf("could not create queue: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	cfg := &config{
		ch: publishCh,
	}

	cfg.repl()
}

func (cfg *config) repl() {
	gamelogic.PrintServerHelp()

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

func loadEnv() (EnvConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return EnvConfig{}, fmt.Errorf("error while loading environment file: %v", err)
	}

	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		return EnvConfig{}, fmt.Errorf("error while reading AMQP_URL url: %v", err)
	}

	config := EnvConfig{
		AMQPURL: amqpURL,
	}

	return config, nil
}
