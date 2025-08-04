package commands

import (
	"net/http"

	"example.com/m/v2/auth"
)

type NextCommand struct{}

func (nc NextCommand) Name() string { return "n" }

func (nc NextCommand) Description() string { return "" }

func (nc NextCommand) Exec(token *auth.FreshToken, args []string) error {
	resp, err := NewSpotRequest(http.MethodPost, "https://api.spotify.com/v1/me/player/next").WithAuth(token).Do()
	return ValidateResponse(resp, err)
}

func init() {
	registerCommand(NextCommand{})
}
