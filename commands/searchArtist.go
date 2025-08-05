package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"example.com/m/v2/auth"
)

type searchArtistCommand struct{}

func (sc searchArtistCommand) Name() string { return "sat" }

func (sc searchArtistCommand) Description() string { return "" }

func (sc searchArtistCommand) Exec(token *auth.FreshToken, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("pls send me params")
	}

	artistInfo, err := SearchArtist(10, 0, args[0], token)
	if err != nil {
		return err
	}
	url := "https://api.spotify.com/v1/artists/" + getIdFromUri(artistInfo[0].Uri) + "/top-tracks"
	resp, err := NewSpotRequest(http.MethodGet, url).
		WithAuth(token).Do()

	err = ValidateResponse(resp, err)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var ResponseData struct {
		Info []TrackInfo `json:"tracks"`
	}

	err = json.Unmarshal(data, &ResponseData)

	if err != nil {
		return err
	}

	return PlayTrack(token, ResponseData.Info[0].Album.Uri, ResponseData.Info[0].Uri)

}

func init() {
	registerCommand(searchArtistCommand{})
}
