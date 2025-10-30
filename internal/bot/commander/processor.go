package commander

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Processor(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	content, _ := strings.CutPrefix(m.Content, commanderPrefix)
	args := strings.Split(strings.TrimSpace(content), " ")

	cmdID := args[0]
	args = args[1:]

	msn := NewMessenger(s, m.Message)
	ctx := NewContext(s, msn, args)

	ctx.SetAuthorID(msn.rootMessage.Author.ID)
	ctx.SetChannelID(msn.rootMessage.ChannelID)
	ctx.SetGuildID(msn.rootMessage.GuildID)

	cmd, ok := command(cmdID)

	if !strings.HasPrefix(m.Content, commanderPrefix) {
		return
	}

	if !ok {
		msn.Reply("Invalid command: `%v`. Type `%v help` or `/help` to list all commands.", cmdID, commanderPrefix)
		return
	}

	if cmd == nil {
		msn.Reply("This command hasn't a handler yet.")
		return
	}

	cmd.Handler(ctx)
}
