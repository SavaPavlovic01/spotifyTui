package commands

import (
	"fmt"

	"example.com/m/v2/auth"
)

type searchSongCommand struct{}

func (sc searchSongCommand) Name() string { return "st" }

func (sc searchSongCommand) Description() string { return "" }

// TODO: take argument like -i, that prints results and then you type like 2 for second result
func (sc searchSongCommand) Exec(token *auth.FreshToken, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("send search params pls")
	}

	trackInfo, err := SearchTracks(10, 0, args[0], token)
	if err != nil {
		return err
	}
	return PlayTrack(token, trackInfo[0].Album.Uri, trackInfo[0].Uri)
}

func init() {
	registerCommand(searchSongCommand{})
}
