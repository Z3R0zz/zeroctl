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

	err = initSwww()
	if err != nil {
		logrus.Errorf("Shit went sideways: %v", err)
		return err
	}

	cmd := exec.Command("swww", "img", "--transition-type", "grow", "--transition-pos", "0.0854,0.977", "--transition-step", "90", wallpaperPath)
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
			output, err := cmd.CombinedOutput()
			if err != nil {
				logrus.Errorf("Failed to execute command %s: %v - Output: %s", cmdStr, err, string(output))
			} else {
				logrus.Infof("Executed post-wallpaper command: %s", cmdStr)
			}
		}
	}

	logrus.Infof("Changed wallpaper to %s", wallpaper)
	return nil
}

func initSwww() error {
	cmd := exec.Command("swww", "query")
	if err := cmd.Run(); err == nil {
		// swww is already running because our fuckass query didn't error out
		return nil
	}

	cmd = exec.Command("swww-daemon")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}

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
