package commander

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Processor(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if !strings.HasPrefix(m.Content, commanderPrefix) {
		fmt.Println("Hasn't prefix")
		return
	}

	content, _ := strings.CutPrefix(m.Content, commanderPrefix)
	args := strings.Split(strings.TrimSpace(content), commanderPrefix)

	cmdID := args[0]
	args = args[1:]

	cmd, ok := command(cmdID)

	msn := NewMessenger(s, m.Message)
	ctx := NewContext(s, msn, args)

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
