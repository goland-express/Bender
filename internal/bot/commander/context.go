package commander

import "github.com/bwmarrin/discordgo"

type Context struct {
	Session   *discordgo.Session
	Messenger Messenger
	Args      []string
}

func NewContext(session *discordgo.Session, messenger *Messenger, args []string) *Context {
	return &Context{
		Session:   session,
		Messenger: *messenger,
		Args:      args,
	}
}
