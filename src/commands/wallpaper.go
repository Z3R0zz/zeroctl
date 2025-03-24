package commands

import (
	"errors"
	"os"
	"strings"
	"zeroctl/src/handlers"
	"zeroctl/src/types"

	"github.com/sirupsen/logrus"
)

func init() {
	types.RegisterCommand(types.Command{
		Name:        "wallpaper",
		Description: "Set a wallpaper or leave it blank to randomize",
		Handler: func(args []string) string {
			if len(args) < 1 {
				return randomize()
			}

			if len(args) > 1 {
				return "Please provide only one argument\n"
			}

			wallpaperPath := strings.TrimSpace(args[0])
			if wallpaperPath == "" {
				return randomize()
			}

			return setWallpaper(wallpaperPath)
		},
	})
}

func randomize() string {
	err := handlers.RandomWallpaper()
	if err != nil {
		logrus.Errorf("Failed to change wallpaper: %v", err)
		return "Failed to change wallpaper: " + err.Error()
	}

	return "Wallpaper changed successfully"
}

func setWallpaper(wallpaperPath string) string {
	if _, err := os.Stat(wallpaperPath); errors.Is(err, os.ErrNotExist) {
		return "Wallpaper does not exist, please check the path"
	}

	err := handlers.SetWallpaper(wallpaperPath)
	if err != nil {
		logrus.Errorf("Failed to set wallpaper: %v", err)
		return "Failed to set wallpaper: " + err.Error()
	}
	return "Wallpaper set successfully"
}
