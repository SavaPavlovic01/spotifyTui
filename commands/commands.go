package commands

import (
	"fmt"

	"example.com/m/v2/auth"
)

type Command interface {
	Name() string
	Description() string
	Exec(tokens *auth.FreshToken, args []string) error
}

type WrappedCommand struct {
	command Command
}

func withAuth(command Command) *WrappedCommand {
	return &WrappedCommand{command: command}
}

func (wc *WrappedCommand) Exec(tokens *auth.FreshToken, args []string) error {
	if tokens.Expired() {
		tokens.RefreshTokens()
	}
	return wc.command.Exec(tokens, args)
}

func (wc *WrappedCommand) Name() string {
	return wc.command.Name()
}

func (wc *WrappedCommand) Description() string {
	return wc.command.Description()
}

var CommandRegistry = make(map[string]Command)

func registerCommand(command Command) {
	_, ok := CommandRegistry[command.Name()]
	if ok {
		fmt.Println("Command with the name " + command.Name() + " alredy exists")
	}
	CommandRegistry[command.Name()] = withAuth(command) // does everything need auth?
}
