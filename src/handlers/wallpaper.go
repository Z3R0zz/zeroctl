package handlers

import (
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// On hood this is some nasa shit just to change a wallpaper
func RandomWallpaper() error {
	wallpaper, err := getWallpaper()
	if err != nil {
		return err
	}

	wallpaperPath := os.Getenv("WALLPAPERS_DIR") + wallpaper

	cmd := exec.Command("waypaper", "--wallpaper", wallpaperPath)
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("wal", "-i", wallpaperPath)
	if err := cmd.Run(); err != nil {
		return err
	}

	afterCommandsStr := os.Getenv("AFTER_WALLPAPER_COMMANDS")
	if afterCommandsStr != "" {
		afterCommands := strings.Split(afterCommandsStr, ";")
		for _, cmdStr := range afterCommands {
			cmdStr = strings.TrimSpace(cmdStr)
			if cmdStr == "" {
				continue
			}

			parts := strings.Fields(cmdStr)
			if len(parts) == 0 {
				continue
			}

			cmd := exec.Command(parts[0], parts[1:]...)
			if err := cmd.Run(); err != nil {
				logrus.Errorf("Failed to execute command %s: %v", cmdStr, err)
			} else {
				logrus.Infof("Executed post-wallpaper command: %s", cmdStr)
			}
		}
	}

	logrus.Infof("Changed wallpaper to %s", wallpaper)
	return nil
}

func getWallpaper() (string, error) {
	entries, err := os.ReadDir(os.Getenv("WALLPAPERS_DIR"))
	if err != nil {
		return "", err
	}

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	return entries[r.Intn(len(entries))].Name(), nil
}
