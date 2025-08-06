package commands

import (
	"example.com/m/v2/auth"
)

type getAlbumTracksCommand struct{}

func (gc getAlbumTracksCommand) Name() string { return "tracks" }

func (gc getAlbumTracksCommand) Description() string { return "" }

func (gc getAlbumTracksCommand) Exec(token *auth.FreshToken, args []string) error {
	track, err := GetCurrentTrack(token)
	if err != nil {
		return err
	}

	albumTracks, err := GetAlbumTracks(token, track.Album.Uri)
	if err != nil {
		return err
	}

	temp := make([]Playable, len(albumTracks))
	for i, cur := range albumTracks {
		cur.Album = AlbumInfo{Uri: track.Album.Uri}
		temp[i] = cur
	}
	return Interactive(token, temp)
}

func init() {
	registerCommand(getAlbumTracksCommand{})
}
