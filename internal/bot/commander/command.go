package commander

type CommandHandler func(ctx *Context)

type Command struct {
	Identifier  string
	Description string
	Handler     CommandHandler
}

func NewCommand(identifier, description string, handler CommandHandler) *Command {
	return &Command{
		Identifier:  identifier,
		Description: description,
		Handler:     handler,
	}
}
