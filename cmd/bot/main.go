package main

import (
	"bender/internal/player"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}
}

func main() {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	if err != nil {
		fmt.Println("Error creating discord session,", err)
		return
	}

	dg.AddHandler(messageHandler)
	dg.AddHandler(guildCreateHandler)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	err = dg.Open()
	defer dg.Close()

	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	fmt.Println("Bender is working now")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func guildCreateHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	player.Init(s, g.ID)
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if !strings.HasPrefix(m.Content, "bender") {
		return
	}

	content, _ := strings.CutPrefix(m.Content, "bender")

	args := strings.Split(strings.TrimSpace(content), " ")

	command := args[0]
	args = args[1:]

	if command == "stop" {
		player.Stop(m.GuildID)

		return
	}

	if command == "skip" {
		player.Skip(m.GuildID)

		return
	}

	if command == "play" {
		player.Play(s, player.PlayQuery{
			GuildID:   m.GuildID,
			UserID:    m.Author.ID,
			ChannelID: m.ChannelID,
			MsgRef:    m.Reference(),
			Query:     strings.Join(args, " "),
		})

		return
	}

	_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Invalid command `%s`.", command), m.Reference())

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
}
