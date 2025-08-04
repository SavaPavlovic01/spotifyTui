package commands

import (
	"net/http"

	"example.com/m/v2/auth"
)

type PreviousCommand struct{}

func (pc PreviousCommand) Name() string { return "b" }

func (pc PreviousCommand) Description() string { return "" }

func (pc PreviousCommand) Exec(token *auth.FreshToken, args []string) error {
	resp, err := NewSpotRequest(http.MethodPost, "https://api.spotify.com/v1/me/player/previous").WithAuth(token).Do()
	return ValidateResponse(resp, err)
}

func init() {
	registerCommand(PreviousCommand{})
}
