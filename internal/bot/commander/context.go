package commander

import "github.com/bwmarrin/discordgo"

type Context struct {
	Session   *discordgo.Session
	GuildID   string
	AuthorID  string
	ChannelID string
	Messenger *Messenger
	Args      []string
}

func NewContext(session *discordgo.Session, messenger *Messenger, args []string) *Context {
	return &Context{
		Session:   session,
		Messenger: messenger,
		Args:      args,
	}
}

func (ctx *Context) SetSession(session *discordgo.Session) {
	ctx.Session = session
}

func (ctx *Context) SetGuildID(guildID string) {
	ctx.GuildID = guildID
}

func (ctx *Context) SetAuthorID(authorID string) {
	ctx.AuthorID = authorID
}

func (ctx *Context) SetChannelID(channelID string) {
	ctx.ChannelID = channelID
}

func (ctx *Context) SetMessenger(messenger *Messenger) {
	ctx.Messenger = messenger
}
