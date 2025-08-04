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
	//_ = commands.CommandRegistry["s"].Exec(tokens, []string{})
	err = commands.CommandRegistry["b"].Exec(tokens, []string{})
	if err != nil {
		panic(err)
	}
}
