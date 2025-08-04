package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	json, _ := json.Marshal(body)

	url := "https://api.spotify.com/v1/me/player/play"

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", " Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")
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
	registerCommand(StartCommand{})
}
