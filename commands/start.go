package commands

import (
	"net/http"

	"example.com/m/v2/auth"
)

type StartCommand struct{}

func (sc StartCommand) Name() string { return "s" }

func (sc StartCommand) Description() string { return "Pause playback" }

func (sc StartCommand) Exec(token *auth.FreshToken, args []string) error {

	body := map[string]int{
		"position_ms": 0,
	}

	url := "https://api.spotify.com/v1/me/player/play"

	resp, err := NewSpotRequest(http.MethodPut, url).WithAuth(token).WithJson(body).Do()
	return ValidateResponse(resp, err)
}

func init() {
	registerCommand(StartCommand{})
}
