package commander

import (
	"errors"
)

const (
	errCommandIndentifierAlreadyExists = "command identifier already exists"
	errInvalidCommandIdentifier        = "command hasn't an identifier"
)

type Register struct {
	prefix   string
	commands map[string]*Command
}

func NewRegister(prefix string) *Register {
	return &Register{
		prefix:   prefix,
		commands: make(map[string]*Command),
	}
}

func (r *Register) AddCommand(identifier, description string, handler CommandHandler) (*Command, error) {
	if identifier == "" {
		return nil, errors.New(errInvalidCommandIdentifier)
	}

	_, commandExists := r.commands[identifier]

	if commandExists {
		return nil, errors.New(errCommandIndentifierAlreadyExists)
	}

	cmd := NewCommand(identifier, description, handler)

	r.commands[identifier] = cmd

	return cmd, nil
}

func (r *Register) SetPrefix(prefix string) {
	r.prefix = prefix
}

func (r *Register) Prefix() string {
	return r.prefix
}

func (r *Register) Command(identifier string) (*Command, bool) {
	command, ok := r.commands[identifier]

	if !ok {
		return nil, ok
	}

	return command, true
}
