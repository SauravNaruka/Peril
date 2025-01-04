package main

import (
	"fmt"

	"github.com/SauravNaruka/Peril/internal/gamelogic"
	"github.com/SauravNaruka/Peril/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")

		gs.HandlePause(ps)
	}
}
