package bot

import (
	"bender/internal/bot/player"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

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

	if command == "quote" {
		if m.ReferencedMessage == nil {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "You need to reply to a message in order to make a quote.", m.Reference())

			if err != nil {
				log.Printf("Error: %v", err)
				return
			}
		}

		displayName := m.ReferencedMessage.Author.DisplayName()

		msgContent := fmt.Sprintf("\"%s\" — %s", m.ReferencedMessage.Content, displayName)

		_, err := s.ChannelMessageSendReply(m.ChannelID, msgContent, m.Reference())

		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		return
	}

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
