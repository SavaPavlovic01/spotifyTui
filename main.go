package main

import (
	"os"
	"strings"

	"example.com/m/v2/auth"
	"example.com/m/v2/commands"
)

func ParseInput() []string {
	var query strings.Builder
	ret := []string{}
	i := 2
	for i < len(os.Args) {
		if !strings.HasPrefix(os.Args[i], "-") {
			break
		}
		ret = append(ret, os.Args[i])
		i += 1
	}

	for i < len(os.Args) {
		query.WriteString(os.Args[i])
		query.WriteString(" ")
		i += 1
	}

	return append([]string{query.String()}, ret...)
}

func main() {
	tokens, err := auth.LoadTokens()
	if err != nil {
		tokens = auth.InitialAuth()
	}

	input := ParseInput()
	if len(os.Args) == 1 {
		panic("Pls send a command")
	}

	err = commands.CommandRegistry[os.Args[1]].Exec(tokens, input)
	if err != nil {
		panic(err)
	}
}
