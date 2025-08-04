package commands

import (
	"net/http"

	"example.com/m/v2/auth"
)

type PlayCommand struct{}

func (pc PlayCommand) Name() string { return "p" }

func (pc PlayCommand) Description() string { return "Pause playback" }

func (pc PlayCommand) Exec(token *auth.FreshToken, args []string) error {
	resp, err := NewSpotRequest(http.MethodPut, "https://api.spotify.com/v1/me/player/pause").WithAuth(token).Do()
	return ValidateResponse(resp, err)
}

func init() {
	registerCommand(PlayCommand{})
}
