package minigames

import (
	"math/rand"
)

func rollDice(sides int) int {
	return rand.Intn(sides-1) + 1
}
