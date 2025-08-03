package main

import (
	"fmt"

	"example.com/m/v2/auth"
	"example.com/m/v2/commands"
)

func main() {
	tokens, err := auth.LoadTokens()
	if err != nil {
		tokens = auth.InitialAuth()
	}
	fmt.Println(tokens)
	_ = commands.CommandRegistry["p"].Exec(tokens, []string{})
}
