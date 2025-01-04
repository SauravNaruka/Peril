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
	connection     *amqp.Connection
	publishChannel *amqp.Channel
	username       string
	gameState      gamelogic.GameState
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

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Peril game client couldn't get username: %v", err)
	}

	gs := gamelogic.NewGameState(userName)

	cfg := &config{
		connection:     conn,
		publishChannel: publishCh,
		username:       userName,
		gameState:      *gs,
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+gs.GetUsername(),
		routing.ArmyMovesPrefix+".*",
		pubsub.SimpleQueueTransient,
		cfg.handlerMove(gs),
	)
	if err != nil {
		log.Fatalf("could not subscribe to army moves: %v", err)
	}

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

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.WarRecognitionsPrefix,
		routing.WarRecognitionsPrefix+".*",
		pubsub.SimpleQueueDurable,
		handlerWar(gs),
	)
	if err != nil {
		log.Fatalf("could not subscribe to war declarations: %v", err)
	}

	cfg.repl()
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
			move, err := cfg.gameState.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
			routingKey := routing.ArmyMovesPrefix + "." + cfg.gameState.GetUsername()
			err = pubsub.PublishJSON(
				cfg.publishChannel,
				routing.ExchangePerilTopic,
				routingKey,
				move,
			)
			if err != nil {
				fmt.Printf("error while publishing changes to exchange=%s key=%s. Error: %v\n", routing.ExchangePerilTopic, routingKey, err)
				continue
			}
			fmt.Printf("Moved %v units to %s\n", len(move.Units), move.ToLocation)
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
