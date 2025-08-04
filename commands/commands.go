package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"example.com/m/v2/auth"
)

type Command interface {
	Name() string
	Description() string
	Exec(tokens *auth.FreshToken, args []string) error
}

type WrappedCommand struct {
	command Command
}

func withAuth(command Command) *WrappedCommand {
	return &WrappedCommand{command: command}
}

func (wc *WrappedCommand) Exec(tokens *auth.FreshToken, args []string) error {
	if tokens.Expired() {
		tokens.RefreshTokens()
	}
	return wc.command.Exec(tokens, args)
}

func (wc *WrappedCommand) Name() string {
	return wc.command.Name()
}

func (wc *WrappedCommand) Description() string {
	return wc.command.Description()
}

var CommandRegistry = make(map[string]Command)

func registerCommand(command Command) {
	_, ok := CommandRegistry[command.Name()]
	if ok {
		fmt.Println("Command with the name " + command.Name() + " alredy exists")
	}
	CommandRegistry[command.Name()] = withAuth(command) // does everything need auth?
}

type DeviceInfo struct {
	Id               string  `json:"id"`
	IsActive         bool    `json:"is_active"`
	IsPrivateSession bool    `json:"is_private_session"`
	IsRestrited      bool    `json:"is_resticted"`
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	VolumePercent    float64 `json:"volume_percent"`
	SupportsVolume   bool    `json:"supports_volume"`
}

func GetDevices(token *auth.FreshToken) ([]DeviceInfo, error) {

	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me/player/devices", http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", " Bearer "+token.AccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 204 {
		// TODO: Add parseError func or something like that
		fmt.Println("ERROR on change playback")
		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
	var info struct {
		List []DeviceInfo `json:"devices"`
	}
	data, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(data, &info)
	defer resp.Body.Close()
	return info.List, nil
}
