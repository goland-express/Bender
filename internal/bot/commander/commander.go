package commander

import (
	"errors"
)

const (
	errCommandIndentifierAlreadyExists = "command identifier already exists"
	errInvalidCommandIdentifier        = "command hasn't an identifier"
)

var commanderPrefix = "bender"
var commands = make(map[string]*Command)

func AddCommand(identifier, description string, handler CommandHandler) (*Command, error) {
	if identifier == "" {
		return nil, errors.New(errInvalidCommandIdentifier)
	}

	_, commandExists := commands[identifier]

	if commandExists {
		return nil, errors.New(errCommandIndentifierAlreadyExists)
	}

	cmd := NewCommand(identifier, description, handler)

	commands[identifier] = cmd

	return cmd, nil
}

func SetPrefix(prefix string) {
	commanderPrefix = prefix
}

func command(identifier string) (*Command, bool) {
	command, ok := commands[identifier]

	if !ok {
		return nil, ok
	}

	return command, true
}
