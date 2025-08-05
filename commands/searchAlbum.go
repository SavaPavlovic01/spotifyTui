package commands

import (
	"fmt"

	"example.com/m/v2/auth"
)

type searchAlbumCommand struct{}

func (sc searchAlbumCommand) Name() string { return "sa" }

func (sc searchAlbumCommand) Description() string { return "" }

func (sc searchAlbumCommand) Exec(token *auth.FreshToken, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("pls send arguments")
	}

	albumInfo, err := SearchAlbums(10, 0, args[0], token)
	if err != nil {
		return err
	}

	return PlayTrack(token, albumInfo[0].Uri, "")
}

func init() {
	registerCommand(searchAlbumCommand{})
}
