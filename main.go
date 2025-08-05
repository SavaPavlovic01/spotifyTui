package main

import (
	"example.com/m/v2/auth"
	"example.com/m/v2/commands"
)

func main() {
	tokens, err := auth.LoadTokens()
	if err != nil {
		tokens = auth.InitialAuth()
	}
	//_ = commands.CommandRegistry["s"].Exec(tokens, []string{})
	err = commands.CommandRegistry["saa"].Exec(tokens, []string{"dev lemons"})

	if err != nil {
		panic(err)
	}
}
