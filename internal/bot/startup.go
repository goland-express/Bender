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
	defaultRegister.WithHelpCommad()

	player.RegisterCommands(defaultRegister)
	minigames.RegisterCommands(defaultRegister)

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
