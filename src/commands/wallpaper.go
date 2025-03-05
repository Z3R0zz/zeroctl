package commands

import (
	"zeroctl/src/handlers"
	"zeroctl/src/types"

	"github.com/sirupsen/logrus"
)

func init() {
	types.RegisterCommand(types.Command{
		Name:        "wallpaper",
		Description: "Randomize wallpaper",
		Handler:     handleRandomWallpaper,
	})
}

func handleRandomWallpaper() string {
	err := handlers.RandomWallpaper()
	if err != nil {
		logrus.Errorf("Failed to change wallpaper: %v", err)
		return "Failed to change wallpaper: " + err.Error()
	}

	return "Wallpaper changed successfully"
}
