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

type config struct {
	connection *amqp.Connection
	username   string
	gameState  gamelogic.GameState
}

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

	gs := gamelogic.NewGameState(userName)

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+gs.GetUsername(),
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gs),
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

	cfg := &config{
		connection: conn,
		username:   userName,
		gameState:  *gs,
	}

	cfg.repl()

	// // wait for ctrl+c
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan
}

func (cfg *config) repl() {
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "spawn":
			err := cfg.gameState.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "move":
			_, err := cfg.gameState.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "status":
			cfg.gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
		default:
			fmt.Println("unknown command")
		}
	}

}
