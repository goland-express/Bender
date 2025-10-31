package bot

import (
	"bender/internal/bot/commander"
	"bender/internal/bot/minigames"
	"bender/internal/bot/player"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func Start() error {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	if err != nil {
		fmt.Println("Error creating discord session,", err)
		return err
	}

	defaultRegister := commander.NewRegister(".b")

	defaultRegister.AddCommand("play", "It plays a song.", player.CommandPlay)
	defaultRegister.AddCommand("skip", "It skips the current song.", player.CommandSkip)
	defaultRegister.AddCommand("stop", "It stop all the songs.", player.CommandStop)

	defaultRegister.AddCommand("dice", "Roll a dice.", minigames.CommandDice)

	dg.AddHandler(defaultRegister.Processor)
	dg.AddHandler(guildCreateHandler)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	err = dg.Open()
	defer dg.Close()

	if err != nil {
		fmt.Println("Error opening connection,", err)
		return err
	}

	fmt.Println("Bender is working now")

	return err
}
