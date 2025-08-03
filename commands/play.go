package commands

import (
	"fmt"
	"io"
	"net/http"

	"example.com/m/v2/auth"
)

type PlayCommand struct{}

func (pc PlayCommand) Name() string { return "p" }

func (pc PlayCommand) Description() string { return "Start playback" }

func (pc PlayCommand) Exec(token *auth.FreshToken, args []string) error {
	req, err := http.NewRequest(http.MethodPut, "https://api.spotify.com/v1/me/player/pause", http.NoBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", " Bearer "+token.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode > 204 {
		// TODO: Add parseError func or something like that
		fmt.Println("ERROR on change playback")
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
	return nil
}

func init() {
	registerCommand(PlayCommand{})
}
