package minigames

import (
	"bender/internal/bot/commander"
	"strconv"
)

func Dice(ctx *commander.Context) {
	sides := 6
	result := rollDice(sides)

	if len(ctx.Args) > 0 {
		strconv.Atoi(ctx.Args[0])
	}

	ctx.Messenger.Reply("Here is the result: **%d**.", result)
}
