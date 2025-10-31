package commander

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (r *Register) Processor(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	content, _ := strings.CutPrefix(m.Content, r.prefix)
	args := strings.Split(strings.TrimSpace(content), " ")

	cmdID := args[0]
	args = args[1:]

	ctx := NewContext(r, s, m, args)

	ctx.SetAuthorID(m.Author.ID)
	ctx.SetChannelID(m.ChannelID)
	ctx.SetGuildID(m.GuildID)

	cmd, ok := r.Command(cmdID)

	if !strings.HasPrefix(m.Content, r.prefix) {
		return
	}

	if !ok {
		ctx.Messenger.Reply("Invalid command: `%v`. Type `%v help` or `/help` to list all commands.", cmdID, r.prefix)
		return
	}

	if cmd == nil {
		ctx.Messenger.Reply("This command hasn't a handler yet.")
		return
	}

	cmd.Handler(ctx)
}
