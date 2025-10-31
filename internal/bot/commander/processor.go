package commander

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (r *Register) Processor(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if !strings.HasPrefix(m.Content, r.prefix) {
		return
	}

	messenger := NewMessenger(s, m.Message)

	content, _ := strings.CutPrefix(m.Content, r.prefix)

	invalidCommand, _ := regexp.Match("^(\\s+).*?$", []byte(content))

	if invalidCommand {
		messenger.Reply("Invalid usage, try `%v<command>`. Type `%vhelp` or `/help` to list all commands.", r.prefix, r.prefix)
		return
	}

	splitContent := strings.Split(content, " ")

	cmdID := splitContent[0]
	args := splitContent[1:]

	ctx := NewContext(r, s, m, args)

	ctx.SetMessenger(messenger)
	ctx.SetAuthorID(m.Author.ID)
	ctx.SetChannelID(m.ChannelID)
	ctx.SetGuildID(m.GuildID)

	cmd, ok := r.Command(cmdID)

	if !ok {
		messenger.Reply("Invalid command: `%v`. Type `%v help` or `/help` to list all commands.", cmdID, r.prefix)
		return
	}

	if cmd == nil {
		messenger.Reply("This command hasn't a handler yet.")
		return
	}

	cmd.Handler(ctx)
}
