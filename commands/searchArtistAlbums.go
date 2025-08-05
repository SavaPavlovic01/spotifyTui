package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"example.com/m/v2/auth"
)

type searchArtistAlbumsCommand struct{}

func (sc searchArtistAlbumsCommand) Name() string { return "saa" }

func (sc searchArtistAlbumsCommand) Description() string { return "" }

func (sc searchArtistAlbumsCommand) Exec(token *auth.FreshToken, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("pls send me params")
	}

	artistInfo, err := SearchArtist(10, 0, args[0], token)
	if err != nil {
		return err
	}

	url := "https://api.spotify.com/v1/artists/" + getIdFromUri(artistInfo[0].Uri) + "/albums"
	resp, err := NewSpotRequest(http.MethodGet, url).WithAuth(token).Do()
	err = ValidateResponse(resp, err)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	var Response struct {
		Albums []AlbumInfo `json:"items"`
	}
	err = json.Unmarshal(data, &Response)
	if err != nil {
		return err
	}

	tmp := make([]Playable, len(Response.Albums))
	for i, album := range Response.Albums {
		tmp[i] = album
	}
	return Interactive(token, tmp)
}

func init() {
	registerCommand(searchArtistAlbumsCommand{})
}
