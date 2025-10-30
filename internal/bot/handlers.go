package bot

import (
	"bender/internal/bot/player"

	"github.com/bwmarrin/discordgo"
)

// if command == "quote" {
// 	if s.State.User.ID == m.ReferencedMessage.Author.ID {
// 		_, err := s.ChannelMessageSendReply(m.ChannelID, "You can't quote me! I'm unquotable, dude.", m.Reference())
//
// 		if err != nil {
// 			log.Printf("Error: %v", err)
// 			return
// 		}
//
// 		return
// 	}
//
// 	if m.ReferencedMessage == nil {
// 		_, err := s.ChannelMessageSendReply(m.ChannelID, "You need to reply to a message in order to make a quote.", m.Reference())
//
// 		if err != nil {
// 			log.Printf("Error: %v", err)
// 			return
// 		}
// 	}
//
// 	displayName := m.ReferencedMessage.Author.DisplayName()
//
// 	msgContent := fmt.Sprintf("\"%s\" — %s", m.ReferencedMessage.Content, displayName)
//
// 	_, err := s.ChannelMessageSendReply(m.ChannelID, msgContent, m.Reference())
//
// 	if err != nil {
// 		log.Printf("Error: %v", err)
// 		return
// 	}
//
// 	return
// }

func guildCreateHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	player.Init(s, g.ID)
}

