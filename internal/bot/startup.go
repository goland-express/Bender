package bot

import (
	"bender/internal/bot/commander"
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

	dg.AddHandler(commander.Processor)
	dg.AddHandler(guildCreateHandler)

	commander.SetPrefix(".b")

	commander.AddCommand("play", "It plays a song.", player.Play)
	commander.AddCommand("skip", "It skips the current song.", player.Skip)
	commander.AddCommand("stop", "It stop all the songs.", player.Stop)

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
