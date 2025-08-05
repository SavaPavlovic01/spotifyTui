package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

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

type TrackInfo struct {
	Name    string       `json:"name"`
	Id      string       `json:"id"`
	Album   AlbumInfo    `json:"album"`
	Artists []ArtistInfo `json:"artists"`
	Uri     string       `json:"uri"`
}

type AlbumInfo struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

type ArtistInfo struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
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

func PlayTrack(token *auth.FreshToken, albumUri string, trackUri string) error {
	var body struct {
		Context_uri string `json:"context_uri"`
		//Uris []string `json:"uris"` (this plays just the track and deletes the rest of the queue)
		// (this play the song and puts the rest of the album in the queue, this is what happens when you press in ui)
		Offset *struct {
			Uri string `json:"uri"`
		} `json:"offset"`
	}
	body.Context_uri = albumUri
	//body.Uris = []string{trackUri}
	if trackUri != "" {
		body.Offset = &struct {
			Uri string `json:"uri"`
		}{Uri: trackUri}
	}
	resp, err := NewSpotRequest(http.MethodPut, "https://api.spotify.com/v1/me/player/play").WithAuth(token).WithJson(body).Do()
	return ValidateResponse(resp, err)
}

func Search(searchTarget string, limit int, offset int, query string, token *auth.FreshToken) ([]byte, error) {
	resp, err := NewSpotRequest(http.MethodGet, "https://api.spotify.com/v1/search").
		WithQueryParam("type", searchTarget).WithQueryParam("q", query).
		WithQueryParam("limit", strconv.Itoa(limit)).WithQueryParam("offset", strconv.Itoa(offset)).WithAuth(token).Do()
	err = ValidateResponse(resp, err)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func ParseTrackInfo(data []byte) ([]TrackInfo, error) {
	var SearchResponse struct {
		Tracks struct {
			Items []TrackInfo `json:"items"`
		} `json:"tracks"`
	}

	err := json.Unmarshal(data, &SearchResponse)
	if err != nil {
		return []TrackInfo{}, err
	}

	return SearchResponse.Tracks.Items, nil
}

func SearchTracks(limit int, offset int, query string, token *auth.FreshToken) ([]TrackInfo, error) {
	data, err := Search("track", limit, offset, query, token)
	if err != nil {
		return []TrackInfo{}, err
	}

	return ParseTrackInfo(data)

}

func SearchAlbums(limit int, offset int, query string, token *auth.FreshToken) ([]AlbumInfo, error) {
	data, err := Search("album", limit, offset, query, token)
	if err != nil {
		return []AlbumInfo{}, err
	}

	var SearchResponse struct {
		Tracks struct {
			Items []AlbumInfo `json:"items"`
		} `json:"albums"`
	}

	err = json.Unmarshal(data, &SearchResponse)
	if err != nil {
		return []AlbumInfo{}, err
	}

	return SearchResponse.Tracks.Items, nil
}

func SearchArtist(limit int, offset int, query string, token *auth.FreshToken) ([]ArtistInfo, error) {
	data, err := Search("artist", limit, offset, query, token)
	if err != nil {
		return []ArtistInfo{}, err
	}

	var SearchResponse struct {
		Tracks struct {
			Items []ArtistInfo `json:"items"`
		} `json:"artists"`
	}

	err = json.Unmarshal(data, &SearchResponse)
	if err != nil {
		return []ArtistInfo{}, err
	}

	return SearchResponse.Tracks.Items, nil
}

func getIdFromUri(uri string) string {
	return strings.Split(uri, ":")[2]
}

type Playable interface {
	Play(token *auth.FreshToken) error
	GetDescription() string
}

func (ti TrackInfo) Play(token *auth.FreshToken) error {
	return PlayTrack(token, ti.Album.Uri, ti.Uri)
}

func (ti TrackInfo) GetDescription() string {
	var builder strings.Builder
	for _, artist := range ti.Artists {
		builder.WriteString(artist.Name)
		builder.WriteString(", ")
	}
	artists := builder.String()
	return fmt.Sprintf("%s - %s", artists[:len(artists)-1], ti.Name)
}

func (ai AlbumInfo) Play(token *auth.FreshToken) error {
	return PlayTrack(token, ai.Uri, "")
}

func (ai AlbumInfo) GetDescription() string {
	return fmt.Sprintf("%s", ai.Name)
}

func Interactive(token *auth.FreshToken, items []Playable) error {
	for i, item := range items {
		fmt.Println(strconv.Itoa(i+1), "-", item.GetDescription())
	}
	fmt.Println("Choose twin")
	var index int
	_, err := fmt.Scan(&index)
	if err != nil {
		return err
	}
	if index > len(items) {
		return fmt.Errorf("OUT OF BOUNDS")
	}
	return items[index-1].Play(token)
}
