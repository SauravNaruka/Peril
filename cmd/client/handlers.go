package main

import (
	"fmt"

	"github.com/SauravNaruka/Peril/internal/gamelogic"
	"github.com/SauravNaruka/Peril/internal/pubsub"
	"github.com/SauravNaruka/Peril/internal/routing"
)

func handlerMove(gs *gamelogic.GameState) func(move gamelogic.ArmyMove) pubsub.Acktype {
	return func(move gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")

		moveOutcome := gs.HandleMove(move)
		switch moveOutcome {
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			return pubsub.Ack
		}

		fmt.Println("error: unknown move outcome")
		return pubsub.NackDiscard
	}
}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.Acktype {
	return func(ps routing.PlayingState) pubsub.Acktype {
		defer fmt.Print("> ")
		gs.HandlePause(ps)

		return pubsub.Ack
	}
}
