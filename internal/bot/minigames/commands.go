package minigames

import (
	"bender/internal/bot/commander"
	"strconv"
)

func CommandDice(ctx *commander.Context) {
	sides := 6
	result := rollDice(sides)

	if len(ctx.Args) > 0 {
		if s, err := strconv.Atoi(ctx.Args[0]); err == nil {
			sides = s
		} else {
			ctx.Messenger.Reply("Misuse of command. Use %s dice <amount of sides>", commander.Prefix())
			return
		}
	}

	ctx.Messenger.Reply("Here is the result: **%d**.", result)
}
