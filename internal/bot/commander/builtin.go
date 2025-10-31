package commander

import "fmt"

func (r *Register) WithHelpCommad() {
	r.AddCommand("help", "It helps you listing all commands", BuiltInHelpCommand)
}

func BuiltInHelpCommand(ctx *Context) {
	commands := ctx.Register.Commands()

	var helpContent string

	for command := range commands {
		helpContent += fmt.Sprintf("%s %s - %s\n", ctx.Register.Prefix(), command.Identifier, command.Description)
	}

	ctx.Messenger.Reply(helpContent)
}
